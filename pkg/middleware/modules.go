package middleware

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	. "github.com/OtchereDev/ris-common-sdk/pkg/types"
	. "github.com/OtchereDev/ris-common-sdk/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// ---------------------------------------------------------------------------
// Module entitlements
//
// A licensed product is enforced here, at each service's API boundary, and not in
// navigation. Hiding a menu item does not stop anyone calling the service directly, and a
// build-time flag cannot be turned off without a rebuild - which means it is not
// enforcement at all, only decoration.
//
// A request may reach a module when two things hold:
//
//	the tenant has it    the module's state in tenant.json (pkg/tenant), which the host's
//	                     agent keeps current: enabled, or read_only for a request that reads
//	the person holds it  what the role in the token grants (types.GrantedModules)
//
// Both are read live, so a module the control plane switches off is refused within seconds,
// and one it switches on works without waiting for a new token.
//
// A service that hasn't loaded a tenant record yet keeps the earlier rule: the modules
// claim auth minted into the token.

// ModuleEnforcement decides whether a refusal is applied or only reported.
type ModuleEnforcement int

const (
	// ModuleWarn reports what the tenant-record check would refuse, and doesn't apply it.
	// Each would-be refusal is logged with the tenant, role, route, module and reason, so a
	// wrong grant or a write route that should be a read is found before it blocks anyone.
	//
	// Warn changes nothing about who gets in. A route behaves as it did before it was
	// guarded: for lab, the one module the token's claim has always been enforced on, a
	// token whose claim lacks it is still refused, which is what keeps a suspended
	// facility's lab from coming back; for every other module there was no check, so
	// there still isn't one, only the report.
	//
	// It is the rollout setting, not a destination.
	ModuleWarn ModuleEnforcement = iota

	// ModuleEnforce applies the tenant-record check. Switch to it after reviewing a week of
	// warn logs for the service.
	ModuleEnforce
)

var moduleEnforcement = ModuleWarn

// SetModuleEnforcement is called once at boot, from the service's config.
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

// Access is whether a route reads or writes, which decides what a read_only module allows.
type Access int

const (
	// AccessByMethod treats GET, HEAD and OPTIONS as reads and everything else as writes.
	AccessByMethod Access = iota
	// Read is for a route that only reads but uses another method, such as a search or an
	// export sent as POST. Guard it with RequireModuleFor(module, Read) on the route itself,
	// outside any group guarded by RequireModule: a group's guard runs before the route's
	// own handlers, so it can't see what the route declares.
	Read
	// Write is for a route that changes something but uses GET.
	Write
)

func (a Access) reads(method string) bool {
	switch a {
	case Read:
		return true
	case Write:
		return false
	}
	switch method {
	case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
		return true
	}
	return false
}

// ErrUnauthenticated means no token reached the guard at all. Distinct from a refusal,
// because it tells a client to re-authenticate rather than to give up.
var ErrUnauthenticated = errors.New("unauthenticated")

// ErrModuleNotLicensed means the tenant doesn't have the module.
var ErrModuleNotLicensed = errors.New("module not licensed")

// ErrModuleReadOnly means the tenant may read the module's records but not add to them.
var ErrModuleReadOnly = errors.New("module is read-only")

// ErrModuleNotGranted means the tenant has the module but the person's role doesn't
// include it.
var ErrModuleNotGranted = errors.New("module not granted to role")

var ModuleNotLicensedError = &fiber.Map{
	"status": fiber.StatusForbidden,
	"data": &fiber.Map{
		"message": "This module is not enabled for your facility"},
}

var ModuleReadOnlyError = &fiber.Map{
	"status": fiber.StatusForbidden,
	"data": &fiber.Map{
		"message": "This module is read-only for your facility. You can view and export its records, but not add to them."},
}

var ModuleNotGrantedError = &fiber.Map{
	"status": fiber.StatusForbidden,
	"data": &fiber.Map{
		"message": "Your role doesn't include this part of the system"},
}

// CheckModule decides whether the request may reach a module, reading or writing by its
// HTTP method, and writes nothing. See CheckModuleFor.
func CheckModule(c *fiber.Ctx, module string) error {
	return CheckModuleFor(c, module, AccessByMethod)
}

// CheckModuleFor decides whether the request may reach a module, and writes nothing.
//
// It exists separately from RequireModuleFor because a service whose routes already run
// through one authorization choke point wants to add this check there, and composing two
// middlewares that each call c.Next() would run the rest of the chain twice.
//
// Who the person is comes from the signed token and from nowhere else. The edge sets an
// X-Tenant-Key header for logging, and a service must never take an entitlement from a
// header - a header is only as trustworthy as everything between the client and here.
func CheckModuleFor(c *fiber.Ctx, module string, access Access) error {
	user, err := SerializeRequestUser(c)
	if err != nil || user == nil {
		return ErrUnauthenticated
	}
	return decide(c, user, []string{module}, access)
}

// decide is the rule for a request that may reach any one of modules.
func decide(c *fiber.Ctx, user *JwtPayload, modules []string, access Access) error {
	reading := access.reads(c.Method())
	label := strings.Join(modules, "|")

	byClaim := ErrModuleNotLicensed
	enforcedByClaim := false
	for _, m := range modules {
		if claimCheck(user, m) == nil {
			byClaim = nil
		}
		enforcedByClaim = enforcedByClaim || claimEnforced[m]
	}
	// A token with no claim at all predates entitlements. Where the claim still decides,
	// warn mode lets it through, and says so.
	if byClaim == nil && !user.Modules.Present() && enforcedByClaim {
		log.Printf("module %q: admitting user %s on a token with no modules claim; set "+
			"module enforcement to enforce once auth is deployed", label, user.UserID)
	}

	t := tenant.Default()
	if t == nil {
		// No tenant record yet: the token's claim is the only rule there is.
		if moduleEnforcement == ModuleEnforce || enforcedByClaim {
			return byClaim
		}
		if byClaim != nil {
			reportWouldRefuse(c, nil, user, label, "", reading, "token's modules claim lacks the module")
		}
		return nil
	}

	refusal, state, reason := tenantCheckAny(t, user, modules, reading)

	if moduleEnforcement == ModuleEnforce {
		return refusal
	}
	if refusal != nil {
		reportWouldRefuse(c, t, user, label, state, reading, reason)
	}
	if enforcedByClaim {
		return byClaim
	}
	return nil
}

// tenantCheckAny allows the request if any one of modules allows it. If none does, it
// reports the refusal that says most: a module the person holds but the tenant doesn't
// allow, before one their role doesn't include at all.
func tenantCheckAny(t *tenant.Tenant, user *JwtPayload, modules []string, reading bool) (error, tenant.ModuleState, string) {
	var (
		refusal error
		state   tenant.ModuleState
		reason  string
	)
	for _, m := range modules {
		err, s, r := tenantCheck(t, user, m, reading)
		if err == nil {
			return nil, s, ""
		}
		if refusal == nil || (errors.Is(refusal, ErrModuleNotGranted) && !errors.Is(err, ErrModuleNotGranted)) {
			refusal, state, reason = err, s, r
		}
	}
	return refusal, state, reason
}

// claimEnforced are the modules the token's claim was enforced on before tenant records:
// lab, by the lab service. In warn mode those keep that rule, and every other module keeps
// having none, so turning a service's guards on in warn mode can't lock anyone out.
var claimEnforced = map[string]bool{Modules.Lab: true}

// tenantCheck is the rule: the role grants the module, and the tenant has it in a state
// that allows this kind of request.
func tenantCheck(t *tenant.Tenant, user *JwtPayload, module string, reading bool) (error, tenant.ModuleState, string) {
	state := t.ModuleState(module)

	grant := Grant{Role: user.UserType, IsAdmin: user.IsAdmin, CanViewAnalytics: user.CanViewAnalytics}
	if !HasGrant(grant, module) {
		return ErrModuleNotGranted, state, "role does not include the module"
	}

	switch state {
	case tenant.ModuleEnabled:
		return nil, state, ""
	case tenant.ModuleReadOnly:
		if reading {
			return nil, state, ""
		}
		return ErrModuleReadOnly, state, "module is read-only and the request writes"
	default:
		return ErrModuleNotLicensed, state, "tenant does not have the module"
	}
}

// claimCheck is the rule before tenant records: the modules claim in the token.
func claimCheck(user *JwtPayload, module string) error {
	if user.HasModule(module) {
		return nil
	}

	// No claim at all is a token minted before entitlements existed. That is a different
	// thing from a claim of nothing, which is what a suspended facility mints and which is
	// refused in either mode.
	if !user.Modules.Present() && moduleEnforcement == ModuleWarn {
		return nil
	}

	return ErrModuleNotLicensed
}

// moduleLog is where would-be refusals are reported. A variable so tests can capture it.
var moduleLog = slog.Default

// wouldRefuse limits each distinct would-be refusal to one log line a minute, so a whole
// route group tagged wrongly doesn't bury everything else in the log.
var wouldRefuse = struct {
	sync.Mutex
	last       map[string]time.Time
	suppressed map[string]int
}{last: map[string]time.Time{}, suppressed: map[string]int{}}

const wouldRefuseEvery = time.Minute

// reportWouldRefuse logs a refusal warn mode didn't apply. t is nil when the service has no
// tenant record and the refusal comes from the token's claim.
func reportWouldRefuse(c *fiber.Ctx, t *tenant.Tenant, user *JwtPayload, module string, state tenant.ModuleState, reading bool, reason string) {
	path := redactPath(c.Path())
	key := fmt.Sprintf("%s|%s|%s|%s|%s", module, reason, strings.ToUpper(user.UserType), c.Method(), path)

	wouldRefuse.Lock()
	now := time.Now()
	if last, ok := wouldRefuse.last[key]; ok && now.Sub(last) < wouldRefuseEvery {
		wouldRefuse.suppressed[key]++
		wouldRefuse.Unlock()
		return
	}
	suppressed := wouldRefuse.suppressed[key]
	wouldRefuse.last[key] = now
	wouldRefuse.suppressed[key] = 0
	wouldRefuse.Unlock()

	access := "write"
	if reading {
		access = "read"
	}
	tenantKey, version := user.TenantKey, int64(0)
	if t != nil {
		tenantKey, version = t.Key(), t.Version()
	}
	moduleLog().Warn("module check would refuse",
		"mode", "warn",
		"tenant", tenantKey,
		"tenant_version", version,
		"module", module,
		"state", string(state),
		"access", access,
		"reason", reason,
		"role", user.UserType,
		"is_admin", user.IsAdmin,
		"user_id", user.UserID,
		"method", c.Method(),
		"path", path,
		"suppressed_since_last", suppressed,
	)
}

var (
	digits   = regexp.MustCompile(`^\d+$`)
	uuidLike = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	hexLong  = regexp.MustCompile(`^[0-9a-fA-F]{16,}$`)
	hasDigit = regexp.MustCompile(`\d`)
)

// redactPath replaces path segments that look like identifiers with ":id", so the log names
// routes rather than patients, studies or receipts.
func redactPath(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if digits.MatchString(p) || uuidLike.MatchString(p) || hexLong.MatchString(p) ||
			(len(p) >= 6 && hasDigit.MatchString(p)) {
			parts[i] = ":id"
		}
	}
	return strings.Join(parts, "/")
}

// DenyModule writes the refusal CheckModule decided on.
func DenyModule(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUnauthenticated):
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"status": fiber.StatusUnauthorized,
			"data":   &fiber.Map{"message": "Unauthorized"}})
	case errors.Is(err, ErrModuleReadOnly):
		return c.Status(fiber.StatusForbidden).JSON(ModuleReadOnlyError)
	case errors.Is(err, ErrModuleNotGranted):
		return c.Status(fiber.StatusForbidden).JSON(ModuleNotGrantedError)
	}
	return c.Status(fiber.StatusForbidden).JSON(ModuleNotLicensedError)
}

// RequireModule guards a route or group on a module, reading or writing by HTTP method.
// Mount it after AuthMiddleware, which is what puts the token on the context.
func RequireModule(module string) fiber.Handler {
	return RequireModuleFor(module, AccessByMethod)
}

// RequireModuleFor guards a route on a module, saying whether it reads or writes.
func RequireModuleFor(module string, access Access) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := CheckModuleFor(c, module, access); err != nil {
			return DenyModule(c, err)
		}

		return c.Next()
	}
}
