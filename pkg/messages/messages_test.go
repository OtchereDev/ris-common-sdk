package messages_test

import (
	"strings"
	"testing"

	"github.com/OtchereDev/ris-common-sdk/pkg/messages"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant/tenanttest"
)

// The catalogue's own defaults must keep its rules: each channel has its words, and the words
// say only what the message can.
func TestCatalogueDefaults(t *testing.T) {
	for _, d := range tenant.Catalogue {
		texts := map[string]string{}
		if d.Has(tenant.ChannelEmail) {
			if d.Email == nil || d.Email.Subject == "" || d.Email.Heading == "" || d.Email.Body == "" {
				t.Errorf("%s: an email message without its default copy", d.ID)
				continue
			}
			texts["subject"], texts["heading"], texts["body"], texts["button"] = d.Email.Subject, d.Email.Heading, d.Email.Body, d.Email.Button
		}
		if d.Has(tenant.ChannelSMS) {
			if d.SMS == "" {
				t.Errorf("%s: an SMS message without its default text", d.ID)
			}
			texts["sms"] = d.SMS
		}
		if d.Has(tenant.ChannelWhatsApp) && len(d.WhatsAppVariables) == 0 {
			t.Errorf("%s: a WhatsApp message without its template variables", d.ID)
		}
		for where, text := range texts {
			for _, name := range tenant.Placeholders(text) {
				if !d.Allows(name) {
					t.Errorf("%s %s says {{%s}}, which it can't", d.ID, where, name)
				}
			}
			if strings.Contains(strings.ToLower(text), "spectra") {
				t.Errorf("%s %s names a customer: defaults are every tenant's", d.ID, where)
			}
		}
	}
}

func record() tenant.Record {
	r := tenanttest.Record()
	r.Profile.Channels.EmailFromName = "Demo Diagnostics"
	return r
}

func TestSMS(t *testing.T) {
	r := record()
	got, err := messages.SMS(&r, "login_code", messages.Values{"first_name": "Ama", "code": "482913"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "Hello Ama, your Demo Diagnostics sign-in code is 482913. It expires in 10 minutes." {
		t.Errorf("default: %q", got)
	}

	// The tenant's own words, and an empty name leaves no gap before the comma.
	r.Messages = &tenant.Messages{SMS: map[string]string{"login_code": "Hi {{first_name}}, {{code}} is your code for {{facility}}."}}
	got, _ = messages.SMS(&r, "login_code", messages.Values{"first_name": "", "code": "482913"})
	if got != "Hi, 482913 is your code for Demo Diagnostics." {
		t.Errorf("own text: %q", got)
	}

	if _, err := messages.SMS(&r, "welcome", nil); err == nil {
		t.Error("welcome is email only")
	}
}

func TestEmail(t *testing.T) {
	r := record()
	r.Profile.Primary, r.Profile.Accent = "#1f6f8b", "#f2a541"
	r.Profile.Logo.Full = "tenants/demo/logo.full-1111111111111111.png"
	powered := "Powered by SpectraRIS"
	r.Profile.Channels.PoweredBy = &powered
	r.Profile.Channels.FooterLinks = []tenant.Link{{Label: "Privacy", Href: "https://demo.example/privacy"}}
	r.Messages = &tenant.Messages{Email: map[string]tenant.EmailCopy{"report_ready": {Heading: "Results for {{first_name}} <now>"}}}

	e, err := messages.RenderEmail(&r, "report_ready", messages.Samples("report_ready"), "https://assets.example.com/")
	if err != nil {
		t.Fatal(err)
	}
	if e.Subject != "Your Demo Diagnostics report is ready" {
		t.Errorf("subject keeps its default: %q", e.Subject)
	}
	for _, want := range []string{
		"Results for Ama &lt;now&gt;",                                            // the tenant's heading, escaped
		"https://assets.example.com/tenants/demo/logo.full-1111111111111111.png", // the logo, from the bucket
		"background:#1f6f8b",                                                     // the button in the tenant's colour
		`href="https://quest.spectraris.com/r/abc123"`, ">View report</a>",
		"Privacy", "Powered by SpectraRIS", "30 000 0000",
	} {
		if !strings.Contains(e.HTML, want) {
			t.Errorf("the email is missing %q", want)
		}
	}
	if !strings.Contains(e.Text, "View report: https://quest.spectraris.com/r/abc123") {
		t.Errorf("plain text: %q", e.Text)
	}

	code, _ := messages.RenderEmail(&r, "login_code", messages.Samples("login_code"), "")
	if !strings.Contains(code.HTML, "482913") || strings.Contains(code.HTML, "text-decoration:none;font-weight:700") {
		t.Errorf("a code email shows the code and has no button:\n%s", code.HTML)
	}
}

func TestWhatsApp(t *testing.T) {
	r := record()
	if _, _, _, ok, _ := messages.WhatsApp(&r, "report_ready", nil); ok {
		t.Error("no template chosen should mean no WhatsApp")
	}
	r.Messages = &tenant.Messages{WhatsApp: &tenant.WhatsApp{
		MessagingServiceSID: "MG" + strings.Repeat("a", 32),
		Templates:           map[string]string{"report_ready": "HX" + strings.Repeat("b", 32)},
	}}
	sender, tpl, vars, ok, err := messages.WhatsApp(&r, "report_ready", messages.Values{"first_name": "Ama", "link": "abc"})
	if err != nil || !ok || sender != "MG"+strings.Repeat("a", 32) || tpl != "HX"+strings.Repeat("b", 32) || vars != `{"1":"Ama","2":"abc"}` {
		t.Errorf("%q %q %q %v %v", sender, tpl, vars, ok, err)
	}
}

func TestCount(t *testing.T) {
	for _, c := range []struct {
		text    string
		parts   int
		unicode bool
	}{
		{"", 0, false},
		{strings.Repeat("a", 160), 1, false},
		{strings.Repeat("a", 161), 2, false},
		{strings.Repeat("a", 159) + "€", 2, false}, // € takes two
		{"Visit us 📍", 1, true},
		{strings.Repeat("a", 70) + "é", 1, false}, // é is in the GSM alphabet
		{strings.Repeat("a", 70) + "✓", 2, true},
	} {
		s := messages.Count(c.text)
		if s.Parts != c.parts || s.Unicode != c.unicode {
			t.Errorf("%.20q…: %+v, want %d parts, unicode %v", c.text, s, c.parts, c.unicode)
		}
	}
}

func TestRecordRules(t *testing.T) {
	r := record()
	r.Messages = &tenant.Messages{
		Email:    map[string]tenant.EmailCopy{"login_code": {Button: "Go", Body: "Your {{password}}"}, "nope": {}},
		SMS:      map[string]string{"welcome": "Hi"},
		WhatsApp: &tenant.WhatsApp{MessagingServiceSID: "MG123", Templates: map[string]string{"report_ready": "HXbad"}},
	}
	err := r.Validate()
	if err == nil {
		t.Fatal("expected problems")
	}
	for _, want := range []string{
		"messages.email.login_code.button: this email has no button",
		"messages.email.login_code.body: {{password}} isn't something this message can say",
		"messages.email.nope: isn't an email message",
		"messages.sms.welcome: isn't an SMS message",
		"messages.whatsapp.messaging_service_sid:",
		"messages.whatsapp.templates.report_ready:",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("missing %q in:\n%v", want, err)
		}
	}

	// Templates without the tenant's own sender won't publish.
	ok := record()
	ok.Messages = &tenant.Messages{WhatsApp: &tenant.WhatsApp{Templates: map[string]string{"report_ready": "HX" + strings.Repeat("b", 32)}}}
	if err := ok.Validate(); err != nil {
		t.Errorf("a service reads it: %v", err)
	}
	if err := ok.ValidateForPublish(nil); err == nil || !strings.Contains(err.Error(), "own WhatsApp sender") {
		t.Errorf("publish: %v", err)
	}

	// Two tenants can't share a sender.
	a, b := record(), record()
	b.Key, b.Profile.Key = "other", "other"
	b.Profile.DICOM.AccessionPrefix, b.Profile.Channels.PushTopic = "OTH", "OtherGeneral"
	sid := "MG" + strings.Repeat("c", 32)
	a.Messages = &tenant.Messages{WhatsApp: &tenant.WhatsApp{MessagingServiceSID: sid}}
	b.Messages = &tenant.Messages{WhatsApp: &tenant.WhatsApp{MessagingServiceSID: sid}}
	if err := tenant.CheckAcross([]tenant.Record{a, b}); err == nil || !strings.Contains(err.Error(), "under the other's name") {
		t.Errorf("shared sender: %v", err)
	}
}
