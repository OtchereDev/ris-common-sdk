package types

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Alias = string

type list struct {
	FrontDesk          Alias
	Cashier            Alias
	Radiographer       Alias
	ReportingAssistant Alias
	Radiologist        Alias
	Accountant         Alias
	Admin              Alias
	// Marketer signs up referring doctors and runs marketing: referring doctors,
	// hospitals and organisations, campaigns and announcements, and referral analytics,
	// without the rest of what an Admin can do.
	Marketer Alias

	// Lab module roles. Module access is a separate axis from role, so these sit alongside
	// the radiology set rather than replacing it: a user is granted radiology, lab, or
	// both. Mirrors section 7 of the LMIS scope, and app/lib/permissons.ts in the web
	// application, which reads the same value off the same token.
	Phlebotomist  Alias
	LabTechnician Alias
	LabScientist  Alias
	Pathologist   Alias
	LabManager    Alias
}

type NotificationEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Enum for public use
var UserTypes = &list{
	FrontDesk:          "FRONTDESK",
	Cashier:            "CASHIER",
	Radiographer:       "RADIOGRAPHER",
	ReportingAssistant: "REPORTINGASSISTANT",
	Radiologist:        "RADIOLOGIST",
	Accountant:         "ACCOUNTANT",
	Admin:              "ADMIN",
	Marketer:           "MARKETER",

	Phlebotomist:  "PHLEBOTOMIST",
	LabTechnician: "LABTECHNICIAN",
	LabScientist:  "LABSCIENTIST",
	Pathologist:   "PATHOLOGIST",
	LabManager:    "LABMANAGER",
}

// RadiologyUserTypes are the roles that staff the imaging side.
var RadiologyUserTypes = []Alias{
	UserTypes.FrontDesk,
	UserTypes.Cashier,
	UserTypes.Radiographer,
	UserTypes.ReportingAssistant,
	UserTypes.Radiologist,
	UserTypes.Accountant,
	UserTypes.Admin,
	UserTypes.Marketer,
}

// LabUserTypes are the roles that staff the lab.
var LabUserTypes = []Alias{
	UserTypes.Phlebotomist,
	UserTypes.LabTechnician,
	UserTypes.LabScientist,
	UserTypes.Pathologist,
	UserTypes.LabManager,
}

// AssignableUserTypes is every role a staff user may hold. Services validating a role on
// create or edit should check against this rather than against a list written out at the
// call site, so a new role reaches every service at once instead of being accepted by one
// and rejected by another.
var AssignableUserTypes = append(append([]Alias{}, RadiologyUserTypes...), LabUserTypes...)

type ApiPaginatedResponse struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	HasError   bool              `json:"has_error"`
	Data       *PaginationResult `json:"data"`
}
type ApiReponse struct {
	StatusCode int        `json:"status_code"`
	Message    string     `json:"message"`
	HasError   bool       `json:"has_error"`
	Data       *fiber.Map `json:"data"`
}

type PaginationResult struct {
	Page       int     `json:"page"`
	TotalCount int64   `json:"total_count"`
	Limit      int     `json:"limit"`
	TotalPage  float64 `json:"total_page"`
	Data       any     `json:"data"`
}

type DisbaleDoctorPayload struct {
	Disabled bool `json:"disabled" validate:"required"`
}

type UpdateUserPassword struct {
	ID                 uint      `json:"id,omitempty" gorm:"primary_key"`
	FirstName          string    `json:"first_name" validate:"required"`
	LastName           string    `json:"last_name" validate:"required"`
	Email              string    `json:"email" validate:"email"`
	Password           string    `json:"-"`
	UserPassword       string    `json:"password,omitempty" validate:"omitempty" gorm:"-"`
	LastLogin          time.Time `json:"last_login,omitempty"`
	IsDisabled         bool      `json:"is_disabled" gorm:"default:false"`
	Role               string    `json:"role" validate:"required"`
	TitleOnReport      string    `json:"title_on_report"`
	HasPasswordUpdate  *bool     `json:"has_password_update" validate:"required"`
	IsAdmin            bool      `json:"is_admin" gorm:"default:false"`
	CanVerfiy          bool      `json:"can_verify" gorm:"default:false"`
	SpecialAppointment bool      `json:"special_appointment" gorm:"default:false"`
	IsProxy            bool      `json:"is_proxy"`
	PhoneNumber        string    `json:"phone_number" validate:"required,e164"`
	UseEmailOTP        bool      `json:"use_email_otp" gorm:"default:false"`
	CenterID           *uint     `json:"center_id,omitempty" validate:"required"`

	// CanEditCompletedAppointment lets the holder adjust the procedure on an
	// appointment that has already been examined. It stops at report verification,
	// which stays blocked for everyone.
	CanEditCompletedAppointment bool `json:"can_edit_completed_appointment" gorm:"default:false"`

	// CanViewAnalytics lets someone who isn't an Admin open Analytics.
	CanViewAnalytics bool `json:"can_view_analytics" gorm:"default:false"`
}

type AllUserParam struct {
	Page   int
	Limit  int
	Search string
	Role   string
}

// ---------------------------------------------------------------------------
// Module entitlements
//
// Access to a product is decided by three independent things, and effective access is the
// intersection of all three:
//
//	Entitlement  Did this facility license lab?   Tenant registry   changes with a contract
//	Grant        Does this user hold lab?         Auth service      changes with a role
//	Permission   May this role do this action?    Existing RBAC     rarely changes
//
// The first is per-tenant configuration, the second is per-user, and conflating them is
// how a build-time flag ends up being the only thing standing between a facility and a
// product it never bought. Auth intersects the first two into the token's modules claim;
// every service checks that claim at its own API boundary, because hiding a navigation
// item does not stop anyone calling the service directly.

type Module = string

type moduleList struct {
	Radiology     Module
	Lab           Module
	AIReporting   Module
	Analytics     Module
	Ask           Module
	Outreach      Module
	PatientApp    Module
	DoctorApp     Module
	Organisations Module
}

// Modules is the vocabulary. It is closed on purpose: a value that exists in only one of
// the places that must agree - here, the tenant record (pkg/tenant), and the web app's
// app/lib/permissons.ts - entitles nobody. The catalogue and what each covers is in
// docs/enterprise/entitlements-and-limits.md in the infra repository.
var Modules = &moduleList{
	Radiology:     "radiology",
	Lab:           "lab",
	AIReporting:   "ai_reporting",
	Analytics:     "analytics",
	Ask:           "ask",
	Outreach:      "outreach",
	PatientApp:    "patient_app",
	DoctorApp:     "doctor_app",
	Organisations: "organisations",
}

// AllModules is every module this build knows about.
var AllModules = []Module{
	Modules.Radiology, Modules.Lab, Modules.AIReporting, Modules.Analytics, Modules.Ask,
	Modules.Outreach, Modules.PatientApp, Modules.DoctorApp, Modules.Organisations,
}

// The roles that are not staff. Their tokens carry these exact values in "role".
const (
	RoleDoctor           Alias = "Doctor"   // a referring doctor
	RoleDoctorAssistant  Alias = "Delegate" // someone a referring doctor added as an assistant
	RolePatient          Alias = "Patient"
	RoleOrganisationUser Alias = "Staff" // a user of a corporate client's portal
)

// moduleByRole maps a staff role to the clinical module it staffs.
var moduleByRole = func() map[Alias]Module {
	m := make(map[Alias]Module, len(RadiologyUserTypes)+len(LabUserTypes))
	for _, role := range RadiologyUserTypes {
		m[role] = Modules.Radiology
	}
	for _, role := range LabUserTypes {
		m[role] = Modules.Lab
	}
	return m
}()

// extraByRole is what a role grants on top of the module it staffs.
var extraByRole = map[Alias][]Module{
	UserTypes.Radiologist:        {Modules.AIReporting},
	UserTypes.ReportingAssistant: {Modules.AIReporting},
	// Marketers run outreach, and read the referral measures in Analytics; which measures
	// is decided inside Analytics, not here.
	UserTypes.Marketer: {Modules.Outreach, Modules.Analytics},
	// The referral and marketing routes take Accountants as well as Marketers.
	UserTypes.Accountant: {Modules.Outreach},
}

// appByRole is the one module each non-staff role exists for.
var appByRole = map[string]Module{
	strings.ToLower(RoleDoctor):           Modules.DoctorApp,
	strings.ToLower(RoleDoctorAssistant):  Modules.DoctorApp,
	strings.ToLower(RolePatient):          Modules.PatientApp,
	strings.ToLower(RoleOrganisationUser): Modules.Organisations,
}

// Grant is what a token says about its holder that decides which modules they may hold.
type Grant struct {
	Role             Alias
	IsAdmin          bool
	CanViewAnalytics bool
}

// GrantedModules is what a person may hold, before the tenant's entitlement is applied.
//
// An admin holds every module: administering a facility that licensed lab means
// administering lab. Anyone else holds the module their role staffs, what that role adds,
// and Analytics and Ask if they were given Analytics access. A role nobody has heard of
// grants nothing - an unrecognised role is not a reason to hand out access.
//
// The order follows AllModules, so the claim minted from it is stable.
func GrantedModules(g Grant) []Module {
	if g.IsAdmin || g.Role == UserTypes.Admin {
		return append([]Module{}, AllModules...)
	}

	held := map[Module]bool{}
	staffModule, staff := moduleByRole[g.Role]
	if staff {
		held[staffModule] = true
		for _, extra := range extraByRole[g.Role] {
			held[extra] = true
		}
		// Corporate clients are served by the radiology desk: booking for them, billing
		// them, running their portal. Whether the tenant has the module decides the rest.
		if staffModule == Modules.Radiology {
			held[Modules.Organisations] = true
		}
	}
	if m, ok := appByRole[strings.ToLower(g.Role)]; ok {
		held[m] = true
	}
	// Analytics access is a staff permission. A doctor's or patient's token carrying the
	// flag by mistake gets nothing from it.
	if g.CanViewAnalytics && staff {
		held[Modules.Analytics] = true
		held[Modules.Ask] = true
	}

	var out []Module
	for _, m := range AllModules {
		if held[m] {
			out = append(out, m)
		}
	}
	return out
}

// HasGrant reports whether g holds module m.
func HasGrant(g Grant, m Module) bool {
	for _, held := range GrantedModules(g) {
		if held == m {
			return true
		}
	}
	return false
}

// ModulesForRole is GrantedModules for a role alone, without Analytics access. Kept for
// callers that only have the role; new code should use GrantedModules.
func ModulesForRole(role Alias, isAdmin bool) []Module {
	return GrantedModules(Grant{Role: role, IsAdmin: isAdmin})
}

// IntersectModules is the entitlement applied to the grant: what the facility licensed, and
// of that, what this person holds. Order is taken from entitled so the claim is stable
// across tokens regardless of how the grant was built.
func IntersectModules(entitled, granted []Module) []Module {
	held := make(map[Module]struct{}, len(granted))
	for _, m := range granted {
		held[m] = struct{}{}
	}

	out := make([]Module, 0, len(entitled))
	seen := make(map[Module]struct{}, len(entitled))
	for _, m := range entitled {
		if _, ok := held[m]; !ok {
			continue
		}
		if _, dup := seen[m]; dup {
			continue
		}
		seen[m] = struct{}{}
		out = append(out, m)
	}
	return out
}

// ParseModules reads the comma-separated form a tenant's rendered env carries in
// MS_AUTH_TENANT_MODULES. Unknown values are dropped rather than trusted: a typo in a
// tenant record must not become an entitlement no service can enforce.
func ParseModules(csv string) []Module {
	known := make(map[Module]struct{}, len(AllModules))
	for _, m := range AllModules {
		known[m] = struct{}{}
	}

	out := make([]Module, 0, len(AllModules))
	seen := make(map[Module]struct{}, len(AllModules))
	for _, part := range strings.Split(csv, ",") {
		value := strings.ToLower(strings.TrimSpace(part))
		if value == "" {
			continue
		}
		if _, ok := known[value]; !ok {
			continue
		}
		if _, dup := seen[value]; dup {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
