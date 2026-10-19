package types

import (
	"reflect"
	"testing"
)

func TestModulesForRole(t *testing.T) {
	cases := []struct {
		name    string
		role    Alias
		isAdmin bool
		want    []Module
	}{
		{"radiographer staffs radiology", UserTypes.Radiographer, false, []Module{Modules.Radiology}},
		{"pathologist staffs lab", UserTypes.Pathologist, false, []Module{Modules.Lab}},
		{"admin by role holds everything", UserTypes.Admin, false, AllModules},
		{"admin by flag holds everything", UserTypes.Cashier, true, AllModules},
		{"a role nobody has heard of grants nothing", "TIMELORD", false, nil},
		{"an empty role grants nothing", "", false, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ModulesForRole(tc.role, tc.isAdmin)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ModulesForRole(%q, %v) = %v, want %v", tc.role, tc.isAdmin, got, tc.want)
			}
		})
	}
}

func TestIntersectModules(t *testing.T) {
	cases := []struct {
		name     string
		entitled []Module
		granted  []Module
		want     []Module
	}{
		{
			// The case that matters commercially: a lab scientist at a facility that never
			// licensed lab reaches nothing, however the role was assigned.
			name:     "a grant the facility never licensed is worth nothing",
			entitled: []Module{Modules.Radiology},
			granted:  []Module{Modules.Lab},
			want:     []Module{},
		},
		{
			// And the other direction: licensing a module does not hand it to everyone.
			name:     "an entitlement nobody holds is worth nothing",
			entitled: []Module{Modules.Radiology, Modules.Lab},
			granted:  []Module{Modules.Radiology},
			want:     []Module{Modules.Radiology},
		},
		{
			name:     "an admin at a two-module facility holds both",
			entitled: []Module{Modules.Radiology, Modules.Lab},
			granted:  AllModules,
			want:     []Module{Modules.Radiology, Modules.Lab},
		},
		{
			name:     "a suspended facility entitling nothing yields nothing",
			entitled: nil,
			granted:  AllModules,
			want:     []Module{},
		},
		{
			name:     "order follows the entitlement, so the claim is stable",
			entitled: []Module{Modules.Lab, Modules.Radiology},
			granted:  []Module{Modules.Radiology, Modules.Lab},
			want:     []Module{Modules.Lab, Modules.Radiology},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntersectModules(tc.entitled, tc.granted)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("IntersectModules(%v, %v) = %v, want %v",
					tc.entitled, tc.granted, got, tc.want)
			}
		})
	}
}

func TestParseModules(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Module
	}{
		{"the rendered form", "radiology,lab", []Module{Modules.Radiology, Modules.Lab}},
		{"whitespace and case", " Radiology , LAB ", []Module{Modules.Radiology, Modules.Lab}},
		{"empty is no entitlement", "", []Module{}},
		{"trailing separator", "radiology,", []Module{Modules.Radiology}},
		{"duplicates collapse", "lab,lab", []Module{Modules.Lab}},
		{
			// A typo in a tenant record must not become an entitlement no service can
			// enforce, so an unknown value is dropped rather than carried into the claim.
			name: "an unknown module is dropped, not trusted",
			in:   "radiology,pathology",
			want: []Module{Modules.Radiology},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseModules(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ParseModules(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
