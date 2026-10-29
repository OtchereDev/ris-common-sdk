package middleware

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	. "github.com/OtchereDev/ris-common-sdk/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// withUser stands in for AuthMiddleware, putting decoded claims on the context exactly as
// it does. Nil claims model an unauthenticated request reaching the guard.
func withUser(claims *jwt.MapClaims) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if claims != nil {
			c.Locals("user", claims)
		}
		return c.Next()
	}
}

// noClaim models a token minted before entitlements existed: the key is absent entirely.
func noClaim() *jwt.MapClaims {
	return &jwt.MapClaims{"user_id": "42", "role": UserTypes.LabScientist}
}

// claims models a token that carries the claim, including the empty claim a suspended
// facility mints.
func claims(modules ...string) *jwt.MapClaims {
	values := make([]any, 0, len(modules))
	for _, module := range modules {
		values = append(values, module)
	}

	return &jwt.MapClaims{"user_id": "42", "role": UserTypes.LabScientist, "modules": values}
}

func call(t *testing.T, user *jwt.MapClaims, module string) int {
	t.Helper()

	app := fiber.New()
	app.Get("/orders", withUser(user), RequireModule(module), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/orders", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode
}

func TestRequireModule(t *testing.T) {
	t.Cleanup(func() { SetModuleEnforcement(ModuleWarn) })

	t.Run("a held module passes", func(t *testing.T) {
		SetModuleEnforcement(ModuleEnforce)
		if got := call(t, claims(Modules.Lab), Modules.Lab); got != fiber.StatusOK {
			t.Fatalf("status = %d, want %d", got, fiber.StatusOK)
		}
	})

	t.Run("a module the facility did not license is refused", func(t *testing.T) {
		SetModuleEnforcement(ModuleEnforce)
		if got := call(t, claims(Modules.Radiology), Modules.Lab); got != fiber.StatusForbidden {
			t.Fatalf("status = %d, want %d", got, fiber.StatusForbidden)
		}
	})

	t.Run("an unauthenticated request is 401, not 403", func(t *testing.T) {
		// The distinction matters to a client deciding whether to re-authenticate.
		if got := call(t, nil, Modules.Lab); got != fiber.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", got, fiber.StatusUnauthorized)
		}
	})

	t.Run("a claimless token is admitted while warning", func(t *testing.T) {
		// Access tokens live 30 minutes. A service that rejected claimless tokens before
		// auth started minting them would lock every user out for up to half an hour.
		SetModuleEnforcement(ModuleWarn)
		if got := call(t, noClaim(), Modules.Lab); got != fiber.StatusOK {
			t.Fatalf("status = %d, want %d", got, fiber.StatusOK)
		}
	})

	t.Run("a claimless token is refused once enforcing", func(t *testing.T) {
		SetModuleEnforcement(ModuleEnforce)
		if got := call(t, noClaim(), Modules.Lab); got != fiber.StatusForbidden {
			t.Fatalf("status = %d, want %d", got, fiber.StatusForbidden)
		}
	})

	t.Run("an empty claim is a real answer, not a missing one", func(t *testing.T) {
		// A suspended facility mints "modules": []. That looks identical to a legacy token
		// under a length check and means the opposite, so warn mode must still refuse it —
		// otherwise suspension for non-payment hands back the product it withdrew.
		SetModuleEnforcement(ModuleWarn)
		if got := call(t, claims([]string{}...), Modules.Lab); got != fiber.StatusForbidden {
			t.Fatalf("status = %d, want %d", got, fiber.StatusForbidden)
		}
	})

	t.Run("a module is matched exactly, not by substring", func(t *testing.T) {
		// The package's Contains helper is a substring test. Entitlements must not be.
		SetModuleEnforcement(ModuleEnforce)
		if got := call(t, claims("notalab"), Modules.Lab); got != fiber.StatusForbidden {
			t.Fatalf("status = %d, want %d", got, fiber.StatusForbidden)
		}
	})
}

func TestParseModuleEnforcement(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want ModuleEnforcement
	}{
		{"enforce", ModuleEnforce},
		{" Enforce ", ModuleEnforce},
		{"warn", ModuleWarn},
		{"", ModuleWarn},
		{"nonsense", ModuleWarn},
	} {
		if got := ParseModuleEnforcement(tc.in); got != tc.want {
			t.Fatalf("ParseModuleEnforcement(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// The claim has to survive the trip through JSON that SerializeRequestUser makes.
func TestModulesClaimRoundTrip(t *testing.T) {
	encoded, err := json.Marshal(jwt.MapClaims{"modules": []string{Modules.Lab}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded struct {
		Modules []string `json:"modules"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(decoded.Modules) != 1 || decoded.Modules[0] != Modules.Lab {
		t.Fatalf("modules = %v, want [%s]", decoded.Modules, Modules.Lab)
	}
}
