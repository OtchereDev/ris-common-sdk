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
