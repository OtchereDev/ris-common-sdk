// Package letterhead is what a printed document repeats: the logo across the top, the contact
// strip along the bottom, the watermark behind the content. It is built from a tenant's
// profile, so pdf-gen prints it and the Console previews it from the same code.
//
// The rendering and the icons are pdf-gen's, from its feature/pluslab branch
// (internal/generator: brand_render.go, icons.go, and facility.go's profile mapping), moved
// here unchanged apart from reading a tenant.Profile and resolving images by key.
package letterhead

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
)

const (
	// DefaultWatermarkOpacity keeps the mark visible as branding without competing with the
	// result values printed over it.
	DefaultWatermarkOpacity = 0.06

	// DefaultFooterColor is the neutral the contact strip falls back to.
	DefaultFooterColor = "#5a6373"

	// DefaultCurrency is what every document was printed in before the profile could say
	// otherwise.
	DefaultCurrency = "GHS"
)

// Letterhead carries the marks a document repeats.
type Letterhead struct {
	Name string
	Site string
	// Logo and Watermark are data: URIs. The browser that prints the header and footer
	// renders them in their own context, which loads nothing from the network, so the bytes
	// travel with the markup.
	Logo             string
	Watermark        string
	WatermarkOpacity float64
	FooterColor      string
	// FooterRaw is the footer markup, used verbatim when set: an escape hatch for a facility
	// whose footer already exists and should not be restyled. Contact and FooterColor are
	// ignored when it is set.
	FooterRaw string
	// Contact is the strip along the bottom of every page: one entry per printed line, each
	// a run of items that may carry an icon.
	Contact  [][]tenant.ContactItem
	Currency string
}

// Images turns what a profile's image slot holds, a key in the asset bucket or artwork from a
// record made before keys, into a data: URI.
type Images func(ctx context.Context, value string) (string, error)

// FromProfile resolves a tenant's letterhead.
//
// The document overrides win where they are set; otherwise the facility's one identity serves
// both screen and page. The mark doubles as the watermark where the profile names no separate
// one: it is the square lockup, which is what sits behind a page of results.
func FromProfile(ctx context.Context, p *tenant.Profile, images Images) (Letterhead, error) {
	pick := func(slots ...tenant.ImageSlot) string {
		for _, s := range slots {
			if v := p.Image(s); v != "" {
				return v
			}
		}
		return ""
	}
	l := Letterhead{
		Name:             orDefault(deref(p.Document.Name), p.Name),
		Site:             p.Site,
		WatermarkOpacity: derefFloat(p.Document.WatermarkOpacity),
		FooterColor:      p.Document.FooterColor,
		FooterRaw:        deref(p.Document.FooterRaw),
		Contact:          p.Document.Contact,
		Currency:         orDefault(p.Locale.CurrencyPrefix, p.Locale.Currency),
	}
	var err error
	if l.Logo, err = resolve(ctx, images, pick(tenant.ImageDocumentLogo, tenant.ImageLogoFull)); err != nil {
		return l, fmt.Errorf("logo: %w", err)
	}
	if l.Watermark, err = resolve(ctx, images, pick(tenant.ImageDocumentWatermark, tenant.ImageLogoMark)); err != nil {
		return l, fmt.Errorf("watermark: %w", err)
	}
	if l.WatermarkOpacity <= 0 {
		l.WatermarkOpacity = DefaultWatermarkOpacity
	}
	if l.FooterColor == "" {
		l.FooterColor = DefaultFooterColor
	}
	if l.Currency == "" {
		l.Currency = DefaultCurrency
	}
	return l, nil
}

func resolve(ctx context.Context, images Images, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	return images(ctx, value)
}

// AddressLines flattens the contact strip into plain lines, for documents that stack the
// facility's address in a corner rather than running it along the foot of the page.
func (l Letterhead) AddressLines() []string {
	out := make([]string, 0, len(l.Contact))
	for _, line := range l.Contact {
		parts := make([]string, 0, len(line))
		for _, item := range line {
			if item.Text != "" {
				parts = append(parts, item.Text)
			}
		}
		if len(parts) > 0 {
			out = append(out, strings.Join(parts, "  "))
		}
	}
	return out
}

// ContactHTML is the contact strip as markup, for a document whose template paints its own
// page furniture (the lab report). Styling stays in that template's stylesheet; only the
// structure is built here.
func (l Letterhead) ContactHTML() string {
	var sb strings.Builder
	for _, line := range l.Contact {
		if len(line) == 0 {
			continue
		}
		sb.WriteString(`<div class="cline">`)
		for _, item := range line {
			sb.WriteString(`<span class="citem">`)
			sb.WriteString(Icon(item.Icon))
			sb.WriteString(html.EscapeString(item.Text))
			sb.WriteString(`</span>`)
		}
		sb.WriteString(`</div>`)
	}
	return sb.String()
}

// MaxLogoHeight keeps the logo inside the one-inch band the header is printed in, with room
// to spare. Sized by width alone, as pdf-gen's header was, a square or round logo ran past
// the band into the report.
const MaxLogoHeight = "0.8in"

// HeaderHTML is the repeating page header the browser draws above the content: the
// facility's logo, right aligned, 30% of the page wide unless that would make it taller than
// MaxLogoHeight. A wide logo prints as before; a tall one is scaled down to fit, still
// right aligned.
func (l Letterhead) HeaderHTML() string {
	if l.Logo == "" {
		return ""
	}
	return fmt.Sprintf(
		`<div style="width:100%%;display:flex;justify-content:flex-end;padding:0 40px;">`+
			`<img src="%s" style="width:30%%;max-height:%s;object-fit:contain;object-position:right center" alt="%s" /></div>`,
		l.Logo, MaxLogoHeight, html.EscapeString(l.Name),
	)
}

// FooterHTML is the repeating contact strip.
//
// Built as one inline-styled string because the browser renders header and footer templates
// in their own context: the page's stylesheet does not reach them, and neither do its web
// fonts, so this states every rule it needs and asks for a generic family it can rely on.
func (l Letterhead) FooterHTML() string {
	// A facility that brought its own footer keeps it, untouched.
	if l.FooterRaw != "" {
		return l.FooterRaw
	}
	if len(l.Contact) == 0 && l.Site == "" {
		return ""
	}
	colour := orDefault(l.FooterColor, DefaultFooterColor)

	var sb strings.Builder
	fmt.Fprintf(&sb,
		`<div style="width:100%%;padding:0 40px;font-family:Helvetica,Arial,sans-serif;`+
			`font-size:8px;font-weight:600;line-height:1.7;text-align:center;color:%s;">`,
		html.EscapeString(colour),
	)
	for _, line := range l.Contact {
		sb.WriteString(`<div>`)
		for i, item := range line {
			if i > 0 {
				sb.WriteString(`<span style="padding:0 6px;"></span>`)
			}
			if svg := Icon(item.Icon); svg != "" {
				// Sized and nudged onto the baseline inline, since the class the page
				// stylesheet would carry does not reach this context.
				sb.WriteString(strings.Replace(svg, `class="ci"`,
					`style="width:8px;height:8px;vertical-align:-1px;margin-right:3px;"`, 1))
			}
			sb.WriteString(html.EscapeString(item.Text))
		}
		sb.WriteString(`</div>`)
	}
	if l.Site != "" {
		fmt.Fprintf(&sb, `<div style="font-weight:700;">%s</div>`, html.EscapeString(l.Site))
	}
	sb.WriteString(`</div>`)
	return sb.String()
}

// WatermarkHTML is the mark behind a page's content, faint, centred, and clipped by its own
// wrapper so it never adds a page. "" when there is no watermark.
func (l Letterhead) WatermarkHTML() string {
	if l.Watermark == "" {
		return ""
	}
	return fmt.Sprintf(
		`<div style="position:absolute;inset:0;overflow:hidden;z-index:-1;">`+
			`<img src="%s" alt="" style="position:absolute;top:50%%;left:50%%;`+
			`transform:translate(-50%%,-50%%);width:148mm;opacity:%g;" /></div>`,
		l.Watermark, l.WatermarkOpacity,
	)
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
