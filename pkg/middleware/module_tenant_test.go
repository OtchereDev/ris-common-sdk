package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
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

// staff is a token for a staff member with the given role, carrying a modules claim. A nil
// modules slice leaves the claim out, as a token minted before entitlements did.
func staff(role string, modules []string, extra ...func(jwt.MapClaims)) *jwt.MapClaims {
	c := jwt.MapClaims{"user_id": "7", "role": role}
	if modules != nil {
		values := make([]any, 0, len(modules))
		for _, m := range modules {
			values = append(values, m)
		}
		c["modules"] = values
	}
	for _, e := range extra {
		e(c)
	}
	return &c
}

func admin(c jwt.MapClaims) { c["is_admin"] = true }

// withTenant loads a tenant record for the test, and the given enforcement mode.
func withTenant(t *testing.T, mode ModuleEnforcement, changes ...func(*tenant.Record)) {
	t.Helper()
	tenant.SetDefault(tenanttest.New(t, changes...))
	SetModuleEnforcement(mode)
	t.Cleanup(func() {
		tenant.SetDefault(nil)
		SetModuleEnforcement(ModuleWarn)
		wouldRefuse.Lock()
		clear(wouldRefuse.last)
		clear(wouldRefuse.suppressed)
		wouldRefuse.Unlock()
	})
}

// captureLog collects what would-be refusals log.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	moduleLog = func() *slog.Logger { return logger }
	t.Cleanup(func() { moduleLog = slog.Default })
	return &buf
}

type result struct {
	status  int
	message string
}

// send makes one request through a guard, and returns the status and refusal message.
func send(t *testing.T, user *jwt.MapClaims, method, path string, guard fiber.Handler) result {
	t.Helper()
	app := fiber.New()
	app.Add(method, "/*", withUser(user), guard, func(c *fiber.Ctx) error { return c.SendString("ok") })

	resp, err := app.Test(httptest.NewRequest(method, path, nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &parsed)
	return result{resp.StatusCode, parsed.Data.Message}
}

func TestModuleCheckWithATenantRecord(t *testing.T) {
	rad := Modules.Radiology
	lab := Modules.Lab
	cases := []struct {
		name    string
		tenant  []func(*tenant.Record)
		user    *jwt.MapClaims
		method  string
		guard   fiber.Handler
		status  int
		message string
	}{
		{"a radiographer reaches radiology", nil,
			staff(UserTypes.Radiographer, []string{rad}), "POST", RequireModule(rad), 200, ""},
		{"a module the tenant doesn't have", nil,
			staff(UserTypes.LabScientist, []string{lab}), "GET", RequireModule(lab), 403, "not enabled for your facility"},
		{"a module the role doesn't include", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleEnabled)},
			staff(UserTypes.Radiographer, []string{rad}), "GET", RequireModule(lab), 403, "role doesn't include"},
		{"an admin reaches any module the tenant has", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleEnabled)},
			staff(UserTypes.Cashier, []string{rad}, admin), "POST", RequireModule(lab), 200, ""},
		{"an admin can't reach one it doesn't", nil,
			staff(UserTypes.Admin, []string{rad}, admin), "GET", RequireModule(lab), 403, "not enabled for your facility"},

		{"read_only: a GET reads", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleReadOnly)},
			staff(UserTypes.LabScientist, nil), "GET", RequireModule(lab), 200, ""},
		{"read_only: a POST writes", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleReadOnly)},
			staff(UserTypes.LabScientist, nil), "POST", RequireModule(lab), 403, "read-only"},
		{"read_only: a POST declared a read", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleReadOnly)},
			staff(UserTypes.LabScientist, nil), "POST", RequireModuleFor(lab, Read), 200, ""},
		{"a GET declared a write", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleReadOnly)},
			staff(UserTypes.LabScientist, nil), "GET", RequireModuleFor(lab, Write), 403, "read-only"},

		{"a read_only tenant can read", []func(*tenant.Record){tenanttest.WithStatus(tenant.StatusReadOnly)},
			staff(UserTypes.Radiologist, []string{rad}), "GET", RequireModule(rad), 200, ""},
		{"a read_only tenant can't write", []func(*tenant.Record){tenanttest.WithStatus(tenant.StatusReadOnly)},
			staff(UserTypes.Radiologist, []string{rad}), "PATCH", RequireModule(rad), 403, "read-only"},
		{"a restricted tenant still does clinical work", []func(*tenant.Record){tenanttest.WithStatus(tenant.StatusRestricted)},
			staff(UserTypes.Radiologist, []string{rad}), "POST", RequireModule(rad), 200, ""},
		{"a terminated tenant reaches nothing", []func(*tenant.Record){tenanttest.WithStatus(tenant.StatusTerminated)},
			staff(UserTypes.Radiologist, []string{rad}), "GET", RequireModule(rad), 403, "not enabled"},

		// The tenant record is live: a module added after the token was minted works at once.
		{"a token minted before the module was added", []func(*tenant.Record){tenanttest.WithModule(lab, tenant.ModuleEnabled)},
			staff(UserTypes.LabScientist, []string{}), "GET", RequireModule(lab), 200, ""},

		{"a referring doctor reaches the doctor app", []func(*tenant.Record){tenanttest.WithModule(Modules.DoctorApp, tenant.ModuleEnabled)},
			staff(RoleDoctor, nil), "GET", RequireModule(Modules.DoctorApp), 200, ""},
		{"a patient can't reach the doctor app", []func(*tenant.Record){tenanttest.WithModule(Modules.DoctorApp, tenant.ModuleEnabled)},
			staff(RolePatient, nil), "GET", RequireModule(Modules.DoctorApp), 403, "role doesn't include"},
		{"analytics access reaches ask", []func(*tenant.Record){tenanttest.WithModule(Modules.Ask, tenant.ModuleEnabled)},
			staff(UserTypes.Cashier, nil, func(c jwt.MapClaims) { c["can_view_analytics"] = true }), "POST", RequireModule(Modules.Ask), 200, ""},

		{"no token at all is 401", nil, nil, "GET", RequireModule(rad), 401, "Unauthorized"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withTenant(t, ModuleEnforce, c.tenant...)
			got := send(t, c.user, c.method, "/x", c.guard)
			if got.status != c.status {
				t.Fatalf("status = %d, want %d (%q)", got.status, c.status, got.message)
			}
			if c.message != "" && !strings.Contains(got.message, c.message) {
				t.Fatalf("message = %q, want it to contain %q", got.message, c.message)
			}
		})
	}
}

func TestWarnModeReportsButNeverLoosens(t *testing.T) {
	lab := Modules.Lab

	t.Run("a refusal the new check would make is logged and not applied", func(t *testing.T) {
		withTenant(t, ModuleWarn) // the tenant has no lab
		logged := captureLog(t)

		got := send(t, staff(UserTypes.LabScientist, []string{lab}), "POST", "/lab/orders/123/results", RequireModule(lab))
		if got.status != 200 {
			t.Fatalf("warn mode applied a refusal: %d", got.status)
		}

		var line map[string]any
		if err := json.Unmarshal(logged.Bytes(), &line); err != nil {
			t.Fatalf("want one JSON log line, got %q", logged.String())
		}
		want := map[string]any{
			"msg": "module check would refuse", "tenant": "demo", "module": "lab", "state": "off",
			"access": "write", "reason": "tenant does not have the module", "role": "LABSCIENTIST",
			"method": "POST", "path": "/lab/orders/:id/results",
		}
		for k, v := range want {
			if line[k] != v {
				t.Errorf("log %s = %v, want %v", k, line[k], v)
			}
		}
	})

	t.Run("what the token's claim refuses is still refused", func(t *testing.T) {
		// A suspended facility's tokens carry an empty claim. Warn mode must not hand the
		// product back just because the new check is only reporting.
		withTenant(t, ModuleWarn, tenanttest.WithModule(lab, tenant.ModuleEnabled))
		captureLog(t)

		got := send(t, staff(UserTypes.LabScientist, []string{}), "GET", "/lab", RequireModule(lab))
		if got.status != 403 {
			t.Fatalf("status = %d, want 403", got.status)
		}
	})

	t.Run("a newly guarded module locks nobody out", func(t *testing.T) {
		// Today's tokens only ever carry radiology, and nothing checked analytics or the
		// doctor app before. Guarding those routes in warn mode must change nothing.
		withTenant(t, ModuleWarn) // the tenant has no analytics
		logged := captureLog(t)

		got := send(t, staff(UserTypes.Cashier, []string{Modules.Radiology}), "POST", "/analytics/query", RequireModule(Modules.Analytics))
		if got.status != 200 {
			t.Fatalf("status = %d, want 200", got.status)
		}
		if !strings.Contains(logged.String(), "would refuse") {
			t.Fatal("the would-be refusal wasn't reported")
		}
	})

	t.Run("a request both checks allow logs nothing", func(t *testing.T) {
		withTenant(t, ModuleWarn, tenanttest.WithModule(lab, tenant.ModuleEnabled))
		logged := captureLog(t)

		if got := send(t, staff(UserTypes.LabScientist, []string{lab}), "GET", "/lab", RequireModule(lab)); got.status != 200 {
			t.Fatalf("status = %d", got.status)
		}
		if logged.Len() != 0 {
			t.Fatalf("logged %q", logged.String())
		}
	})

	t.Run("the same would-be refusal is logged once a minute", func(t *testing.T) {
		withTenant(t, ModuleWarn)
		logged := captureLog(t)

		for range 5 {
			send(t, staff(UserTypes.LabScientist, []string{lab}), "GET", "/lab/orders/1", RequireModule(lab))
			send(t, staff(UserTypes.LabScientist, []string{lab}), "GET", "/lab/orders/2", RequireModule(lab))
		}
		if n := strings.Count(logged.String(), "would refuse"); n != 1 {
			t.Fatalf("logged %d lines for one route, want 1", n)
		}

		// After the minute, the next one is logged with the count it stood for.
		wouldRefuse.Lock()
		for k := range wouldRefuse.last {
			wouldRefuse.last[k] = time.Now().Add(-2 * wouldRefuseEvery)
		}
		wouldRefuse.Unlock()
		send(t, staff(UserTypes.LabScientist, []string{lab}), "GET", "/lab/orders/3", RequireModule(lab))
		if !strings.Contains(logged.String(), `"suppressed_since_last":9`) {
			t.Fatalf("want the suppressed count of 9 reported, got %q", logged.String())
		}
	})
}

func TestRedactPath(t *testing.T) {
	for in, want := range map[string]string{
		"/lab/orders/123/results":                           "/lab/orders/:id/results",
		"/appointments/SPHN204512/report":                   "/appointments/:id/report",
		"/share/9f8c1e2a-3b4d-4c5e-8f60-718293a4b5c6":       "/share/:id",
		"/studies/1f8b3c9d2e7a6b5c":                         "/studies/:id",
		"/api/v1/analytics/boards":                          "/api/v1/analytics/boards",
		"/lab/specimen-types":                               "/lab/specimen-types",
		"/doctors/4/patients/88/appointments?status=active": "/doctors/:id/patients/:id/appointments?status=active",
	} {
		if got := redactPath(in); got != want {
			t.Errorf("redactPath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestWarnModeWithoutATenantRecord(t *testing.T) {
	t.Cleanup(func() { SetModuleEnforcement(ModuleWarn) })

	t.Run("a doctor reaching a newly guarded radiology route is admitted and reported", func(t *testing.T) {
		SetModuleEnforcement(ModuleWarn)
		logged := captureLog(t)
		t.Cleanup(func() {
			wouldRefuse.Lock()
			clear(wouldRefuse.last)
			wouldRefuse.Unlock()
		})

		got := send(t, staff(RoleDoctor, []string{}), "GET", "/doctor/reports", RequireModule(Modules.Radiology))
		if got.status != 200 {
			t.Fatalf("status = %d, want 200", got.status)
		}
		if !strings.Contains(logged.String(), "claim lacks the module") {
			t.Fatalf("want a claim-based report, got %q", logged.String())
		}
	})

	t.Run("enforce mode without a record applies the claim to every module", func(t *testing.T) {
		SetModuleEnforcement(ModuleEnforce)
		got := send(t, staff(RoleDoctor, []string{}), "GET", "/doctor/reports", RequireModule(Modules.Radiology))
		if got.status != 403 {
			t.Fatalf("status = %d, want 403", got.status)
		}
	})
}
