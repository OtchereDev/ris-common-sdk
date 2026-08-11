package seedclock

import (
	"testing"
	"time"
)

func TestResolveIgnoresInputWhenDisabled(t *testing.T) {
	t.Setenv(EnvVar, "")

	before := time.Now()
	got := Resolve("2024-03-01T09:00:00Z")

	if got.Before(before) {
		t.Fatalf("seed mode off must fall back to now, got %s", got)
	}
}

// The property the whole design rests on: a production binary cannot be talked into
// backdating a write, whatever the client sends.
func TestResolveIsInertInProduction(t *testing.T) {
	t.Setenv(EnvVar, "false")

	for _, input := range []string{
		"2024-03-01T09:00:00Z",
		"2020-01-01",
		"1999-12-31 23:59:59",
	} {
		before := time.Now().Add(-time.Second)
		if got := Resolve(input); got.Before(before) {
			t.Errorf("input %q was honoured with seed mode off: %s", input, got)
		}
	}
}

func TestResolveHonoursInputWhenEnabled(t *testing.T) {
	t.Setenv(EnvVar, "true")

	want := time.Date(2024, time.March, 1, 9, 0, 0, 0, time.UTC)
	got := Resolve("2024-03-01T09:00:00Z")

	if !got.Equal(want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestEnabledAcceptsAnyCasing(t *testing.T) {
	for _, v := range []string{"true", "TRUE", "True", " true "} {
		t.Setenv(EnvVar, v)
		if !Enabled() {
			t.Errorf("%q should arm seed mode", v)
		}
	}

	for _, v := range []string{"", "1", "yes", "false", "on"} {
		t.Setenv(EnvVar, v)
		if Enabled() {
			t.Errorf("%q should not arm seed mode", v)
		}
	}
}

func TestParseAcceptsEveryDocumentedLayout(t *testing.T) {
	for _, input := range []string{
		"2024-03-01T09:00:00Z",
		"2024-03-01T09:00:00",
		"2024-03-01 09:00:00",
		"2024-03-01",
	} {
		if _, ok := Parse(input); !ok {
			t.Errorf("layout %q should parse", input)
		}
	}
}

func TestParseRejectsOutOfWindow(t *testing.T) {
	cases := map[string]string{
		"empty":        "",
		"nonsense":     "not a time",
		"before floor": "1998-06-01T00:00:00Z",
		"zero value":   "0001-01-01T00:00:00Z",
		"far future":   "2099-01-01T00:00:00Z",
	}

	for name, input := range cases {
		if _, ok := Parse(input); ok {
			t.Errorf("%s (%q) should be rejected", name, input)
		}
	}
}

// A future timestamp must not survive even with seed mode armed: it would sit at the top
// of every overdue and expiring query indefinitely.
func TestResolveRejectsFutureWhenEnabled(t *testing.T) {
	t.Setenv(EnvVar, "true")

	future := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
	before := time.Now().Add(-time.Second)

	if got := Resolve(future); got.After(time.Now().Add(time.Minute)) {
		t.Fatalf("future input was honoured: %s", got)
	} else if got.Before(before) {
		t.Fatalf("expected fallback to now, got %s", got)
	}
}

// Clock skew between the seeder and the service must not make a write fail to backdate.
func TestParseAllowsSmallSkew(t *testing.T) {
	nearFuture := time.Now().Add(2 * time.Minute).Format(time.RFC3339)

	if _, ok := Parse(nearFuture); !ok {
		t.Fatal("a timestamp inside the skew window should be accepted")
	}
}

func TestResolveAtUsesSuppliedFallback(t *testing.T) {
	t.Setenv(EnvVar, "true")

	fallback := time.Date(2023, time.May, 4, 12, 0, 0, 0, time.UTC)

	if got := ResolveAt("", fallback); !got.Equal(fallback) {
		t.Fatalf("empty input should yield the fallback, got %s", got)
	}

	want := time.Date(2024, time.March, 1, 9, 0, 0, 0, time.UTC)
	if got := ResolveAt("2024-03-01T09:00:00Z", fallback); !got.Equal(want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}
