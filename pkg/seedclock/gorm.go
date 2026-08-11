package seedclock

import (
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// The callback names, exported as constants so a service that needs to order its own
// callbacks around these can name them rather than guess.
const (
	CreateCallback = "seedclock:create"
	UpdateCallback = "seedclock:update"
)

// RegisterGORM installs callbacks that stamp the audit columns from the context.
//
// With this registered, a repository backdating a write sets the time once on the context
// it already hands to WithContext, and every row created or updated under that context
// carries it. Without it, the same effect needs an explicit CreatedAt assignment on every
// model in every write path, which is a long tail of easy omissions in code where the
// omission is invisible until a chart is wrong.
//
// Only zero fields are stamped, so an explicit assignment in the repository still wins —
// the callback fills in what the caller did not state, and never overrides what it did.
//
// Inert unless seed mode is armed: FromContext returns nothing in production, so both
// callbacks return immediately and GORM's own autoCreateTime handling applies untouched.
func RegisterGORM(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:create").
		Register(CreateCallback, stampCreate); err != nil {
		return err
	}

	return db.Callback().Update().Before("gorm:update").
		Register(UpdateCallback, stampUpdate)
}

// autoFields returns the schema's automatic timestamp fields.
//
// Selected from GORM's own metadata rather than by field name. The obvious implementation
// looks for CreatedAt and UpdatedAt, and misses every domain timestamp a service declares
// the same way — finance stamps a payment's PaidAt with autoCreateTime, the lab stamps an
// order's OrderedAt. Those are exactly the columns a revenue chart and a volume chart group
// by, so a name-based version backdates the audit columns and leaves the reporting ones at
// the seed run, which is the failure that looks like it worked.
//
// updateOnly narrows to the fields GORM would touch on an UPDATE.
func autoFields(db *gorm.DB, updateOnly bool) []*schema.Field {
	var out []*schema.Field

	for _, field := range db.Statement.Schema.Fields {
		if field.AutoUpdateTime > 0 || (!updateOnly && field.AutoCreateTime > 0) {
			out = append(out, field)
		}
	}

	return out
}

// stampCreate fills the automatic timestamps before GORM builds the INSERT.
//
// GORM's own handling only fills a field that is still zero, so setting them here is
// preserved rather than overwritten downstream.
func stampCreate(db *gorm.DB) {
	at, ok := seedTime(db)
	if !ok {
		return
	}

	fields := autoFields(db, false)
	each(db.Statement.ReflectValue, func(rv reflect.Value) {
		stamp(db, rv, at, fields...)
	})
}

// stampUpdate keeps the automatic timestamps in step on a backdated write.
//
// Without it a row created in the past is immediately updated in the present, and anything
// sorting or filtering on updated_at sees the seed run rather than the history it is meant
// to portray.
func stampUpdate(db *gorm.DB) {
	at, ok := seedTime(db)
	if !ok {
		return
	}

	fields := autoFields(db, true)

	// Updates(map[string]any{...}) is the common form in these repositories, and GORM adds
	// the timestamp columns to the map itself. Overriding the map entries is the only thing
	// that reaches the generated SQL in that case.
	//
	// Driven off the schema rather than a literal "updated_at". Not every table has the
	// column — the sequence counters are updated by the same code path and have neither
	// audit column — and naming one that does not exist turns a working update into a
	// "no such column" failure.
	if dest, isMap := db.Statement.Dest.(map[string]any); isMap {
		for _, field := range fields {
			if _, set := dest[field.DBName]; !set {
				dest[field.DBName] = at
			}
		}
		return
	}

	each(db.Statement.ReflectValue, func(rv reflect.Value) {
		stamp(db, rv, at, fields...)
	})
}

func seedTime(db *gorm.DB) (time.Time, bool) {
	if db.Statement == nil || db.Statement.Schema == nil {
		return time.Time{}, false
	}

	return FromContext(db.Statement.Context)
}

// each applies fn to a single record or to every record of a batch.
func each(rv reflect.Value, fn func(reflect.Value)) {
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			fn(reflect.Indirect(rv.Index(i)))
		}
	case reflect.Struct:
		fn(rv)
	}
}

// stamp sets the given schema fields to at, leaving any that already hold a value.
func stamp(db *gorm.DB, rv reflect.Value, at time.Time, fields ...*schema.Field) {
	if !rv.IsValid() {
		return
	}

	for _, field := range fields {
		if _, zero := field.ValueOf(db.Statement.Context, rv); !zero {
			continue
		}

		// A failure here means the model's field is not a time, which is a programming
		// error rather than a runtime condition. Recorded on the statement so it surfaces
		// in the seed run instead of silently leaving the column at the current time.
		if err := field.Set(db.Statement.Context, rv, at); err != nil {
			_ = db.AddError(err)
		}
	}
}
