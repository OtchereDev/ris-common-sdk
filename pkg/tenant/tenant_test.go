package tenant_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant/tenanttest"
)

const poll = 5 * time.Millisecond

// publish writes a record the way the agent does: a temporary file, then a rename.
func publish(t *testing.T, path string, r tenant.Record) {
	t.Helper()
	writeRaw(t, path, mustJSON(t, r))
}

func writeRaw(t *testing.T, path string, data []byte) {
	t.Helper()
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tenant-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmp.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		t.Fatal(err)
	}
}

func mustJSON(t *testing.T, r tenant.Record) []byte {
	t.Helper()
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func version(n int64, changes ...func(*tenant.Record)) tenant.Record {
	r := tenanttest.Record()
	r.Version = n
	for _, c := range changes {
		c(&r)
	}
	return r
}

// logs captures what the package logs, safely across the watcher goroutine.
type logs struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (l *logs) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

func (l *logs) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

func (l *logs) count(s string) int { return strings.Count(l.String(), s) }

func open(t *testing.T, path string, opts ...func(*tenant.Options)) (*tenant.Tenant, *logs) {
	t.Helper()
	captured := &logs{}
	o := tenant.Options{Path: path, PollInterval: poll, Logger: slog.New(slog.NewTextHandler(captured, nil))}
	for _, f := range opts {
		f(&o)
	}
	tn, err := tenant.Open(t.Context(), o)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(tn.Close)
	return tn, captured
}

// eventually waits for cond, failing the test after a second.
func eventually(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(poll)
	}
}

// settle gives the watcher a few polls to act on whatever it was just handed.
func settle() { time.Sleep(10 * poll) }

func TestOpenRefusesToStartWithoutAGoodRecord(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tenant.json")

	if _, err := tenant.Open(t.Context(), tenant.Options{Path: path}); err == nil {
		t.Fatal("a missing record must stop the service from starting")
	}

	writeRaw(t, path, []byte(`{"schema": 1, "key": "demo"`))
	if _, err := tenant.Open(t.Context(), tenant.Options{Path: path}); err == nil {
		t.Fatal("a truncated record must stop the service from starting")
	}

	bad := version(1)
	bad.Status = "suspended"
	publish(t, path, bad)
	if _, err := tenant.Open(t.Context(), tenant.Options{Path: path}); err == nil {
		t.Fatal("an invalid record must stop the service from starting")
	}

	if _, err := tenant.Open(t.Context(), tenant.Options{}); err == nil {
		t.Fatal("no path must be an error")
	}
}

func TestOpenRefusesAnotherTenantsRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(1))

	_, err := tenant.Open(t.Context(), tenant.Options{Path: path, ExpectKey: "pluslab"})
	if err == nil || !strings.Contains(err.Error(), `deployed for "pluslab"`) {
		t.Fatalf("want a refusal naming the expected tenant, got %v", err)
	}
}

func TestReadsTheExampleFile(t *testing.T) {
	tn, _ := open(t, "testdata/demo.json", func(o *tenant.Options) { o.PollInterval = 0 })

	if tn.Key() != "demo" || tn.Version() != 3 || tn.DisplayName() != "Demo Diagnostics" {
		t.Fatalf("got %s v%d %q", tn.Key(), tn.Version(), tn.DisplayName())
	}
	if !tn.Feature("whatsapp_reports") || tn.Feature("sso") {
		t.Fatal("features read wrongly")
	}
	if n, ok := tn.Limit("staff_seats"); !ok || n != 25 {
		t.Fatalf("staff_seats = %d, %v", n, ok)
	}
	if _, ok := tn.Limit("storage_gb"); ok {
		t.Fatal("a limit the record doesn't list is unlimited")
	}
	if n, day := tn.UsageMTD("sms"); n != 3120 || day != "2026-11-22" {
		t.Fatalf("usage sms = %d on %q", n, day)
	}
	if n, _ := tn.UsageMTD("whatsapp"); n != 0 {
		t.Fatal("an unrecorded meter is zero")
	}
	if p := tn.Profile(); p.Name != "Demo Diagnostics" || p.Document.Contact[0][1].Text != "+233 30 000 0000" {
		t.Fatalf("profile read wrongly: %+v", p)
	}
}

func TestModuleState(t *testing.T) {
	cases := []struct {
		status tenant.Status
		module string
		want   tenant.ModuleState
	}{
		{tenant.StatusActive, "radiology", tenant.ModuleEnabled},
		{tenant.StatusActive, "lab", tenant.ModuleReadOnly},
		{tenant.StatusActive, "ask", tenant.ModuleOff},
		{tenant.StatusActive, "telemedicine", tenant.ModuleOff}, // not listed
		{tenant.StatusRestricted, "radiology", tenant.ModuleEnabled},
		{tenant.StatusReadOnly, "radiology", tenant.ModuleReadOnly},
		{tenant.StatusReadOnly, "lab", tenant.ModuleReadOnly},
		{tenant.StatusReadOnly, "ask", tenant.ModuleOff},
		{tenant.StatusTerminated, "radiology", tenant.ModuleOff},
	}
	for _, c := range cases {
		tn := tenanttest.New(t,
			tenanttest.WithModule("lab", tenant.ModuleReadOnly),
			tenanttest.WithModule("ask", tenant.ModuleOff),
			tenanttest.WithStatus(c.status))
		if got := tn.ModuleState(c.module); got != c.want {
			t.Errorf("%s tenant, %s: got %s, want %s", c.status, c.module, got, c.want)
		}
	}
}

func TestReloadsAReplacedRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(1))
	tn, _ := open(t, path)

	var changes atomic.Int32
	tn.OnChange(func(*tenant.Tenant) { changes.Add(1) })

	publish(t, path, version(2, func(r *tenant.Record) {
		r.Features = []string{"sso"}
		r.Modules["lab"] = tenant.ModuleEnabled
	}))
	eventually(t, "version 2", func() bool { return tn.Version() == 2 })

	if !tn.Feature("sso") || tn.ModuleState("lab") != tenant.ModuleEnabled {
		t.Fatal("version 2 loaded but its contents weren't")
	}
	eventually(t, "the listener", func() bool { return changes.Load() == 1 })
}

func TestKeepsTheLastGoodRecord(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(2))
	tn, logged := open(t, path)

	cases := []struct {
		name  string
		write func()
		log   string
	}{
		{"half a file", func() { writeRaw(t, path, []byte(`{"schema":1,`)) }, "new tenant record refused"},
		{"an invalid record", func() {
			publish(t, path, version(3, func(r *tenant.Record) { r.Status = "paused" }))
		}, "new tenant record refused"},
		{"another tenant's record", func() {
			publish(t, path, version(4, func(r *tenant.Record) { r.Key, r.Profile.Key = "pluslab", "pluslab" }))
		}, "is for another tenant"},
		{"an older record", func() { publish(t, path, version(1)) }, "older than the one in use"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			before := logged.count(c.log)
			c.write()
			eventually(t, "the log line", func() bool { return logged.count(c.log) > before })
			if tn.Version() != 2 {
				t.Fatalf("version changed to %d", tn.Version())
			}
		})
	}

	// A good record after all that is still picked up.
	publish(t, path, version(5))
	eventually(t, "version 5", func() bool { return tn.Version() == 5 })
}

func TestLogsARefusedFileOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(1))
	_, logged := open(t, path)

	writeRaw(t, path, []byte(`not json`))
	eventually(t, "the refusal", func() bool { return logged.count("refused") == 1 })
	settle()
	if n := logged.count("refused"); n != 1 {
		t.Fatalf("the same bad file was logged %d times", n)
	}
}

func TestSurvivesTheFileDisappearing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(1))
	tn, logged := open(t, path)

	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	eventually(t, "the log line", func() bool { return logged.count("disappeared") == 1 })
	settle()
	if logged.count("disappeared") != 1 || tn.Version() != 1 {
		t.Fatalf("want one log line and version 1, got %d lines, version %d", logged.count("disappeared"), tn.Version())
	}

	publish(t, path, version(2))
	eventually(t, "version 2", func() bool { return tn.Version() == 2 })
}

func TestAListenerPanicDoesntStopReloading(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(1))
	tn, _ := open(t, path)
	tn.OnChange(func(*tenant.Tenant) { panic("listener bug") })

	publish(t, path, version(2))
	eventually(t, "version 2", func() bool { return tn.Version() == 2 })
	publish(t, path, version(3))
	eventually(t, "version 3", func() bool { return tn.Version() == 3 })
}

func TestProfileIsACopy(t *testing.T) {
	tn := tenanttest.New(t, func(r *tenant.Record) {
		r.Profile.Channels.FooterLinks = []tenant.Link{{Label: "Contact", Href: "mailto:a@b.c"}}
	})
	p := tn.Profile()
	p.Name = "changed"
	p.Channels.FooterLinks[0].Label = "changed"

	if q := tn.Profile(); q.Name == "changed" || q.Channels.FooterLinks[0].Label == "changed" {
		t.Fatal("changing a returned profile changed the tenant's")
	}
}

func TestConcurrentReadsDuringReloads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tenant.json")
	publish(t, path, version(1))
	tn, _ := open(t, path)

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = tn.Profile().Name
					_ = tn.ModuleState("radiology")
					_, _ = tn.Limit("staff_seats")
				}
			}
		}()
	}
	for v := int64(2); v <= 6; v++ {
		publish(t, path, version(v))
		eventually(t, "the next version", func() bool { return tn.Version() == v })
	}
	close(stop)
	wg.Wait()
}

func TestOptionsFromEnv(t *testing.T) {
	t.Setenv("TENANT_FILE", "")
	t.Setenv("TENANT_KEY", "")
	if o := tenant.OptionsFromEnv(); o.Path != tenant.DefaultPath || o.ExpectKey != "" {
		t.Fatalf("defaults: %+v", o)
	}
	t.Setenv("TENANT_FILE", "/run/tenant/tenant.json")
	t.Setenv("TENANT_KEY", " pluslab ")
	if o := tenant.OptionsFromEnv(); o.Path != "/run/tenant/tenant.json" || o.ExpectKey != "pluslab" {
		t.Fatalf("from env: %+v", o)
	}
}

func TestDefault(t *testing.T) {
	tn := tenanttest.New(t)
	tenant.SetDefault(tn)
	t.Cleanup(func() { tenant.SetDefault(nil) })
	if tenant.Default() != tn {
		t.Fatal("Default didn't return what SetDefault set")
	}
}

var _ io.Writer = (*logs)(nil)
