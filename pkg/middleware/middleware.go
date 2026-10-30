package middleware

import (
	"errors"
	"fmt"
	"log"
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

// RolesOrAdmin lets in the given roles, and anyone flagged is_admin. For routes
// that a role besides Admin needs, without opening them to all staff: the
// referral and marketing routes take Admins, Accountants and Marketers.
func RolesOrAdmin(roles ...Alias) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)
		if user != nil {
			if user.IsAdmin {
				return c.Next()
			}
			for _, r := range roles {
				if strings.EqualFold(user.UserType, string(r)) {
					return c.Next()
				}
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
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

// OrgAdminMiddleware allow only admin for organization
func OrgAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)

		if user == nil || user.OrganizationID == 0 || strings.ToUpper(user.OrganizationRole) != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)

		}

		return c.Next()
	}
}

// AdminOrgAdminMiddleware this allow either spectra admin or organization admin to go through
func AdminOrgAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := SerializeRequestUser(c)

		if user == nil {
			return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
		}

		// Global admin → allow
		if strings.ToUpper(user.UserType) == UserTypes.Admin || user.IsAdmin {
			return c.Next()
		}

		// Org admin → allow
		if user.OrganizationID != 0 && strings.ToUpper(user.OrganizationRole) == "ADMIN" {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
	}
}

// ---------------------------------------------------------------------------
// Module entitlements
//
// A licensed product is enforced here, at each service's API boundary, and not in
// navigation. Hiding a menu item does not stop anyone calling the service directly, and a
// build-time flag cannot be turned off without a rebuild — which means it is not
// enforcement at all, only decoration.

// ModuleEnforcement decides what happens to a token that carries no modules claim.
type ModuleEnforcement int

const (
	// ModuleWarn allows a token with no modules claim through, and logs it. This is the
	// safe setting during a rollout: access tokens live 30 minutes, so a service that
	// started rejecting claimless tokens before auth started minting them would lock every
	// user out for up to half an hour. It is a transitional setting, not a destination.
	ModuleWarn ModuleEnforcement = iota

	// ModuleEnforce rejects a token with no modules claim. Switch to this once auth is
	// deployed and every live token has been refreshed — a deliberate second step.
	ModuleEnforce
)

var moduleEnforcement = ModuleWarn

// SetModuleEnforcement is called once at boot, from the service's config. Left alone, a
// service enforces modules for tokens that carry them and admits ones that do not.
func SetModuleEnforcement(mode ModuleEnforcement) {
	moduleEnforcement = mode
}

// ParseModuleEnforcement reads the setting as it appears in a service's environment.
// Anything other than "enforce" is warn, because the failure mode of guessing wrong in that
// direction is a log line rather than an outage.
func ParseModuleEnforcement(value string) ModuleEnforcement {
	if strings.EqualFold(strings.TrimSpace(value), "enforce") {
		return ModuleEnforce
	}

	return ModuleWarn
}

// ErrUnauthenticated means no token reached the guard at all. Distinct from a refusal,
// because it tells a client to re-authenticate rather than to give up.
var ErrUnauthenticated = errors.New("unauthenticated")

// ErrModuleNotLicensed means a valid token does not carry the module.
var ErrModuleNotLicensed = errors.New("module not licensed")

var ModuleNotLicensedError = &fiber.Map{
	"status": fiber.StatusForbidden,
	"data": &fiber.Map{
		"message": "This module is not enabled for your facility"},
}

// CheckModule decides whether the request's token may reach a module, and writes nothing.
//
// It exists separately from RequireModule because a service whose routes already run
// through one authorization choke point wants to add this check there, and composing two
// middlewares that each call c.Next() would run the rest of the chain twice.
//
// The claim is read from the signed token and from nowhere else. The edge sets an
// X-Tenant-Key header for logging, and a service must never take an entitlement from a
// header — a header is only as trustworthy as everything between the client and here.
func CheckModule(c *fiber.Ctx, module string) error {
	user, err := SerializeRequestUser(c)

	if err != nil || user == nil {
		return ErrUnauthenticated
	}

	if user.HasModule(module) {
		return nil
	}

	// No claim at all is a token minted before entitlements existed. That is a different
	// thing from a claim of nothing, which is what a suspended facility mints and which is
	// refused in either mode.
	if !user.Modules.Present() && moduleEnforcement == ModuleWarn {
		log.Printf("module %q: admitting user %s on a token with no modules claim; set "+
			"module enforcement to enforce once auth is deployed", module, user.UserID)
		return nil
	}

	return ErrModuleNotLicensed
}

// DenyModule writes the refusal CheckModule decided on.
func DenyModule(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrUnauthenticated) {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status": fiber.StatusUnauthorized,
			"data":   &fiber.Map{"message": "Unauthorized"}})
	}

	return c.Status(fiber.StatusForbidden).JSON(ModuleNotLicensedError)
}

// RequireModule gates a route group on a module the facility has licensed and the user
// holds. Mount it after AuthMiddleware, which is what puts the token on the context.
func RequireModule(module string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := CheckModule(c, module); err != nil {
			return DenyModule(c, err)
		}

		return c.Next()
	}
}
