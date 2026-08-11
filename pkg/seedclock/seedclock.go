// Package seedclock resolves the timestamp a write should carry, so that a staging
// environment can be seeded with data that has a believable history.
//
// The problem it solves: a demo environment populated by replaying real workflows through
// the API ends up with every row stamped at the moment the seeder ran. Volume charts show
// one spike, turnaround times are all near zero, and revenue trends are a single bar. The
// data is structurally correct and tells the viewer nothing.
//
// The fix is to let the seeder state, per write, when the event should be recorded as
// having happened. Handlers that today call time.Now() call Resolve instead, passing
// whatever the caller supplied.
//
// # Why this is not a production time-travel hole
//
// Resolve ignores its input entirely unless SEED_MODE is set. In a production binary the
// timestamp fields on the request payloads are inert: a client may send them, and every
// write still lands at time.Now(). That is the whole safety argument, and it is why the
// env check lives here rather than being repeated at each call site — there is one place
// to audit, and a handler cannot forget to guard.
//
// Even with SEED_MODE on, inputs outside a sane window are rejected rather than honoured,
// so a malformed or hostile value degrades to the current time instead of creating a row
// dated 2099 that then sits at the top of every "overdue" query forever.
//
// SEED_MODE must never be set in production. It belongs in the staging compose file and
// nowhere else.
package seedclock

import (
	"context"
	"os"
	"strings"
	"time"
)

// EnvVar is the switch that arms the whole package.
const EnvVar = "SEED_MODE"

// The layouts a caller may send, matching the set the lab service already accepts on
// collected_at so that a client formats a timestamp one way for the whole product.
var layouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// The window a seeded timestamp must fall in.
//
// The floor rejects zero values and obvious parse accidents: a year-1 or year-1970 date
// reaching the database is far more likely to be a bug than an intent to seed. The ceiling
// allows a little skew — the seeder's clock and the service's clock are not the same
// machine — but nothing meaningfully in the future, because a future timestamp is what
// turns "overdue" and "expiring" queries into permanent false positives.
var (
	floor      = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
	futureSkew = 5 * time.Minute
)

// Enabled reports whether seed mode is armed.
//
// Read from the environment on each call rather than memoised at init. A process cannot
// have its environment changed from outside, so this is no weaker a guarantee, and it
// keeps the package testable with t.Setenv and free of an exported mutator that would
// itself be a way to arm seed mode in production. These calls happen a handful of times
// per request, never in a loop, so the lookup cost does not matter.
func Enabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(EnvVar)), "true")
}

// Parse interprets a caller-supplied timestamp, reporting whether it was usable.
//
// Usable means both well-formed and inside the accepted window. It does not consult
// SEED_MODE: callers wanting the gate want Resolve.
func Parse(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err != nil {
			continue
		}
		if t.Before(floor) || t.After(time.Now().Add(futureSkew)) {
			return time.Time{}, false
		}
		return t, true
	}

	return time.Time{}, false
}

// Resolve returns the timestamp a write should carry.
//
// It is the seeded value when seed mode is armed and the input is usable, and time.Now()
// in every other case — production, an empty input, an unparseable one, or one outside the
// accepted window. Handlers replace a bare time.Now() with this and are otherwise
// unchanged.
func Resolve(input string) time.Time {
	return ResolveAt(input, time.Now())
}

// ResolveAt is Resolve with an explicit fallback.
//
// For the case where a handler has already established the time it would otherwise have
// used — several transitions in a single transaction that must share one timestamp, where
// calling time.Now() again would scatter them across a few microseconds and make the
// ordering within the transaction ambiguous.
func ResolveAt(input string, fallback time.Time) time.Time {
	if !Enabled() {
		return fallback
	}

	if t, ok := Parse(input); ok {
		return t
	}

	return fallback
}

// ========================================
// CARRYING THE RESOLVED TIME
// ========================================

type ctxKey struct{}

// WithTime returns a context carrying the timestamp seeded writes under it should land at.
//
// Paired with RegisterGORM, this is what saves a repository from assigning CreatedAt on
// every model it touches: the handler resolves the time once, puts it on the context it
// already passes to the database, and the audit columns follow. Missing one model is then
// impossible rather than merely unlikely, which matters because the one that gets missed
// is invariably the one a dashboard groups by.
//
// A no-op unless seed mode is armed, so a caller that forgets to check cannot smuggle a
// timestamp into a production write.
func WithTime(ctx context.Context, t time.Time) context.Context {
	if !Enabled() {
		return ctx
	}

	return context.WithValue(ctx, ctxKey{}, t)
}

// NowFrom returns the timestamp a write under ctx should carry.
//
// The seeded time when there is one, the current time otherwise. This is the call that
// replaces a bare time.Now() inside a repository: the handler has already resolved the
// payload's stamp onto the context, so everything below it agrees on one clock without
// having to be handed the payload.
func NowFrom(ctx context.Context) time.Time {
	if t, ok := FromContext(ctx); ok {
		return t
	}

	return time.Now()
}

// FromContext returns the seeded timestamp carried by ctx, if any.
func FromContext(ctx context.Context) (time.Time, bool) {
	if ctx == nil || !Enabled() {
		return time.Time{}, false
	}

	t, ok := ctx.Value(ctxKey{}).(time.Time)

	return t, ok
}
