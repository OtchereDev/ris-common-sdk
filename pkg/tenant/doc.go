// Package tenant gives a service its tenant's record: who the tenant is, how it is branded,
// what it has licensed and what it may use.
//
// The record arrives as a file, tenant.json, written by the host's agent from what the
// control plane has published. A service never calls the control plane. It reads the file
// at boot, refuses to start without a valid one, and picks up a new one while running:
//
//	t, err := tenant.Open(ctx, tenant.OptionsFromEnv())
//	if err != nil {
//		log.Fatal(err) // no record, no service: never guess whose deployment this is
//	}
//	tenant.SetDefault(t)
//
//	t.Profile().Name                      // "Pluslab Diagnostics"
//	t.ModuleState(tenant.ModuleLab)       // enabled, read_only or off
//	t.Feature("whatsapp_reports")         // false unless the tenant has it
//	limit, limited := t.Limit("staff_seats")
//
// # Replacing the file
//
// The agent writes a new record to a temporary file in the same directory and renames it
// over tenant.json, so a reader never sees half a record. Mount the directory into the
// container, not the file: a bind-mounted file keeps pointing at the old inode after a
// rename, and the service would never see a change.
//
// A new file is accepted only if it is valid, names the same tenant, and carries a higher
// version than the one in use. Anything else is logged and ignored, and the service keeps
// running on the last good record. A stale entitlement is far less harmful than a
// radiology department that cannot start.
//
// # Two naming conventions
//
// The record's own fields are snake_case, like the rest of the control plane's API. The
// profile is the facility profile format that already exists (see facilities/README.md in
// the infra repository) and keeps its camelCase names, so the web app's branding code and
// the profile checks read it unchanged.
package tenant
