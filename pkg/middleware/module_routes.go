package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/OtchereDev/ris-common-sdk/pkg/types"
	. "github.com/OtchereDev/ris-common-sdk/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ModuleRoute says which module a part of a service's API belongs to.
type ModuleRoute struct {
	// Path is matched against the start of the request path, segment by segment. A segment
	// written ":name" matches any one segment. "/reports" covers "/reports" and everything
	// under it, and not "/reports-archive".
	Path string

	// Modules are the modules that reach this path; holding any one of them is enough. A
	// path both staff and referring doctors use lists both: radiology and doctor_app.
	Modules []string

	// Access says whether the path reads or writes. Leave it at AccessByMethod unless
	// requests to it read with POST, or write with GET.
	Access Access
}

// ModuleRoutes guards a whole service from one table, mounted once with app.Use before the
// routes. Each request is checked against the rule with the longest matching path; a
// request no rule matches is core and passes unchecked.
//
// It runs before the routes' own authentication, so it reads the token itself, the same way
// AuthMiddleware does, and changes nothing on the request. A request with no token, or one
// that doesn't verify, goes on to the route, whose own authentication decides as it always
// has: a public link stays public, and a protected route still refuses it.
//
// The table is checked when the service starts. A module name this build doesn't know, a
// rule with no modules, or the same path twice stops it, because a typo here would silently
// leave part of the API unguarded.
func ModuleRoutes(routes ...ModuleRoute) fiber.Handler {
	rules := compileModuleRoutes(routes)

	return func(c *fiber.Ctx) error {
		rule := matchModuleRoute(rules, c.Path())
		if rule == nil {
			return c.Next()
		}

		user := requestUser(c)
		if user == nil {
			return c.Next()
		}

		if err := decide(c, user, rule.Modules, rule.Access); err != nil {
			return DenyModule(c, err)
		}
		return c.Next()
	}
}

type compiledRoute struct {
	ModuleRoute
	segments []string
	literals int
}

func compileModuleRoutes(routes []ModuleRoute) []compiledRoute {
	known := map[string]bool{}
	for _, m := range types.AllModules {
		known[m] = true
	}

	seen := map[string]bool{}
	out := make([]compiledRoute, 0, len(routes))
	for _, r := range routes {
		if !strings.HasPrefix(r.Path, "/") {
			panic(fmt.Sprintf("middleware.ModuleRoutes: path %q must start with /", r.Path))
		}
		if len(r.Modules) == 0 {
			panic(fmt.Sprintf("middleware.ModuleRoutes: %s lists no modules", r.Path))
		}
		for _, m := range r.Modules {
			if !known[m] {
				panic(fmt.Sprintf("middleware.ModuleRoutes: %s names %q, which isn't a module", r.Path, m))
			}
		}
		segs := splitPath(r.Path)
		key := strings.Join(segs, "/")
		if seen[key] {
			panic(fmt.Sprintf("middleware.ModuleRoutes: %s appears twice", r.Path))
		}
		seen[key] = true

		literals := 0
		for _, s := range segs {
			if !strings.HasPrefix(s, ":") {
				literals++
			}
		}
		out = append(out, compiledRoute{ModuleRoute: r, segments: segs, literals: literals})
	}
	return out
}

// matchModuleRoute is the rule whose path matches the most leading segments, preferring
// literal segments over parameters when two match as many.
func matchModuleRoute(rules []compiledRoute, path string) *compiledRoute {
	segs := splitPath(path)
	var best *compiledRoute
	for i := range rules {
		r := &rules[i]
		if len(r.segments) > len(segs) {
			continue
		}
		ok := true
		for j, want := range r.segments {
			if !strings.HasPrefix(want, ":") && !strings.EqualFold(want, segs[j]) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if best == nil || len(r.segments) > len(best.segments) ||
			(len(r.segments) == len(best.segments) && r.literals > best.literals) {
			best = r
		}
	}
	return best
}

func splitPath(path string) []string {
	var out []string
	for _, s := range strings.Split(path, "/") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// authDecode is the token decoder of the AuthMiddleware this service set up, so ModuleRoutes
// verifies tokens exactly as the routes will.
var authDecode atomic.Pointer[func(*fiber.Ctx) (*jwt.MapClaims, error)]

// requestUser is the person making the request: from the context if authentication has
// already run, otherwise by verifying the bearer token. Nil if there is no token, or it
// doesn't verify.
func requestUser(c *fiber.Ctx) *JwtPayload {
	if c.Locals("user") != nil {
		if u, err := SerializeRequestUser(c); err == nil && u != nil {
			return u
		}
	}

	decode := authDecode.Load()
	if decode == nil {
		return nil
	}
	// The decoder takes the second word of the header on trust; anything not shaped
	// "Bearer <token>" is left for the route's own authentication to turn away.
	parts := strings.Fields(c.Get(fiber.HeaderAuthorization))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return nil
	}
	claims, err := (*decode)(c)
	if err != nil || claims == nil {
		return nil
	}

	user, err := payloadFromClaims(claims)
	if err != nil {
		return nil
	}
	return user
}

// payloadFromClaims reads verified claims the way SerializeRequestUser reads them from the
// context, through JSON.
func payloadFromClaims(claims *jwt.MapClaims) (*JwtPayload, error) {
	raw, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	var u JwtPayload
	if err := json.Unmarshal(raw, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
