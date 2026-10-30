package middleware

import (
	"net/http/httptest"
	"testing"

	. "github.com/OtchereDev/ris-common-sdk/pkg/types"
	"github.com/gofiber/fiber/v2"
)

func TestRolesOrAdmin(t *testing.T) {
	cases := []struct {
		name   string
		claims map[string]interface{}
		want   int
	}{
		{"marketer", map[string]interface{}{"role": "MARKETER"}, 200},
		{"accountant", map[string]interface{}{"role": "ACCOUNTANT"}, 200},
		{"admin by flag", map[string]interface{}{"role": "RADIOLOGIST", "is_admin": true}, 200},
		{"cashier", map[string]interface{}{"role": "CASHIER"}, 403},
		{"nobody", nil, 403},
	}
	for _, c := range cases {
		app := fiber.New()
		app.Get("/", func(ctx *fiber.Ctx) error {
			if c.claims != nil {
				ctx.Locals("user", c.claims)
			}
			return ctx.Next()
		}, RolesOrAdmin(UserTypes.Admin, UserTypes.Accountant, UserTypes.Marketer), func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(200)
		})
		res, err := app.Test(httptest.NewRequest("GET", "/", nil))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != c.want {
			t.Errorf("%s: got %d, want %d", c.name, res.StatusCode, c.want)
		}
	}
}

func TestMarketerHoldsRadiology(t *testing.T) {
	got := ModulesForRole(UserTypes.Marketer, false)
	if len(got) != 1 || got[0] != Modules.Radiology {
		t.Errorf("a marketer's modules = %v, want radiology", got)
	}
}
