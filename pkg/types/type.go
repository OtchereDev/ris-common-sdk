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
	Radiology Module
	Lab       Module
}

// Modules is the vocabulary. It is closed on purpose: it appears in three places that must
// agree - here, the enum in the infra repo's tenants/schema.json, and MS_AUTH_TENANT_MODULES
// in a tenant's rendered env - and a value that exists in only one of them entitles nobody.
var Modules = &moduleList{
	Radiology: "radiology",
	Lab:       "lab",
}

// AllModules is every module this build knows about.
var AllModules = []Module{Modules.Radiology, Modules.Lab}

// moduleByRole maps a role to the module it staffs. This is the grant layer: what a person
// could hold, before the facility's entitlement is applied.
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

// ModulesForRole is what a role grants, before entitlement.
//
// An admin grants every module rather than only radiology. Administration is not a
// radiology job that happens to appear in that list; an admin of a facility that licensed
// lab has to be able to administer it. Any other role grants exactly the one module it
// staffs, and an unrecognised role grants nothing - a role the service has never heard of
// is not a reason to hand out access.
func ModulesForRole(role Alias, isAdmin bool) []Module {
	if isAdmin || role == UserTypes.Admin {
		return append([]Module{}, AllModules...)
	}
	if module, ok := moduleByRole[role]; ok {
		return []Module{module}
	}
	return nil
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
