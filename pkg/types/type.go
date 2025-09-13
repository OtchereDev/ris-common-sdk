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
}

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
}

type AllUserParam struct {
	Page   int
	Limit  int
	Search string
	Role   string
}
