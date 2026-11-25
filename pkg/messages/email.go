package messages

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
)

// Email is a rendered email: its subject, its HTML, and the same as plain text for clients
// that show no HTML.
type Email struct {
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text"`
}

// defaultPrimary is the app's own action colour, for a tenant that sets none.
const defaultPrimary = "#597a49"

// RenderEmail renders a message's email for a tenant. assetsURL is where the asset bucket is
// read (PLATFORM_ASSETS_URL): an email client fetches the logo from there, as it can't show
// an image carried inside the message.
func RenderEmail(r *tenant.Record, id string, v Values, assetsURL string) (Email, error) {
	d, err := def(id, tenant.ChannelEmail)
	if err != nil {
		return Email{}, err
	}
	c, err := Copy(r, id)
	if err != nil {
		return Email{}, err
	}
	all := withRecord(r, v)
	p := &r.Profile

	data := emailData{
		Subject:     fill(c.Subject, all),
		Heading:     fill(c.Heading, all),
		Paragraphs:  paragraphs(fill(c.Body, all)),
		Primary:     orDefault(p.Primary, defaultPrimary),
		Logo:        logoURL(r, assetsURL),
		Facility:    all["facility"],
		Contact:     contactLine(p),
		FooterLinks: p.Channels.FooterLinks,
		PoweredBy:   deref(p.Channels.PoweredBy),
	}
	if d.Code {
		data.Code = all["code"]
	}
	if c.Button != "" {
		data.Button, data.Link = fill(c.Button, all), all["link"]
	}

	var html bytes.Buffer
	if err := layout.Execute(&html, data); err != nil {
		return Email{}, err
	}
	return Email{Subject: data.Subject, HTML: html.String(), Text: data.text()}, nil
}

type emailData struct {
	Subject, Heading string
	Paragraphs       [][]string
	Code             string
	Button, Link     string
	Primary, Logo    string
	Facility         string
	Contact          string
	FooterLinks      []tenant.Link
	PoweredBy        string
}

func (d emailData) text() string {
	var sb strings.Builder
	sb.WriteString(d.Heading + "\n\n")
	for _, p := range d.Paragraphs {
		sb.WriteString(strings.Join(p, "\n") + "\n\n")
	}
	if d.Code != "" {
		sb.WriteString(d.Code + "\n\n")
	}
	if d.Link != "" {
		fmt.Fprintf(&sb, "%s: %s\n\n", d.Button, d.Link)
	}
	sb.WriteString(d.Facility)
	if d.Contact != "" {
		sb.WriteString("\n" + d.Contact)
	}
	return sb.String()
}

// paragraphs splits a body at blank lines, and each paragraph into its lines.
func paragraphs(body string) [][]string {
	var out [][]string
	for _, para := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n\n") {
		if strings.TrimSpace(para) == "" {
			continue
		}
		out = append(out, strings.Split(strings.TrimSpace(para), "\n"))
	}
	return out
}

func logoURL(r *tenant.Record, assetsURL string) string {
	if u := deref(r.Profile.Channels.EmailLogoURL); u != "" {
		return u
	}
	key := r.Profile.Image(tenant.ImageLogoFull)
	if strings.HasPrefix(key, "tenants/") && assetsURL != "" {
		return strings.TrimRight(assetsURL, "/") + "/" + key
	}
	return ""
}

func contactLine(p *tenant.Profile) string {
	var parts []string
	for _, s := range []string{p.Contact.Address, p.Contact.Phone, p.Contact.Email} {
		if s = strings.TrimSpace(s); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " · ")
}

func orDefault(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// layout is the one email layout: tables and inline styles, as email clients need, 600 wide.
// Every value is escaped by html/template; the copy is text, never markup.
var layout = template.Must(template.New("email").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="x-apple-disable-message-reformatting"><title>{{.Subject}}</title></head>
<body style="margin:0;padding:0;background:#f4f5f0;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="background:#f4f5f0;">
<tr><td align="center" style="padding:32px 16px;">
<table role="presentation" width="600" cellpadding="0" cellspacing="0" style="width:100%;max-width:600px;background:#ffffff;border-radius:12px;font-family:Helvetica,Arial,sans-serif;color:#191d17;">
{{if .Logo}}<tr><td style="padding:28px 32px 0;"><img src="{{.Logo}}" alt="{{.Facility}}" style="max-height:48px;max-width:220px;display:block;"></td></tr>{{end}}
<tr><td style="padding:24px 32px 8px;"><h1 style="margin:0;font-size:24px;line-height:32px;font-weight:700;">{{.Heading}}</h1></td></tr>
{{range .Paragraphs}}<tr><td style="padding:8px 32px;font-size:16px;line-height:24px;">{{range $i, $line := .}}{{if $i}}<br>{{end}}{{$line}}{{end}}</td></tr>
{{end}}{{if .Code}}<tr><td style="padding:16px 32px;"><div style="display:inline-block;padding:14px 24px;border-radius:10px;background:#f4f5f0;font-family:Menlo,Consolas,monospace;font-size:28px;letter-spacing:6px;font-weight:700;">{{.Code}}</div></td></tr>
{{end}}{{if .Link}}<tr><td style="padding:16px 32px;"><a href="{{.Link}}" style="display:inline-block;padding:12px 24px;border-radius:10px;background:{{.Primary}};color:#ffffff;text-decoration:none;font-weight:700;font-size:16px;">{{.Button}}</a></td></tr>
{{end}}<tr><td style="padding:24px 32px 28px;font-size:13px;line-height:20px;color:#60665c;border-top:1px solid #eaede6;">
<strong style="color:#191d17;">{{.Facility}}</strong>{{if .Contact}}<br>{{.Contact}}{{end}}
{{if .FooterLinks}}<br>{{range $i, $l := .FooterLinks}}{{if $i}} · {{end}}<a href="{{$l.Href}}" style="color:#60665c;">{{$l.Label}}</a>{{end}}{{end}}
{{if .PoweredBy}}<br><span style="font-size:12px;">{{.PoweredBy}}</span>{{end}}
</td></tr></table></td></tr></table></body></html>`))
