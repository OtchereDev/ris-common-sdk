package types

import (
	"slices"
	"testing"
)

// The web application reads the same role string off the same token, so these values are a
// contract with app/lib/permissons.ts. A rename on either side is a user who can see a
// screen the server then refuses.
func TestUserTypeValues(t *testing.T) {
	want := map[string]string{
		"FrontDesk":          "FRONTDESK",
		"Cashier":            "CASHIER",
		"Radiographer":       "RADIOGRAPHER",
		"ReportingAssistant": "REPORTINGASSISTANT",
		"Radiologist":        "RADIOLOGIST",
		"Accountant":         "ACCOUNTANT",
		"Admin":              "ADMIN",
		"Phlebotomist":       "PHLEBOTOMIST",
		"LabTechnician":      "LABTECHNICIAN",
		"LabScientist":       "LABSCIENTIST",
		"Pathologist":        "PATHOLOGIST",
		"LabManager":         "LABMANAGER",
	}

	got := map[string]string{
		"FrontDesk":          UserTypes.FrontDesk,
		"Cashier":            UserTypes.Cashier,
		"Radiographer":       UserTypes.Radiographer,
		"ReportingAssistant": UserTypes.ReportingAssistant,
		"Radiologist":        UserTypes.Radiologist,
		"Accountant":         UserTypes.Accountant,
		"Admin":              UserTypes.Admin,
		"Phlebotomist":       UserTypes.Phlebotomist,
		"LabTechnician":      UserTypes.LabTechnician,
		"LabScientist":       UserTypes.LabScientist,
		"Pathologist":        UserTypes.Pathologist,
		"LabManager":         UserTypes.LabManager,
	}

	for name, expected := range want {
		if got[name] != expected {
			t.Errorf("%s = %q, want %q", name, got[name], expected)
		}
	}
}

// The facility service validates a role on create and on edit against AssignableUserTypes.
// Anything missing here is a role an administrator cannot assign.
func TestAssignableUserTypesCoversBothModules(t *testing.T) {
	for _, role := range append(append([]Alias{}, RadiologyUserTypes...), LabUserTypes...) {
		if !slices.Contains(AssignableUserTypes, role) {
			t.Errorf("%q is not assignable", role)
		}
	}

	// The seven roles facility accepted before the lab roles were added. Losing one would
	// lock existing staff out of user administration.
	for _, role := range []Alias{
		UserTypes.FrontDesk,
		UserTypes.Radiographer,
		UserTypes.Cashier,
		UserTypes.Accountant,
		UserTypes.Admin,
		UserTypes.Radiologist,
		UserTypes.ReportingAssistant,
	} {
		if !slices.Contains(AssignableUserTypes, role) {
			t.Errorf("previously assignable role %q was dropped", role)
		}
	}
}

func TestAssignableUserTypesHasNoDuplicates(t *testing.T) {
	seen := make(map[Alias]bool, len(AssignableUserTypes))
	for _, role := range AssignableUserTypes {
		if seen[role] {
			t.Errorf("%q appears more than once", role)
		}
		seen[role] = true
	}
}

// A role belongs to one module or the other, never both. getUserModules in the web
// application routes on exactly this split, so an overlap would put a phlebotomist in
// radiology.
func TestModuleRoleListsDoNotOverlap(t *testing.T) {
	for _, role := range LabUserTypes {
		if slices.Contains(RadiologyUserTypes, role) {
			t.Errorf("%q is in both module lists", role)
		}
	}
}
