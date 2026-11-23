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

func TestGrantedModules(t *testing.T) {
	m := Modules
	cases := []struct {
		name  string
		grant Grant
		want  []Module
	}{
		{"a radiographer staffs radiology only", Grant{Role: UserTypes.Radiographer}, []Module{m.Radiology}},
		{"a radiologist also gets AI reporting", Grant{Role: UserTypes.Radiologist}, []Module{m.Radiology, m.AIReporting}},
		{"a reporting assistant also gets AI reporting", Grant{Role: UserTypes.ReportingAssistant}, []Module{m.Radiology, m.AIReporting}},
		{"a marketer runs outreach and reads analytics", Grant{Role: UserTypes.Marketer}, []Module{m.Radiology, m.Analytics, m.Outreach}},
		{"an accountant reaches outreach", Grant{Role: UserTypes.Accountant}, []Module{m.Radiology, m.Outreach}},
		{"a cashier with analytics access gets analytics and ask", Grant{Role: UserTypes.Cashier, CanViewAnalytics: true}, []Module{m.Radiology, m.Analytics, m.Ask}},
		{"a lab scientist staffs lab", Grant{Role: UserTypes.LabScientist}, []Module{m.Lab}},
		{"a lab manager with analytics access", Grant{Role: UserTypes.LabManager, CanViewAnalytics: true}, []Module{m.Lab, m.Analytics, m.Ask}},
		{"an admin by role holds everything", Grant{Role: UserTypes.Admin}, AllModules},
		{"an admin by flag holds everything", Grant{Role: UserTypes.Cashier, IsAdmin: true}, AllModules},
		{"a referring doctor gets the doctor app", Grant{Role: RoleDoctor}, []Module{m.DoctorApp}},
		{"a doctor's assistant gets the doctor app", Grant{Role: RoleDoctorAssistant}, []Module{m.DoctorApp}},
		{"a patient gets the patient app", Grant{Role: RolePatient}, []Module{m.PatientApp}},
		{"an organisation user gets organisations", Grant{Role: RoleOrganisationUser}, []Module{m.Organisations}},
		{"non-staff roles match whatever their case", Grant{Role: "patient"}, []Module{m.PatientApp}},
		{"analytics access means nothing to a doctor", Grant{Role: RoleDoctor, CanViewAnalytics: true}, []Module{m.DoctorApp}},
		{"analytics access means nothing without a role", Grant{CanViewAnalytics: true}, nil},
		{"a role nobody has heard of grants nothing", Grant{Role: "TIMELORD"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := GrantedModules(tc.grant); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("GrantedModules(%+v) = %v, want %v", tc.grant, got, tc.want)
			}
		})
	}
}

// Every role a staff user may be assigned must grant something, or someone given that role
// can reach nothing once modules are enforced.
func TestEveryAssignableRoleGrantsAModule(t *testing.T) {
	for _, role := range AssignableUserTypes {
		if len(GrantedModules(Grant{Role: role})) == 0 {
			t.Errorf("role %s grants no module", role)
		}
	}
}

func TestEveryModuleIsGrantedToSomeoneBesidesAdmins(t *testing.T) {
	roles := append(append([]Alias{}, AssignableUserTypes...), RoleDoctor, RoleDoctorAssistant, RolePatient, RoleOrganisationUser)
	for _, module := range AllModules {
		found := false
		for _, role := range roles {
			for _, flag := range []bool{false, true} {
				if HasGrant(Grant{Role: role, CanViewAnalytics: flag}, module) {
					found = true
				}
			}
		}
		if !found {
			t.Errorf("only admins can ever hold %s", module)
		}
	}
}
