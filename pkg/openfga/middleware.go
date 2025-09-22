package openfga

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/cache"
	"github.com/OtchereDev/ris-common-sdk/pkg/roles"
	"github.com/OtchereDev/ris-common-sdk/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Filter       func(c *fiber.Ctx) bool
	Unauthorized fiber.Handler
	Secret       string
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
	Cache        cache.Cache
	OpenFGA      *OpenFGAService // Add OpenFGA service
}

type PermissionConfig struct {
	Permission roles.Permission
	OpenFGA    *OpenFGAService
	Optional   bool // If true, continues even if permission check fails (logs warning)
}

var DefaultConfig = Config{
	Filter:       nil,
	Unauthorized: nil,
	Secret:       "TestSenior",
	Decode:       nil,
	Cache:        nil,
	OpenFGA:      nil,
}

var (
	AuthMiddleware              func(*fiber.Ctx) error
	PermissionNotFulfilledError = &fiber.Map{
		"status": fiber.StatusForbidden,
		"data": &fiber.Map{
			"message": "You dont have permission to perform this action",
		},
	}
	PermissionCheckFailedError = &fiber.Map{
		"status": fiber.StatusInternalServerError,
		"data": &fiber.Map{
			"message": "Permission check failed",
		},
	}
)

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

func SetupMiddleAuthMiddleware(s string, cache cache.Cache, fgaService *OpenFGAService) {
	if cache == nil {
		AuthMiddleware = New(Config{Secret: s, Cache: nil, OpenFGA: fgaService})
	} else {
		AuthMiddleware = New(Config{Secret: s, Cache: cache, OpenFGA: fgaService})
	}
}

// OpenFGA-based permission middleware
func RequirePermission(config PermissionConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.OpenFGA == nil {
			log.Printf("OpenFGA service not configured for permission check")
			if config.Optional {
				return c.Next()
			}
			return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
		}

		user, err := utils.SerializeRequestUser(c)
		if err != nil {
			log.Printf("Failed to serialize user from request: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status": fiber.StatusUnauthorized,
				"data":   &fiber.Map{"message": "Invalid user context"},
			})
		}

		userIDStr := fmt.Sprintf("%d", user.UserID)

		// Check permission with OpenFGA
		hasPermission, err := config.OpenFGA.CheckPermission(c.Context(), userIDStr, config.Permission)
		if err != nil {
			log.Printf("OpenFGA permission check failed for user %s, permission %s: %v", userIDStr, config.Permission, err)
			if config.Optional {
				log.Printf("Permission check failed but marked as optional, continuing...")
				return c.Next()
			}
			return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
		}

		if !hasPermission {
			log.Printf("User %s denied access to permission %s", userIDStr, config.Permission)
			return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
		}

		// Store permission info in context for downstream handlers
		c.Locals("checked_permission", config.Permission)
		c.Locals("user_id", user.UserID)

		return c.Next()
	}
}

// Batch permission check middleware
func RequireAnyPermission(fgaService *OpenFGAService, permissions []roles.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if fgaService == nil {
			log.Printf("OpenFGA service not configured for batch permission check")
			return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
		}

		user, err := utils.SerializeRequestUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status": fiber.StatusUnauthorized,
				"data":   &fiber.Map{"message": "Invalid user context"},
			})
		}

		userIDStr := fmt.Sprintf("%d", user.UserID)

		// Use batch check for efficiency
		results, err := fgaService.BatchCheckPermissions(c.Context(), userIDStr, permissions)
		if err != nil {
			log.Printf("OpenFGA batch permission check failed for user %s: %v", userIDStr, err)
			return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
		}

		// Check if user has at least one of the required permissions
		hasAnyPermission := false
		var allowedPermissions []roles.Permission

		for permission, allowed := range results {
			if allowed {
				hasAnyPermission = true
				allowedPermissions = append(allowedPermissions, permission)
			}
		}

		if !hasAnyPermission {
			log.Printf("User %s denied access - no required permissions found", userIDStr)
			return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
		}

		// Store allowed permissions in context
		c.Locals("allowed_permissions", allowedPermissions)
		c.Locals("user_id", user.UserID)

		return c.Next()
	}
}

// Resource-specific permission middleware
func RequireResourcePermission(fgaService *OpenFGAService, relation, resourceType string, resourceIDParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if fgaService == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
		}

		user, err := utils.SerializeRequestUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status": fiber.StatusUnauthorized,
				"data":   &fiber.Map{"message": "Invalid user context"},
			})
		}

		// Extract resource ID from URL parameters
		resourceID := c.Params(resourceIDParam)
		if resourceID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"status": fiber.StatusBadRequest,
				"data":   &fiber.Map{"message": "Resource ID is required"},
			})
		}

		userIDStr := fmt.Sprintf("%d", user.UserID)

		// Check resource-specific permission
		hasPermission, err := fgaService.CheckResourcePermission(c.Context(), userIDStr, relation, resourceType, resourceID)
		if err != nil {
			log.Printf("OpenFGA resource permission check failed for user %s on %s:%s: %v", userIDStr, resourceType, resourceID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
		}

		if !hasPermission {
			log.Printf("User %s denied access to %s %s:%s", userIDStr, relation, resourceType, resourceID)
			return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
		}

		// Store resource info in context
		c.Locals("resource_type", resourceType)
		c.Locals("resource_id", resourceID)
		c.Locals("relation", relation)
		c.Locals("user_id", user.UserID)

		return c.Next()
	}
}

// Backward compatibility - Legacy role-based middleware
func RoleMiddleware(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := utils.SerializeRequestUser(c)

		if slices.Contains(roles, user.UserType) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
	}
}

// OpenFGA-enhanced middleware equivalents
func ReportingTeamMiddleware(fgaService *OpenFGAService) fiber.Handler {
	return RequireAnyPermission(fgaService, []roles.Permission{
		roles.ViewReports, roles.CreateReports, roles.EditReports, roles.VerifyReports,
	})
}

func StaffMiddleware(fgaService *OpenFGAService) fiber.Handler {
	return RequireAnyPermission(fgaService, []roles.Permission{
		roles.ViewAppointments, roles.ViewReports, roles.ViewFinancials, roles.ViewUsers,
	})
}

func AdminMiddleware(fgaService *OpenFGAService) fiber.Handler {
	return RequirePermission(PermissionConfig{
		Permission: roles.ManageSystem,
		OpenFGA:    fgaService,
		Optional:   false,
	})
}

func FinanceMiddleware(fgaService *OpenFGAService) fiber.Handler {
	return RequireAnyPermission(fgaService, []roles.Permission{
		roles.ViewFinancials, roles.EditFinancials, roles.ProcessPayments, roles.CreateInvoices,
	})
}

// Hybrid middleware - checks both role and OpenFGA permissions
func HybridRolePermissionMiddleware(fgaService *OpenFGAService, roles []string, permissions []roles.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.SerializeRequestUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status": fiber.StatusUnauthorized,
				"data":   &fiber.Map{"message": "Invalid user context"},
			})
		}

		// Check role first (fast path)
		if slices.Contains(roles, user.UserType) {
			return c.Next()
		}

		// If role doesn't match, check OpenFGA permissions
		if fgaService != nil && len(permissions) > 0 {
			userIDStr := fmt.Sprintf("%d", user.UserID)

			results, err := fgaService.BatchCheckPermissions(c.Context(), userIDStr, permissions)
			if err != nil {
				log.Printf("Hybrid permission check failed for user %s: %v", userIDStr, err)
				return c.Status(fiber.StatusInternalServerError).JSON(PermissionCheckFailedError)
			}

			// Check if user has any of the required permissions
			for _, allowed := range results {
				if allowed {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
	}
}

// Utility function to check multiple permissions and return detailed results
func CheckUserPermissions(c *fiber.Ctx, fgaService *OpenFGAService, permissions []roles.Permission) (map[roles.Permission]bool, error) {
	user, err := utils.SerializeRequestUser(c)
	if err != nil {
		return nil, err
	}

	userIDStr := fmt.Sprintf("%d", user.UserID)
	return fgaService.BatchCheckPermissions(c.Context(), userIDStr, permissions)
}

// Middleware factory for creating permission-based route groups
func CreatePermissionGroup(fgaService *OpenFGAService, permission roles.Permission) fiber.Handler {
	return RequirePermission(PermissionConfig{
		Permission: permission,
		OpenFGA:    fgaService,
		Optional:   false,
	})
}

// Audit middleware that logs permission checks
func AuditPermissionMiddleware(fgaService *OpenFGAService, permission roles.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.SerializeRequestUser(c)
		if err != nil {
			return c.Next()
		}

		userIDStr := fmt.Sprintf("%d", user.UserID)

		hasPermission, err := fgaService.CheckPermission(c.Context(), userIDStr, permission)

		// Log the permission check regardless of result
		log.Printf("AUDIT: User %s (%s) attempted %s on %s %s - Result: %v, Error: %v",
			userIDStr, user.UserType, permission, c.Method(), c.Path(), hasPermission, err)

		if err != nil || !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(PermissionNotFulfilledError)
		}

		return c.Next()
	}
}
