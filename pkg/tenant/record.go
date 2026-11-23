package tenant

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"
)

// SchemaVersion is the record format this build reads. Fields may be added without
// changing it; a service ignores fields it doesn't know. It changes only when an existing
// field changes meaning, and a service refuses a record whose schema it doesn't read.
const SchemaVersion = 1

// Record is one tenant's published record, as delivered in tenant.json.
type Record struct {
	Schema      int       `json:"schema"`
	Key         string    `json:"key"`
	Version     int64     `json:"version"`
	PublishedAt time.Time `json:"published_at"`

	// DisplayName is the customer's name for consoles, contracts and invoices. What the
	// product shows patients and staff is Profile.Name.
	DisplayName string `json:"display_name"`

	Status Status `json:"status"`

	// Domains are the hostnames the tenant answers on; the first is canonical.
	Domains []string `json:"domains,omitempty"`

	// Modules holds a state for each module. A module that isn't listed is off.
	Modules map[string]ModuleState `json:"modules"`

	// Features are the feature flags turned on for this tenant. Absent means off.
	Features []string `json:"features,omitempty"`

	// Limits are the plan's limits. A limit that isn't listed is unlimited.
	Limits map[string]int64 `json:"limits,omitempty"`

	// UsageMTD is the control plane's month-to-date usage, written once a day, so a
	// service whose local counter was lost doesn't under-count by more than a day.
	UsageMTD *UsageMTD `json:"usage_mtd,omitempty"`

	Profile Profile `json:"profile"`
}

// UsageMTD is month-to-date usage per meter, as of a day.
type UsageMTD struct {
	AsOf   string           `json:"as_of"` // YYYY-MM-DD
	Meters map[string]int64 `json:"meters"`
}

// Status is where the tenant stands commercially.
type Status string

const (
	StatusTrial      Status = "trial"
	StatusActive     Status = "active"
	StatusGrace      Status = "grace"
	StatusRestricted Status = "restricted"
	StatusReadOnly   Status = "read_only"
	StatusTerminated Status = "terminated"
)

func (s Status) valid() bool {
	switch s {
	case StatusTrial, StatusActive, StatusGrace, StatusRestricted, StatusReadOnly, StatusTerminated:
		return true
	}
	return false
}

// ModuleState is what a tenant may do with a module.
type ModuleState string

const (
	// ModuleEnabled allows new work and reading existing records.
	ModuleEnabled ModuleState = "enabled"
	// ModuleReadOnly allows reading and exporting existing records, and no new work. A
	// clinical module the tenant no longer has stays here, never off, because the records
	// it holds are part of patients' medical records.
	ModuleReadOnly ModuleState = "read_only"
	// ModuleOff allows nothing.
	ModuleOff ModuleState = "off"
)

func (s ModuleState) valid() bool {
	return s == ModuleEnabled || s == ModuleReadOnly || s == ModuleOff
}

// Module names this package refers to. The full catalogue lives in pkg/types; a record may
// carry modules this build doesn't know, and they are kept rather than refused, so a newer
// control plane never stops an older service from starting.
const (
	ModuleRadiology = "radiology"
	ModuleLab       = "lab"
)

// clinicalModules hold medical records, so they may never be set to off once used.
var clinicalModules = []string{ModuleRadiology, ModuleLab}

// Profile is the facility profile: brand, contacts, channels, identifiers and printed
// documents. The format is the one in facilities/README.md, with its camelCase names.
type Profile struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Site string `json:"site,omitempty"`

	Primary string   `json:"primary,omitempty"`
	Accent  string   `json:"accent,omitempty"`
	Radius  *float64 `json:"radius,omitempty"`
	Logo    Logo     `json:"logo"`

	Locale          Locale    `json:"locale"`
	Contact         Contact   `json:"contact"`
	Channels        Channels  `json:"channels"`
	Mobile          Mobile    `json:"mobile"`
	Links           Links     `json:"links"`
	PatientIDPrefix string    `json:"patientIdPrefix,omitempty"`
	Radiology       Radiology `json:"radiology"`
	DICOM           DICOM     `json:"dicom"`
	Document        Document  `json:"document"`
}

// Logo holds artwork as a data: URI or bare base64.
type Logo struct {
	Full    string `json:"full,omitempty"`
	Mark    string `json:"mark,omitempty"`
	Favicon string `json:"favicon,omitempty"`
}

type Locale struct {
	Tag            string `json:"tag,omitempty"`
	Currency       string `json:"currency,omitempty"`
	CurrencyPrefix string `json:"currencyPrefix,omitempty"`
	CurrencySymbol string `json:"currencySymbol,omitempty"`
}

// Contact is the single "call us" line, for SMS and the web app.
type Contact struct {
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
}

type Channels struct {
	SMSSenderID      string  `json:"smsSenderId,omitempty"`
	EmailFromName    string  `json:"emailFromName,omitempty"`
	EmailFromAddress string  `json:"emailFromAddress,omitempty"`
	PushTopic        string  `json:"pushTopic,omitempty"`
	EmailLogoURL     *string `json:"emailLogoUrl,omitempty"`
	PoweredBy        *string `json:"poweredBy,omitempty"`
	FooterLinks      []Link  `json:"footerLinks,omitempty"`
}

type Link struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type Mobile struct {
	AppleAppIDs []string `json:"appleAppIds,omitempty"`
}

type Links struct {
	System string `json:"system,omitempty"`
	App    string `json:"app,omitempty"`
	Viewer string `json:"viewer,omitempty"`
}

type Radiology struct {
	DefaultRadiologist      string `json:"defaultRadiologist,omitempty"`
	DefaultRadiologistTitle string `json:"defaultRadiologistTitle,omitempty"`
}

type DICOM struct {
	AccessionPrefix string `json:"accessionPrefix,omitempty"`
	StudyUIDRoot    string `json:"studyUidRoot,omitempty"`
}

// Document overrides the identity for printed documents only.
type Document struct {
	Name             *string         `json:"name,omitempty"`
	Logo             *string         `json:"logo,omitempty"`
	FooterColor      string          `json:"footerColor,omitempty"`
	Watermark        *string         `json:"watermark,omitempty"`
	WatermarkOpacity *float64        `json:"watermarkOpacity,omitempty"`
	FooterRaw        *string         `json:"footerRaw,omitempty"`
	Contact          [][]ContactItem `json:"contact,omitempty"`
}

type ContactItem struct {
	Icon string `json:"icon,omitempty"`
	Text string `json:"text"`
}

// maxSenderID is Hubtel's limit (system: internal/providers/sms/hubtel.go). An id over it
// is rejected by the provider, and an unregistered one silently fails to deliver.
const maxSenderID = 11

var (
	keyPattern    = regexp.MustCompile(`^[a-z][a-z0-9]{1,30}$`) // as provision-facility.sh
	namePattern   = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)
	hexPattern    = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	uidPattern    = regexp.MustCompile(`^[0-9]+(\.[0-9]+)*$`)
	urlPattern    = regexp.MustCompile(`(?i)^https?://[^\s/]+`)
	tagPattern    = regexp.MustCompile(`^[A-Za-z]{2,3}(-[A-Za-z0-9]{2,8})*$`)
	hostPattern   = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$`)
	dayPattern    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	currencyShape = regexp.MustCompile(`^[A-Z]{3}$`)
)

// Validate checks that a record is well formed. It is what a service runs on every file it
// reads, so it checks shape only, never commercial policy: a service must not refuse to
// start because of a rule only the control plane can judge. See ValidateForPublish.
func (r *Record) Validate() error {
	var errs []error
	bad := func(field, format string, args ...any) {
		errs = append(errs, fmt.Errorf("%s: "+format, append([]any{field}, args...)...))
	}

	if r.Schema != SchemaVersion {
		bad("schema", "is %d; this build reads %d", r.Schema, SchemaVersion)
	}
	if !keyPattern.MatchString(r.Key) {
		bad("key", "%q must be lowercase letters and digits, starting with a letter, 2 to 31 long", r.Key)
	}
	if r.Version < 1 {
		bad("version", "must be 1 or more, got %d", r.Version)
	}
	if !r.Status.valid() {
		bad("status", "%q is not a status", r.Status)
	}
	for i, d := range r.Domains {
		if !hostPattern.MatchString(d) {
			bad(fmt.Sprintf("domains[%d]", i), "%q is not a lowercase hostname", d)
		}
	}
	for name, state := range r.Modules {
		if !namePattern.MatchString(name) {
			bad("modules", "%q is not a module name", name)
		}
		if !state.valid() {
			bad("modules."+name, "%q is not enabled, read_only or off", state)
		}
	}
	seen := make(map[string]bool, len(r.Features))
	for _, f := range r.Features {
		if !namePattern.MatchString(f) {
			bad("features", "%q is not a feature name", f)
		}
		if seen[f] {
			bad("features", "%q is listed twice", f)
		}
		seen[f] = true
	}
	for name, n := range r.Limits {
		if !namePattern.MatchString(name) {
			bad("limits", "%q is not a limit name", name)
		}
		if n < 0 {
			bad("limits."+name, "is %d; a limit can't be negative", n)
		}
	}
	if u := r.UsageMTD; u != nil {
		if !dayPattern.MatchString(u.AsOf) {
			bad("usage_mtd.as_of", "%q is not YYYY-MM-DD", u.AsOf)
		}
		for name, n := range u.Meters {
			if n < 0 {
				bad("usage_mtd.meters."+name, "is %d; usage can't be negative", n)
			}
		}
	}

	errs = append(errs, r.Profile.validate(r.Key)...)
	return errors.Join(errs...)
}

func (p *Profile) validate(recordKey string) []error {
	var errs []error
	bad := func(field, format string, args ...any) {
		errs = append(errs, fmt.Errorf("profile.%s: "+format, append([]any{field}, args...)...))
	}

	if p.Key != recordKey {
		bad("key", "is %q but the record is %q; the profile must be the tenant's own", p.Key, recordKey)
	}
	if strings.TrimSpace(p.Name) == "" {
		bad("name", "is required")
	}
	for field, value := range map[string]string{"primary": p.Primary, "accent": p.Accent, "document.footerColor": p.Document.FooterColor} {
		if value != "" && !hexPattern.MatchString(value) {
			bad(field, "%q is not a hex colour", value)
		}
	}
	if (p.Primary == "") != (p.Accent == "") {
		bad("primary", "primary and accent go together; one without the other is ignored")
	}
	if s := p.Channels.SMSSenderID; len(s) > maxSenderID {
		bad("channels.smsSenderId", "%q is %d characters; Hubtel rejects anything over %d", s, len(s), maxSenderID)
	}
	if s := p.Links.System; s != "" && !urlPattern.MatchString(s) {
		bad("links.system", "%q is not an absolute http(s) URL", s)
	}
	if c := p.Locale.Currency; c != "" && !currencyShape.MatchString(c) {
		bad("locale.currency", "%q is not a three-letter ISO 4217 code", c)
	}
	if t := p.Locale.Tag; t != "" && !tagPattern.MatchString(t) {
		bad("locale.tag", "%q is not a BCP 47 language tag", t)
	}
	if root := p.DICOM.StudyUIDRoot; root != "" {
		if !uidPattern.MatchString(root) {
			bad("dicom.studyUidRoot", "%q is not a DICOM UID", root)
		} else if len(root) > 40 {
			bad("dicom.studyUidRoot", "is %d characters, leaving no room for the suffix within DICOM's 64", len(root))
		}
	}
	return errs
}

// ValidateForPublish is what the control plane checks before publishing a record, on top of
// Validate. prev is the tenant's currently published record, or nil for its first.
func (r *Record) ValidateForPublish(prev *Record) error {
	errs := []error{r.Validate()}

	hasClinical := false
	for _, m := range clinicalModules {
		if r.Modules[m] != "" && r.Modules[m] != ModuleOff {
			hasClinical = true
		}
	}
	if !hasClinical {
		errs = append(errs, errors.New("modules: a tenant needs radiology or lab"))
	}

	if prev != nil {
		if prev.Key != r.Key {
			errs = append(errs, fmt.Errorf("key: %q can't become %q", prev.Key, r.Key))
		}
		if r.Version <= prev.Version {
			errs = append(errs, fmt.Errorf("version: %d must be higher than the published %d", r.Version, prev.Version))
		}
		// Once a clinical module has been on, its records are medical records. Off would
		// hide them from the doctors and patients they belong to.
		for _, m := range clinicalModules {
			was := prev.Modules[m]
			if (was == ModuleEnabled || was == ModuleReadOnly) && r.moduleOrOff(m) == ModuleOff {
				errs = append(errs, fmt.Errorf("modules.%s: was %s, so it can go to read_only but never off", m, was))
			}
		}
	}
	return errors.Join(errs...)
}

// CheckAcross finds collisions that only show when tenants are read together: identifiers
// two tenants must never share.
func CheckAcross(records []Record) error {
	var errs []error
	for _, c := range []struct {
		field, consequence string
		value              func(*Record) string
	}{
		{"profile.dicom.accessionPrefix", "accession numbers collide in any archive holding both",
			func(r *Record) string { return r.Profile.DICOM.AccessionPrefix }},
		{"profile.channels.pushTopic", "one tenant's push notifications arrive on the other's devices",
			func(r *Record) string { return r.Profile.Channels.PushTopic }},
	} {
		owners := map[string][]string{}
		for i := range records {
			if v := strings.TrimSpace(c.value(&records[i])); v != "" {
				owners[v] = append(owners[v], records[i].Key)
			}
		}
		for _, v := range slices.Sorted(maps.Keys(owners)) {
			if keys := owners[v]; len(keys) > 1 {
				errs = append(errs, fmt.Errorf("%s %q is shared by %s: %s", c.field, v, strings.Join(keys, ", "), c.consequence))
			}
		}
	}

	hosts := map[string]string{}
	for _, r := range records {
		for _, d := range r.Domains {
			if other, ok := hosts[d]; ok && other != r.Key {
				errs = append(errs, fmt.Errorf("domains: %q is claimed by both %s and %s", d, other, r.Key))
			}
			hosts[d] = r.Key
		}
	}
	return errors.Join(errs...)
}

func (r *Record) moduleOrOff(m string) ModuleState {
	if s, ok := r.Modules[m]; ok {
		return s
	}
	return ModuleOff
}
