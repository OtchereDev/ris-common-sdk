package letterhead

import "strings"

// Contact strip icons.
//
// Inline SVG rather than an icon font or emoji, because the renderer prints in a container
// with only the fonts it was built with: an emoji glyph shows as a tofu box, and a webfont
// is one more network fetch the browser has no route for. These draw at any size, inherit
// the surrounding colour through currentColor, and cost nothing to load.
//
// Stroked at 1.8 so they still read at the 7.5pt the footer sets.
const (
	IconLocation = "location"
	IconPhone    = "phone"
	IconEmail    = "email"
	IconWeb      = "web"
	IconMail     = "mail"
	IconPerson   = "person"
)

var icons = map[string]string{
	IconLocation: `<path d="M12 21s7-6.4 7-11a7 7 0 1 0-14 0c0 4.6 7 11 7 11Z"/><circle cx="12" cy="10" r="2.6"/>`,
	IconPhone:    `<path d="M6.3 3.5h3l1.5 3.8-2 1.4a11.5 11.5 0 0 0 5.5 5.5l1.4-2 3.8 1.5v3a2 2 0 0 1-2.2 2A16.5 16.5 0 0 1 4.3 5.7a2 2 0 0 1 2-2.2Z"/>`,
	IconEmail:    `<rect x="2.8" y="5" width="18.4" height="14" rx="2.2"/><path d="m3.4 6.6 8.6 6.2 8.6-6.2"/>`,
	IconWeb:      `<circle cx="12" cy="12" r="9"/><path d="M3.2 12h17.6"/><path d="M12 3a15 15 0 0 1 0 18 15 15 0 0 1 0-18Z"/>`,
	IconMail:     `<path d="M3 8.5 12 3l9 5.5V19a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2Z"/><path d="m3 8.5 9 6 9-6"/>`,
	IconPerson:   `<circle cx="12" cy="8" r="3.6"/><path d="M4.8 20.4a7.2 7.2 0 0 1 14.4 0"/>`,
}

// Icon returns the markup for an icon name, or an empty string for an unknown or absent
// one. A contact entry with no icon is a plain run of text, which is what a facility that
// has not configured icons gets.
func Icon(name string) string {
	path, ok := icons[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return ""
	}

	return `<svg class="ci" viewBox="0 0 24 24" fill="none" stroke="currentColor" ` +
		`stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" ` +
		`aria-hidden="true">` + path + `</svg>`
}
