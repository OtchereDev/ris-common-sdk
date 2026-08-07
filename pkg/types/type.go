package types

import (
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
}

type AllUserParam struct {
	Page   int
	Limit  int
	Search string
	Role   string
}
