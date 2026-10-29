package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type JwtPayload struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	UserType         string `json:"role"`
	UserID           string `json:"user_id"`
	IsAdmin          bool   `json:"is_admin"`
	CanVerify        bool   `json:"can_verify"`
	IsRefresh        bool   `json:"is_refresh"`
	OrganizationID   uint   `json:"organization_id"`
	DeviceID         string `json:"device_id,omitempty"`
	OrganizationRole string `json:"org_role"`
	CenterID         *uint  `json:"center_id"`

	// CanEditCompletedAppointment allows adjusting the procedure on an appointment
	// that has already been examined. It never reaches past report verification —
	// that stays blocked for everyone, because commission has been paid by then.
	CanEditCompletedAppointment bool `json:"can_edit_completed_appointment"`

	// CanViewAnalytics lets someone who isn't an Admin open Analytics.
	CanViewAnalytics bool `json:"can_view_analytics"`

	// TenantKey is the facility this token was minted for. Present so a service can tell
	// that a valid token belongs to a different tenant than the one being addressed.
	TenantKey string `json:"tid"`

	// Modules is what this user may reach: what the facility licensed, intersected with
	// what the role grants. Minted by auth, checked by every service at its own API
	// boundary.
	Modules ModuleClaim `json:"modules"`
}

// ModuleClaim is the modules claim, and it distinguishes "no claim" from "a claim of
// nothing". The two look identical to a length check and mean opposite things: a token
// minted before entitlements existed carries no claim, while a suspended facility's token
// carries an empty one. Admitting the first during a rollout is a transition. Admitting the
// second would hand a suspended facility the product it stopped paying for.
type ModuleClaim struct {
	present bool
	values  []string
}

// Present reports whether the token carried a modules claim at all.
func (m ModuleClaim) Present() bool { return m.present }

// Values is the modules held, which may legitimately be none.
func (m ModuleClaim) Values() []string { return m.values }

// Has reports whether a module is held. Exact match, deliberately: the package's Contains
// helper is a substring test, and a substring test on entitlements would let "lab" be
// satisfied by a claim reading "notalab".
func (m ModuleClaim) Has(module string) bool {
	for _, held := range m.values {
		if held == module {
			return true
		}
	}

	return false
}

// UnmarshalJSON records that the claim was there. A JSON null is treated as absent rather
// than as an empty entitlement, because nothing mints null deliberately.
func (m *ModuleClaim) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if err := json.Unmarshal(data, &m.values); err != nil {
		return err
	}

	m.present = true

	return nil
}

// MarshalJSON keeps the claim round-tripping as the plain list it is on the wire.
func (m ModuleClaim) MarshalJSON() ([]byte, error) {
	if m.values == nil {
		return json.Marshal([]string{})
	}

	return json.Marshal(m.values)
}

// NewModuleClaim builds a claim that is present and holds exactly these modules. Used by
// auth when it mints a token, and by tests.
func NewModuleClaim(modules []string) ModuleClaim {
	if modules == nil {
		modules = []string{}
	}

	return ModuleClaim{present: true, values: modules}
}

// HasModule reports whether this token carries a module.
func (j *JwtPayload) HasModule(module string) bool {
	return j.Modules.Has(module)
}

func SerializeRequestUser(c *fiber.Ctx) (*JwtPayload, error) {
	var userPayload = c.Locals("user")

	if userPayload == nil {
		return nil, errors.New("unauthorized")
	}

	jsonData, error := json.Marshal(userPayload)

	var JwtPayload JwtPayload
	json.Unmarshal(jsonData, &JwtPayload)

	return &JwtPayload, error
}

func SerializeErrors(validationErr error) map[string]string {

	var errorResponse = make(map[string]string)

	for _, err := range validationErr.(validator.ValidationErrors) {
		fmt.Println(err.Field(), err.Tag(), err.Error(), err.ActualTag())

		errorResponse[err.Field()] = err.Error()
	}

	return errorResponse

}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 30
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}

	return false
}

// ParseUint parses a string to uint safely.
// Returns an error if the string is not a valid unsigned integer.
func ParseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
