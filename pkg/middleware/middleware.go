package middleware

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/cache"
	. "github.com/OtchereDev/ris-common-sdk/pkg/types"
	. "github.com/OtchereDev/ris-common-sdk/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Secret       string
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
	Cache        cache.Cache
}

var DefaultConfig = Config{
	Filter:       nil,
	Unauthorized: nil,
	Secret:       "TestSenior",
	Decode:       nil,
	Cache:        nil,
}

var ReportingTeamMiddleware = RoleMiddleware([]string{
	UserTypes.Radiologist,
	UserTypes.ReportingAssistant,
	UserTypes.FrontDesk,
})

var DoctorMiddleware = RoleMiddleware([]string{"Doctor"})
var PatientMiddleware = RoleMiddleware([]string{"Patient"})
var OrganizationMiddleware = RoleMiddleware([]string{"Staff"})
var AuthMiddleware func(*fiber.Ctx) error
var PermissionNotFulfilledError = &fiber.Map{
	"status": fiber.StatusForbidden,
	"data": &fiber.Map{
		"message": "You dont have permission to perform this action"},
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return DefaultConfig
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = DefaultConfig.Filter
	}

	if strings.Trim(cfg.Secret, " ") == "" {
		cfg.Secret = os.Getenv("JWT_SECRET")
	}

	if cfg.Decode == nil {
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {
			authHeaders := c.Get("Authorization")

			if authHeaders == "" {
				return nil, errors.New("authorization Headers not provided")
			}

			token, err := jwt.Parse(strings.Split(authHeaders, " ")[1], func(t *jwt.Token) (interface{}, error) {

				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing error: %v", t.Header["alg"])
				}
				return []byte(cfg.Secret), nil
			})

			if err != nil {
				return nil, errors.New("invalid Token")
			}

			claim, ok := token.Claims.(jwt.MapClaims)

			if !(ok && token.Valid) {
				return nil, errors.New("invalid Token")
			}

			if expireAt, ok := claim["exp"]; ok && int64(expireAt.(float64)) < time.Now().UTC().Unix() {
				return nil, errors.New("token is expired")
			}

			if isRefresh, ok := claim["is_refresh"].(bool); ok && isRefresh {
				return nil, errors.New("invalid token, use the access token instead of the refresh token")
			}

			if cfg.Cache != nil {
				userId, ok := claim["user_id"].(string)
				if !ok {
					return nil, errors.New("invalid or missing user_id in token")
				}

				role, ok := claim["role"].(string)
				if !ok {
					return nil, errors.New("invalid or missing role in token")
				}

				id, err := strconv.ParseUint(userId, 10, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid user_id format: %w", err)
				}

				if cfg.Cache.IsUnSafe(uint32(id), role) {
					return nil, errors.New("this account has been disabled or deleted")
				}
			}

			return &claim, nil

		}
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"status": fiber.StatusUnauthorized, "data": &fiber.Map{"message": "Unauthorized"}})
		}
	}

	return cfg
}

func New(config Config) fiber.Handler {
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			fmt.Println("AuthMiddleware skipped")
		}

		claims, err := cfg.Decode(c)

		if err == nil {
			c.Locals("user", claims)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}

func SetupMiddleAuthMiddleware(s string, cache cache.Cache) {
	if cache == nil {
		AuthMiddleware = New(Config{Secret: s, Cache: nil})
	} else {
		AuthMiddleware = New(Config{Secret: s, Cache: cache})
	}
}

func RoleMiddleware(roles []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)

		if Contains(roles, user.UserType) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)

	}
}

func StaffMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)

		if Contains([]string{
			UserTypes.FrontDesk,
			UserTypes.Radiographer,
			UserTypes.Cashier,
			UserTypes.Accountant,
			UserTypes.Admin,
			UserTypes.Radiologist,
			UserTypes.ReportingAssistant,
		}, user.UserType) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)

	}
}

func AdminMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)

		if Contains([]string{
			UserTypes.Accountant,
			UserTypes.Admin,
		}, user.UserType) || user.IsAdmin {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)

	}
}

func FinanceMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)

		if Contains([]string{
			UserTypes.Accountant,
			UserTypes.Cashier,
		}, user.UserType) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)

	}
}
