package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant/tenanttest"
	. "github.com/OtchereDev/ris-common-sdk/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const routesSecret = "module-routes-test"

// sign makes a bearer token the configured AuthMiddleware accepts.
func sign(t *testing.T, role string, extra ...func(jwt.MapClaims)) string {
	t.Helper()
	c := jwt.MapClaims{"user_id": "9", "role": role, "exp": float64(time.Now().Add(time.Hour).Unix())}
	for _, e := range extra {
		e(c)
	}
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(routesSecret))
	if err != nil {
		t.Fatal(err)
	}
	return "Bearer " + s
}

// table is a small service: staff radiology, a path shared with referring doctors, a
// doctor-only path, a lab section with a search that reads by POST, and a public link.
var table = []ModuleRoute{
	{Path: "/appointment", Modules: []string{Modules.Radiology}},
	{Path: "/tmp-appointment", Modules: []string{Modules.Radiology, Modules.DoctorApp}},
	{Path: "/tmp-appointment/shared-with-me", Modules: []string{Modules.DoctorApp}},
	{Path: "/lab", Modules: []string{Modules.Lab}},
	{Path: "/lab/results/search", Modules: []string{Modules.Lab}, Access: Read},
	{Path: "/reports/:id/messages", Modules: []string{Modules.Radiology}},
}

// service mounts the table, then routes that authenticate themselves the way the
// appointment service's do, plus one public route.
func service(t *testing.T) *fiber.App {
	t.Helper()
	SetupMiddleAuthMiddleware(routesSecret, nil)
	t.Cleanup(func() { authDecode.Store(nil) })

	app := fiber.New()
	app.Use(ModuleRoutes(table...))
	ok := func(c *fiber.Ctx) error { return c.SendString("ok") }
	for _, p := range []string{"/appointment", "/tmp-appointment/appointments/:id", "/tmp-appointment/shared-with-me",
		"/lab/orders", "/lab/results/search", "/reports/:id/messages", "/patients"} {
		app.Get(p, AuthMiddleware, ok)
		app.Post(p, AuthMiddleware, ok)
	}
	app.Get("/share/:token", ok) // public
	app.Get("/tmp-appointment/public-form", ok)
	return app
}

func hit(t *testing.T, app *fiber.App, method, path, bearer string) int {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", bearer)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode
}

func TestModuleRoutesEnforcing(t *testing.T) {
	withTenant(t, ModuleEnforce,
		tenanttest.WithModule(Modules.DoctorApp, tenant.ModuleEnabled),
		tenanttest.WithModule(Modules.Lab, tenant.ModuleReadOnly))
	app := service(t)

	radiographer := sign(t, UserTypes.Radiographer)
	doctor := sign(t, RoleDoctor)
	scientist := sign(t, UserTypes.LabScientist)

	cases := []struct {
		name, method, path, bearer string
		want                       int
	}{
		{"staff reach radiology", "POST", "/appointment", radiographer, 200},
		{"a doctor can't reach staff radiology", "GET", "/appointment", doctor, 403},
		{"a doctor reaches a shared path through the doctor app", "GET", "/tmp-appointment/appointments/4", doctor, 200},
		{"staff reach the same shared path through radiology", "GET", "/tmp-appointment/appointments/4", radiographer, 200},
		{"the longer rule wins: staff can't reach a doctor-only path", "GET", "/tmp-appointment/shared-with-me", radiographer, 403},
		{"a read_only module reads", "GET", "/lab/orders", scientist, 200},
		{"a read_only module doesn't write", "POST", "/lab/orders", scientist, 403},
		{"a search declared a read, under a read_only module", "POST", "/lab/results/search", scientist, 200},
		{"a parameter segment matches", "POST", "/reports/77/messages", radiographer, 200},
		{"a parameter segment matches, and refuses", "POST", "/reports/77/messages", doctor, 403},
		{"no rule matches: core, unchecked", "GET", "/patients", doctor, 200},
		{"a prefix is matched by segment, not by text", "GET", "/appointments-archive", doctor, 404},

		// The guard never decides about a request without a verified token; the route does.
		{"a public route stays public", "GET", "/share/abc123", "", 200},
		{"a public route under a guarded prefix stays public", "GET", "/tmp-appointment/public-form", "", 200},
		{"a protected route still refuses no token", "GET", "/appointment", "", 401},
		{"a protected route still refuses a bad token", "GET", "/appointment", "Bearer not-a-token", 401},
		{"a malformed header doesn't crash anything", "GET", "/share/abc123", "Bearer", 200},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := hit(t, app, c.method, c.path, c.bearer); got != c.want {
				t.Fatalf("%s %s: %d, want %d", c.method, c.path, got, c.want)
			}
		})
	}
}

func TestModuleRoutesWarnModeLocksNobodyOut(t *testing.T) {
	withTenant(t, ModuleWarn) // no doctor app, no lab
	logged := captureLog(t)
	app := service(t)

	for _, c := range []struct{ path, bearer string }{
		{"/tmp-appointment/shared-with-me", sign(t, RoleDoctor)},
		{"/appointment", sign(t, RoleDoctor)},
	} {
		if got := hit(t, app, "GET", c.path, c.bearer); got != 200 {
			t.Fatalf("warn mode refused %s: %d", c.path, got)
		}
	}
	if n := strings.Count(logged.String(), "would refuse"); n != 2 {
		t.Fatalf("want 2 reports, got %d: %s", n, logged.String())
	}
	if !strings.Contains(logged.String(), `"module":"doctor_app"`) {
		t.Fatalf("report should name the module: %s", logged.String())
	}
}

func TestModuleRoutesPrefersTheMostTellingRefusal(t *testing.T) {
	// A lab scientist at a tenant whose lab is read_only, writing to a path lab and
	// radiology share: "read-only" says more than "your role doesn't include radiology".
	withTenant(t, ModuleEnforce, tenanttest.WithModule(Modules.Lab, tenant.ModuleReadOnly))
	SetupMiddleAuthMiddleware(routesSecret, nil)
	t.Cleanup(func() { authDecode.Store(nil) })

	app := fiber.New()
	app.Use(ModuleRoutes(ModuleRoute{Path: "/shared", Modules: []string{Modules.Radiology, Modules.Lab}}))
	app.Post("/shared", func(c *fiber.Ctx) error { return c.SendString("ok") })

	req := httptest.NewRequest("POST", "/shared", nil)
	req.Header.Set("Authorization", sign(t, UserTypes.LabScientist))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 403 || !strings.Contains(string(body), "read-only") {
		t.Fatalf("got %d %s", resp.StatusCode, body)
	}
}

func TestModuleRoutesRefusesABadTable(t *testing.T) {
	for name, routes := range map[string][]ModuleRoute{
		"an unknown module": {{Path: "/x", Modules: []string{"pathology"}}},
		"no modules":        {{Path: "/x"}},
		"a relative path":   {{Path: "x", Modules: []string{Modules.Lab}}},
		"a path twice":      {{Path: "/x", Modules: []string{Modules.Lab}}, {Path: "/x/", Modules: []string{Modules.Radiology}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("want the service to refuse to start")
				}
			}()
			ModuleRoutes(routes...)
		})
	}
}
