package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DefaultPath is where a service looks for its record when TENANT_FILE isn't set. The
// directory, not the file, is what gets mounted; see the package documentation.
const DefaultPath = "/etc/spectra/tenant/tenant.json"

// DefaultPollInterval is how often a service checks whether its record has been replaced.
// A stat call every two seconds costs nothing measurable, and it works on every mount type,
// which file-change notifications through a Docker bind mount don't reliably do.
const DefaultPollInterval = 2 * time.Second

// Options says where the record is and whose it must be.
type Options struct {
	// Path to tenant.json. Required.
	Path string

	// ExpectKey, if set, is the tenant this service was deployed for. A record for any
	// other tenant is refused, at boot and on reload. This is the check that would have
	// stopped one laboratory's results printing another's telephone number.
	ExpectKey string

	// PollInterval defaults to DefaultPollInterval. Zero or negative disables reloading.
	PollInterval time.Duration

	// Logger defaults to slog.Default().
	Logger *slog.Logger
}

// OptionsFromEnv reads TENANT_FILE (default DefaultPath) and TENANT_KEY.
func OptionsFromEnv() Options {
	path := strings.TrimSpace(os.Getenv("TENANT_FILE"))
	if path == "" {
		path = DefaultPath
	}
	return Options{
		Path:         path,
		ExpectKey:    strings.TrimSpace(os.Getenv("TENANT_KEY")),
		PollInterval: DefaultPollInterval,
	}
}

// Tenant is a service's live view of its record. It is safe for concurrent use. Reads never
// block and never touch the disk; reloading happens in the background.
type Tenant struct {
	rec atomic.Pointer[Record]

	path      string
	expectKey string
	log       *slog.Logger

	mu        sync.Mutex
	listeners []func(*Tenant)
	stamp     stamp // the file last read, accepted or not

	stop     chan struct{}
	stopped  chan struct{}
	stopOnce sync.Once
}

// Open reads the record at opts.Path and, unless opts.PollInterval is zero or negative,
// watches it for replacement until ctx is done or Close is called.
//
// It fails if the file is missing or invalid, or names a tenant other than opts.ExpectKey.
// A service should not start without its record: running on a guess about whose deployment
// this is, is how one tenant's identity ends up on another's documents.
func Open(ctx context.Context, opts Options) (*Tenant, error) {
	if opts.Path == "" {
		return nil, errors.New("tenant: no record path; set TENANT_FILE")
	}
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	t := &Tenant{
		path:      opts.Path,
		expectKey: opts.ExpectKey,
		log:       opts.Logger.With("component", "tenant", "path", opts.Path),
		stop:      make(chan struct{}),
		stopped:   make(chan struct{}),
	}

	rec, st, err := t.read()
	if err != nil {
		return nil, err
	}
	t.rec.Store(rec)
	t.stamp = st
	t.log.Info("tenant record loaded", "tenant", rec.Key, "version", rec.Version, "status", rec.Status)

	if opts.PollInterval > 0 {
		go t.watch(ctx, opts.PollInterval)
	} else {
		close(t.stopped)
	}
	return t, nil
}

// FromRecord makes a Tenant from a record in memory, with no file behind it and no
// reloading. For the control plane, which holds records rather than files, and for tests.
func FromRecord(r Record) (*Tenant, error) {
	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("tenant: invalid record: %w", err)
	}
	t := &Tenant{log: slog.Default(), stop: make(chan struct{}), stopped: make(chan struct{})}
	close(t.stopped)
	t.rec.Store(&r)
	return t, nil
}

// Close stops watching for a new record. Reads keep returning the last record.
func (t *Tenant) Close() {
	t.stopOnce.Do(func() { close(t.stop) })
	<-t.stopped
}

// OnChange registers fn to run after a new record is accepted. It runs on the watcher's
// goroutine, so it should be quick; a panic in it is logged and doesn't stop reloading.
func (t *Tenant) OnChange(fn func(*Tenant)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.listeners = append(t.listeners, fn)
}

// Key is the tenant's identifier.
func (t *Tenant) Key() string { return t.rec.Load().Key }

// Version is the version of the record in use. Reported in /healthz so the control plane can
// show what each service is running against.
func (t *Tenant) Version() int64 { return t.rec.Load().Version }

// DisplayName is the customer's name for consoles and contracts. For what the product shows,
// use Profile().Name.
func (t *Tenant) DisplayName() string { return t.rec.Load().DisplayName }

// Status is the tenant's commercial status.
func (t *Tenant) Status() Status { return t.rec.Load().Status }

// Domains are the hostnames the tenant answers on; the first is canonical.
func (t *Tenant) Domains() []string { return slices.Clone(t.rec.Load().Domains) }

// ModuleState is what the tenant may do with module m right now, taking its status into
// account: a read_only tenant may only read, whatever it has licensed, and a terminated one
// may do nothing. A module the record doesn't list is off.
func (t *Tenant) ModuleState(m string) ModuleState {
	r := t.rec.Load()
	state := r.moduleOrOff(m)
	switch r.Status {
	case StatusTerminated:
		return ModuleOff
	case StatusReadOnly:
		if state == ModuleEnabled {
			return ModuleReadOnly
		}
	}
	return state
}

// Modules is every module the record lists, with its state as recorded. For the effective
// state, which also reflects the tenant's status, use ModuleState.
func (t *Tenant) Modules() map[string]ModuleState { return maps.Clone(t.rec.Load().Modules) }

// Feature reports whether a feature flag is on for this tenant. A flag the record doesn't
// list is off, so a flag removed from the control plane switches off rather than on.
func (t *Tenant) Feature(name string) bool { return slices.Contains(t.rec.Load().Features, name) }

// Features is every flag turned on.
func (t *Tenant) Features() []string { return slices.Clone(t.rec.Load().Features) }

// Limit returns the tenant's limit for name, and false if it has none, which means
// unlimited. A limit of zero is a real limit: nothing is allowed.
func (t *Tenant) Limit(name string) (int64, bool) {
	n, ok := t.rec.Load().Limits[name]
	return n, ok
}

// UsageMTD is the control plane's month-to-date count for a meter, and the day it was taken.
// Zero and "" if the record carries none. A service counting usage locally should use the
// higher of this and its own counter.
func (t *Tenant) UsageMTD(meter string) (int64, string) {
	u := t.rec.Load().UsageMTD
	if u == nil {
		return 0, ""
	}
	return u.Meters[meter], u.AsOf
}

// Profile is the tenant's brand, contacts, channels and identifiers. The returned value is a
// copy; changing it changes nothing.
func (t *Tenant) Profile() Profile { return t.rec.Load().Profile.clone() }

func (p Profile) clone() Profile {
	p.Channels.FooterLinks = slices.Clone(p.Channels.FooterLinks)
	p.Mobile.AppleAppIDs = slices.Clone(p.Mobile.AppleAppIDs)
	if p.Document.Contact != nil {
		lines := make([][]ContactItem, len(p.Document.Contact))
		for i, line := range p.Document.Contact {
			lines[i] = slices.Clone(line)
		}
		p.Document.Contact = lines
	}
	return p
}

// stamp identifies a version of the file on disk, so an unchanged file isn't read again and
// a rejected one isn't logged every poll.
type stamp struct {
	info    os.FileInfo
	modTime time.Time
	size    int64
}

func (s stamp) same(info os.FileInfo) bool {
	return s.info != nil && os.SameFile(s.info, info) && s.modTime.Equal(info.ModTime()) && s.size == info.Size()
}

func (t *Tenant) read() (*Record, stamp, error) {
	f, err := os.Open(t.path)
	if err != nil {
		return nil, stamp{}, fmt.Errorf("tenant: open record: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, stamp{}, fmt.Errorf("tenant: stat record: %w", err)
	}
	st := stamp{info: info, modTime: info.ModTime(), size: info.Size()}

	var r Record
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return nil, st, fmt.Errorf("tenant: %s is not a valid record: %w", t.path, err)
	}
	if err := r.Validate(); err != nil {
		return nil, st, fmt.Errorf("tenant: %s is not a valid record: %w", t.path, err)
	}
	if t.expectKey != "" && r.Key != t.expectKey {
		return nil, st, fmt.Errorf("tenant: %s is the record for %q, but this service was deployed for %q", t.path, r.Key, t.expectKey)
	}
	return &r, st, nil
}

func (t *Tenant) watch(ctx context.Context, every time.Duration) {
	defer close(t.stopped)
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stop:
			return
		case <-tick.C:
			t.reload()
		}
	}
}

// reload reads the file if it has changed and accepts it if it is valid, for the same
// tenant, and newer. Otherwise it keeps the record in use and logs why, once per file.
func (t *Tenant) reload() {
	info, err := os.Stat(t.path)
	if err != nil {
		t.mu.Lock()
		missingBefore := t.stamp.info == nil
		t.stamp = stamp{}
		t.mu.Unlock()
		if !missingBefore {
			t.log.Error("tenant record disappeared; keeping the one in use", "error", err)
		}
		return
	}

	t.mu.Lock()
	unchanged := t.stamp.same(info)
	t.mu.Unlock()
	if unchanged {
		return
	}

	next, st, err := t.read()
	if st.info == nil {
		// Unreadable after a successful stat: remember the stat, so a lasting problem such
		// as permissions is logged once rather than every poll.
		st = stamp{info: info, modTime: info.ModTime(), size: info.Size()}
	}

	t.mu.Lock()
	t.stamp = st
	t.mu.Unlock()

	cur := t.rec.Load()
	switch {
	case err != nil:
		t.log.Error("new tenant record refused; keeping the one in use", "version", cur.Version, "error", err)
		return
	case next.Key != cur.Key:
		t.log.Error("new tenant record is for another tenant; keeping the one in use",
			"tenant", cur.Key, "offered", next.Key)
		return
	case next.Version < cur.Version:
		t.log.Warn("new tenant record is older than the one in use; ignoring it",
			"version", cur.Version, "offered", next.Version)
		return
	case next.Version == cur.Version:
		return
	}

	t.rec.Store(next)
	t.log.Info("tenant record updated", "tenant", next.Key, "from", cur.Version, "to", next.Version, "status", next.Status)

	t.mu.Lock()
	listeners := slices.Clone(t.listeners)
	t.mu.Unlock()
	for _, fn := range listeners {
		t.notify(fn)
	}
}

func (t *Tenant) notify(fn func(*Tenant)) {
	defer func() {
		if p := recover(); p != nil {
			t.log.Error("tenant change listener panicked", "panic", p)
		}
	}()
	fn(t)
}

var defaultTenant atomic.Pointer[Tenant]

// SetDefault makes t the tenant Default returns. A service calls it once at boot.
func SetDefault(t *Tenant) { defaultTenant.Store(t) }

// Default is the tenant set by SetDefault, or nil before it is called. For code that can't
// be handed a *Tenant, such as middleware registered before the service is wired up.
func Default() *Tenant { return defaultTenant.Load() }
