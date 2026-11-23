package tenant_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant/tenanttest"
)

func TestExampleRecordIsValid(t *testing.T) {
	data, err := os.ReadFile("testdata/demo.json")
	if err != nil {
		t.Fatal(err)
	}
	var r tenant.Record
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	if err := r.Validate(); err != nil {
		t.Fatalf("testdata/demo.json should be valid: %v", err)
	}
	if r.Profile.Channels.SMSSenderID != "DemoDx" || r.Profile.DICOM.AccessionPrefix != "DMO" {
		t.Fatalf("profile fields not read by their camelCase names: %+v", r.Profile)
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name   string
		change func(*tenant.Record)
		want   string // a fragment of the error; "" means valid
	}{
		{"the demo record", func(*tenant.Record) {}, ""},
		{"a schema this build doesn't read", func(r *tenant.Record) { r.Schema = 2 }, "schema: is 2"},
		{"an uppercase key", func(r *tenant.Record) { r.Key, r.Profile.Key = "Demo", "Demo" }, "key:"},
		{"a key with an underscore", func(r *tenant.Record) { r.Key, r.Profile.Key = "demo_lab", "demo_lab" }, "key:"},
		{"version zero", func(r *tenant.Record) { r.Version = 0 }, "version:"},
		{"an unknown status", func(r *tenant.Record) { r.Status = "suspended" }, "status:"},
		{"a module state that isn't one", func(r *tenant.Record) { r.Modules["lab"] = "on" }, "modules.lab:"},
		{"a module this build doesn't know", func(r *tenant.Record) { r.Modules["telemedicine"] = tenant.ModuleEnabled }, ""},
		{"a feature listed twice", func(r *tenant.Record) { r.Features = []string{"sso", "sso"} }, "listed twice"},
		{"a negative limit", func(r *tenant.Record) { r.Limits = map[string]int64{"staff_seats": -1} }, "limits.staff_seats:"},
		{"a zero limit", func(r *tenant.Record) { r.Limits = map[string]int64{"staff_seats": 0} }, ""},
		{"a bad usage date", func(r *tenant.Record) { r.UsageMTD = &tenant.UsageMTD{AsOf: "22/11/2026"} }, "usage_mtd.as_of:"},
		{"an uppercase domain", func(r *tenant.Record) { r.Domains = []string{"Demo.example.com"} }, "domains[0]:"},
		{"another tenant's profile", func(r *tenant.Record) { r.Profile.Key = "pluslab" }, "profile.key:"},
		{"a profile with no name", func(r *tenant.Record) { r.Profile.Name = " " }, "profile.name:"},
		{"a colour that isn't hex", func(r *tenant.Record) { r.Profile.Primary = "purple" }, "profile.primary:"},
		{"a primary without an accent", func(r *tenant.Record) { r.Profile.Accent = "" }, "go together"},
		{"a sender id Hubtel rejects", func(r *tenant.Record) { r.Profile.Channels.SMSSenderID = "SpectraHealth" }, "Hubtel rejects"},
		{"an 11-character sender id", func(r *tenant.Record) { r.Profile.Channels.SMSSenderID = "SpectraHlth" }, ""},
		{"a relative system link", func(r *tenant.Record) { r.Profile.Links.System = "spectraris.com" }, "links.system:"},
		{"a lowercase currency", func(r *tenant.Record) { r.Profile.Locale.Currency = "ghs" }, "locale.currency:"},
		{"a study UID root that isn't one", func(r *tenant.Record) { r.Profile.DICOM.StudyUIDRoot = "1.2.x" }, "not a DICOM UID"},
		{"a study UID root too long", func(r *tenant.Record) {
			r.Profile.DICOM.StudyUIDRoot = "1." + strings.Repeat("2.", 20) + "3"
		}, "leaving no room"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := tenanttest.Record()
			c.change(&r)
			err := r.Validate()
			switch {
			case c.want == "" && err != nil:
				t.Fatalf("want valid, got %v", err)
			case c.want != "" && err == nil:
				t.Fatalf("want an error containing %q, got none", c.want)
			case c.want != "" && !strings.Contains(err.Error(), c.want):
				t.Fatalf("want an error containing %q, got %v", c.want, err)
			}
		})
	}
}

func TestValidateReportsEveryProblem(t *testing.T) {
	r := tenanttest.Record()
	r.Version = 0
	r.Status = "paused"
	r.Profile.Name = ""
	err := r.Validate()
	for _, want := range []string{"version:", "status:", "profile.name:"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("want %q among the errors, got %v", want, err)
		}
	}
}

func TestValidateForPublish(t *testing.T) {
	withModules := func(m map[string]tenant.ModuleState, version int64) *tenant.Record {
		r := tenanttest.Record()
		r.Modules, r.Version = m, version
		return &r
	}
	on := tenant.ModuleEnabled
	ro := tenant.ModuleReadOnly
	off := tenant.ModuleOff

	cases := []struct {
		name       string
		prev, next *tenant.Record
		want       string
	}{
		{"a first record", nil, withModules(map[string]tenant.ModuleState{"radiology": on}, 1), ""},
		{"no clinical module", nil, withModules(map[string]tenant.ModuleState{"analytics": on}, 1), "needs radiology or lab"},
		{"only off clinical modules", nil, withModules(map[string]tenant.ModuleState{"radiology": off}, 1), "needs radiology or lab"},
		{"a lab-only tenant", nil, withModules(map[string]tenant.ModuleState{"lab": on, "radiology": off}, 1), ""},
		{"lab dropped to read_only", withModules(map[string]tenant.ModuleState{"radiology": on, "lab": on}, 1),
			withModules(map[string]tenant.ModuleState{"radiology": on, "lab": ro}, 2), ""},
		{"lab dropped to off", withModules(map[string]tenant.ModuleState{"radiology": on, "lab": on}, 1),
			withModules(map[string]tenant.ModuleState{"radiology": on, "lab": off}, 2), "modules.lab: was enabled"},
		{"lab removed from the list", withModules(map[string]tenant.ModuleState{"radiology": on, "lab": on}, 1),
			withModules(map[string]tenant.ModuleState{"radiology": on}, 2), "modules.lab: was enabled"},
		{"read_only lab to off", withModules(map[string]tenant.ModuleState{"radiology": on, "lab": ro}, 1),
			withModules(map[string]tenant.ModuleState{"radiology": on, "lab": off}, 2), "modules.lab: was read_only"},
		{"a lab never licensed stays off", withModules(map[string]tenant.ModuleState{"radiology": on, "lab": off}, 1),
			withModules(map[string]tenant.ModuleState{"radiology": on, "lab": off}, 2), ""},
		{"analytics turned off", withModules(map[string]tenant.ModuleState{"radiology": on, "analytics": on}, 1),
			withModules(map[string]tenant.ModuleState{"radiology": on, "analytics": off}, 2), ""},
		{"the same version again", withModules(map[string]tenant.ModuleState{"radiology": on}, 4),
			withModules(map[string]tenant.ModuleState{"radiology": on}, 4), "must be higher"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.next.ValidateForPublish(c.prev)
			switch {
			case c.want == "" && err != nil:
				t.Fatalf("want publishable, got %v", err)
			case c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)):
				t.Fatalf("want an error containing %q, got %v", c.want, err)
			}
		})
	}
}

func TestCheckAcross(t *testing.T) {
	a := tenanttest.Record()
	b := tenanttest.Record()
	b.Key, b.Profile.Key = "pluslab", "pluslab"
	b.Profile.DICOM.AccessionPrefix = "PLB"
	b.Profile.Channels.PushTopic = "PluslabGeneral"
	b.Domains = []string{"pluslabgh.com"}

	if err := tenant.CheckAcross([]tenant.Record{a, b}); err != nil {
		t.Fatalf("distinct tenants should pass: %v", err)
	}

	b.Profile.DICOM.AccessionPrefix = a.Profile.DICOM.AccessionPrefix
	b.Profile.Channels.PushTopic = a.Profile.Channels.PushTopic
	b.Domains = a.Domains
	err := tenant.CheckAcross([]tenant.Record{a, b})
	for _, want := range []string{"accessionPrefix", "pushTopic", "claimed by both"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("want %q among the errors, got %v", want, err)
		}
	}
}
