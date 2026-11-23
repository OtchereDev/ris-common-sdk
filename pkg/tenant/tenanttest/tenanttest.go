// Package tenanttest builds tenant records for tests, so a service's tests don't need a
// tenant.json on disk or a real customer's details.
package tenanttest

import (
	"testing"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
)

// Record is a valid record for a made-up radiology tenant, "demo". Change what a test cares
// about and leave the rest.
func Record() tenant.Record {
	return tenant.Record{
		Schema:      tenant.SchemaVersion,
		Key:         "demo",
		Version:     1,
		PublishedAt: time.Date(2026, 11, 23, 9, 0, 0, 0, time.UTC),
		DisplayName: "Demo Diagnostics",
		Status:      tenant.StatusActive,
		Domains:     []string{"demo.spectraris.com"},
		Modules:     map[string]tenant.ModuleState{tenant.ModuleRadiology: tenant.ModuleEnabled},
		Profile: tenant.Profile{
			Key:     "demo",
			Name:    "Demo Diagnostics",
			Primary: "#31287f",
			Accent:  "#0b9bdc",
			Locale:  tenant.Locale{Tag: "en-GH", Currency: "GHS", CurrencyPrefix: "GHS", CurrencySymbol: "GH₵"},
			Contact: tenant.Contact{Address: "1 Test Street, Accra", Phone: "+233 30 000 0000", Email: "hello@demo.example"},
			Channels: tenant.Channels{
				SMSSenderID:      "DemoDx",
				EmailFromName:    "Demo Diagnostics",
				EmailFromAddress: "hello@demo.example",
				PushTopic:        "DemoGeneral",
			},
			Links:           tenant.Links{System: "https://demo.spectraris.com"},
			PatientIDPrefix: "DMO-",
			DICOM:           tenant.DICOM{AccessionPrefix: "DMO"},
		},
	}
}

// New returns a Tenant for the demo record after applying each change in order. It fails the
// test if the result isn't a valid record.
func New(tb testing.TB, changes ...func(*tenant.Record)) *tenant.Tenant {
	tb.Helper()
	r := Record()
	for _, change := range changes {
		change(&r)
	}
	t, err := tenant.FromRecord(r)
	if err != nil {
		tb.Fatalf("tenanttest: %v", err)
	}
	return t
}

// WithModule sets one module's state.
func WithModule(name string, state tenant.ModuleState) func(*tenant.Record) {
	return func(r *tenant.Record) {
		r.Modules = cloneModules(r.Modules)
		r.Modules[name] = state
	}
}

// WithStatus sets the tenant's status.
func WithStatus(s tenant.Status) func(*tenant.Record) {
	return func(r *tenant.Record) { r.Status = s }
}

// WithFeatures turns on the given feature flags, replacing any others.
func WithFeatures(names ...string) func(*tenant.Record) {
	return func(r *tenant.Record) { r.Features = append([]string(nil), names...) }
}

// WithLimit sets one limit.
func WithLimit(name string, n int64) func(*tenant.Record) {
	return func(r *tenant.Record) {
		limits := make(map[string]int64, len(r.Limits)+1)
		for k, v := range r.Limits {
			limits[k] = v
		}
		limits[name] = n
		r.Limits = limits
	}
}

func cloneModules(m map[string]tenant.ModuleState) map[string]tenant.ModuleState {
	out := make(map[string]tenant.ModuleState, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}
