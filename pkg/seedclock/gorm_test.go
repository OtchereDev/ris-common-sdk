package seedclock

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type widget struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	Kind      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func newDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := RegisterGORM(db); err != nil {
		t.Fatalf("register callbacks: %v", err)
	}

	if err := db.AutoMigrate(&widget{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return db
}

var seeded = time.Date(2024, time.March, 1, 9, 0, 0, 0, time.UTC)

func TestCreateIsBackdatedFromContext(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	ctx := WithTime(context.Background(), seeded)
	w := widget{Name: "tube"}

	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	var got widget
	if err := db.First(&got, w.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}

	if !got.CreatedAt.UTC().Equal(seeded) {
		t.Errorf("created_at = %s, want %s", got.CreatedAt.UTC(), seeded)
	}
	if !got.UpdatedAt.UTC().Equal(seeded) {
		t.Errorf("updated_at = %s, want %s", got.UpdatedAt.UTC(), seeded)
	}
}

// The production path: the same call with seed mode off must land at the current time.
func TestCreateIsNotBackdatedInProduction(t *testing.T) {
	t.Setenv(EnvVar, "false")
	db := newDB(t)

	// A context built while seed mode was off carries nothing, which is the first guard.
	ctx := WithTime(context.Background(), seeded)

	before := time.Now().Add(-time.Minute)
	w := widget{Name: "tube"}

	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	if w.CreatedAt.Before(before) {
		t.Fatalf("created_at was backdated with seed mode off: %s", w.CreatedAt)
	}
}

// Defence in depth: even a context that somehow carries a seed time is ignored when the
// flag is off, because FromContext consults it too.
func TestSeededContextIsIgnoredWhenFlagOff(t *testing.T) {
	t.Setenv(EnvVar, "true")
	ctx := WithTime(context.Background(), seeded)

	t.Setenv(EnvVar, "false")
	db := newDB(t)

	before := time.Now().Add(-time.Minute)
	w := widget{Name: "tube"}

	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	if w.CreatedAt.Before(before) {
		t.Fatalf("a stale seeded context was honoured: %s", w.CreatedAt)
	}
}

func TestExplicitAssignmentWins(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	explicit := time.Date(2022, time.July, 4, 8, 30, 0, 0, time.UTC)
	ctx := WithTime(context.Background(), seeded)
	w := widget{Name: "tube", CreatedAt: explicit}

	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	var got widget
	if err := db.First(&got, w.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}

	if !got.CreatedAt.UTC().Equal(explicit) {
		t.Errorf("created_at = %s, want the explicitly assigned %s", got.CreatedAt.UTC(), explicit)
	}
}

func TestBatchCreateIsBackdated(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	ctx := WithTime(context.Background(), seeded)
	batch := []widget{{Name: "a"}, {Name: "b"}, {Name: "c"}}

	if err := db.WithContext(ctx).Create(&batch).Error; err != nil {
		t.Fatalf("create batch: %v", err)
	}

	var got []widget
	if err := db.Find(&got).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d rows, want 3", len(got))
	}

	for _, w := range got {
		if !w.CreatedAt.UTC().Equal(seeded) {
			t.Errorf("%s created_at = %s, want %s", w.Name, w.CreatedAt.UTC(), seeded)
		}
	}
}

// Updates(map[string]any) is the dominant form in these repositories, and GORM fills
// updated_at itself, so the map entry is the only thing that reaches the SQL.
func TestMapUpdateIsBackdated(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	ctx := WithTime(context.Background(), seeded)
	w := widget{Name: "tube"}
	if err := db.WithContext(ctx).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	later := seeded.Add(6 * time.Hour)
	if err := db.WithContext(WithTime(context.Background(), later)).
		Model(&widget{}).Where("id = ?", w.ID).
		Updates(map[string]any{"kind": "edta"}).Error; err != nil {
		t.Fatalf("update: %v", err)
	}

	var got widget
	if err := db.First(&got, w.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}

	if got.Kind != "edta" {
		t.Errorf("kind = %q, want edta", got.Kind)
	}
	if !got.UpdatedAt.UTC().Equal(later) {
		t.Errorf("updated_at = %s, want %s", got.UpdatedAt.UTC(), later)
	}
}

// A domain timestamp declared with autoCreateTime must be backdated too.
//
// This is the case a name-based implementation gets wrong. Finance stamps a payment's
// PaidAt this way and the lab stamps an order's OrderedAt — the columns the revenue and
// volume charts group by. Backdating created_at while leaving those at the seed run
// produces a demo that looks seeded until somebody opens the reports.
type receipt struct {
	ID        uint `gorm:"primary_key"`
	Amount    float64
	PaidAt    time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func TestDomainTimestampsAreBackdated(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	if err := db.AutoMigrate(&receipt{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	ctx := WithTime(context.Background(), seeded)
	rec := receipt{Amount: 250}

	if err := db.WithContext(ctx).Create(&rec).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	var got receipt
	if err := db.First(&got, rec.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}

	if !got.PaidAt.UTC().Equal(seeded) {
		t.Errorf("paid_at = %s, want %s — a revenue chart would show the seed run",
			got.PaidAt.UTC(), seeded)
	}
	if !got.CreatedAt.UTC().Equal(seeded) {
		t.Errorf("created_at = %s, want %s", got.CreatedAt.UTC(), seeded)
	}
}

// A table without audit columns must still be updatable under a seeded context.
//
// The sequence counters these services use to issue order and specimen numbers have
// neither created_at nor updated_at, and they are incremented inside the same transaction
// as the backdated write. Naming a column that does not exist turned every seeded order
// into "no such column: updated_at".
type counter struct {
	ID    uint `gorm:"primary_key"`
	Scope string
	Value uint
}

func TestUpdateOnATableWithoutAuditColumns(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	if err := db.AutoMigrate(&counter{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	ctx := WithTime(context.Background(), seeded)

	c := counter{Scope: "order", Value: 1}
	if err := db.WithContext(ctx).Create(&c).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := db.WithContext(ctx).Model(&counter{}).Where("id = ?", c.ID).
		Updates(map[string]any{"value": 2}).Error; err != nil {
		t.Fatalf("update a table with no audit columns: %v", err)
	}

	var got counter
	if err := db.First(&got, c.ID).Error; err != nil {
		t.Fatalf("read back: %v", err)
	}
	if got.Value != 2 {
		t.Errorf("value = %d, want 2", got.Value)
	}
}

func TestUnseededContextLeavesGORMAlone(t *testing.T) {
	t.Setenv(EnvVar, "true")
	db := newDB(t)

	before := time.Now().Add(-time.Minute)
	w := widget{Name: "tube"}

	// Seed mode is armed, but this write carries no seeded time.
	if err := db.WithContext(context.Background()).Create(&w).Error; err != nil {
		t.Fatalf("create: %v", err)
	}

	if w.CreatedAt.Before(before) {
		t.Fatalf("created_at = %s, want the current time", w.CreatedAt)
	}
}
