// Package messages renders a tenant's transactional messages: the text of an SMS, an email in
// the one branded layout, and what a WhatsApp template is sent with. The catalogue, defaults
// and rules are pkg/tenant's (messages.go); a tenant's record holds only what it changes.
//
// The system service sends with this package and the Console previews with it, so a preview
// is what is sent.
package messages

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
)

// Values fills a message's variables. "facility" and "support_phone" come from the record
// when they aren't given.
type Values map[string]string

func def(id string, c tenant.Channel) (*tenant.MessageDef, error) {
	d := tenant.Message(id)
	if d == nil {
		return nil, fmt.Errorf("no message %q", id)
	}
	if !d.Has(c) {
		return nil, fmt.Errorf("%s isn't sent by %s", id, c)
	}
	return d, nil
}

// withRecord adds what the record knows to the values given. The record's own name and phone
// win over any given; a given one only fills a gap the record leaves, as a preview's sample.
func withRecord(r *tenant.Record, v Values) Values {
	out := Values{}
	for k, val := range v {
		out[k] = val
	}
	for k, val := range map[string]string{"facility": Facility(r), "support_phone": r.Profile.Contact.Phone} {
		if strings.TrimSpace(val) != "" {
			out[k] = val
		}
	}
	return out
}

// Facility is the tenant's name as its messages are signed: the name its emails come from,
// or else its name in the app.
func Facility(r *tenant.Record) string {
	if n := strings.TrimSpace(r.Profile.Channels.EmailFromName); n != "" {
		return n
	}
	return r.Profile.Name
}

var (
	placeholder = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_]*)\s*\}\}`)
	// A variable that was empty leaves a space before punctuation: "Hi ," for "Hi {{first_name}},".
	strayGap = regexp.MustCompile(`[ \t]+([,.!?])`)
)

// fill replaces each {{variable}} with its value. One with no value is left showing, so a
// preview or a test makes the gap plain rather than sending "Dear ,".
func fill(text string, v Values) string {
	out := placeholder.ReplaceAllStringFunc(text, func(m string) string {
		name := placeholder.FindStringSubmatch(m)[1]
		if val, ok := v[name]; ok {
			return val
		}
		return m
	})
	return strayGap.ReplaceAllString(out, "$1")
}

// SMS is a message's text message.
func SMS(r *tenant.Record, id string, v Values) (string, error) {
	d, err := def(id, tenant.ChannelSMS)
	if err != nil {
		return "", err
	}
	text := d.SMS
	if r.Messages != nil && strings.TrimSpace(r.Messages.SMS[id]) != "" {
		text = r.Messages.SMS[id]
	}
	return fill(text, withRecord(r, v)), nil
}

// Copy is an email's words for the tenant: its own where it set them, the defaults elsewhere.
func Copy(r *tenant.Record, id string) (tenant.EmailCopy, error) {
	d, err := def(id, tenant.ChannelEmail)
	if err != nil {
		return tenant.EmailCopy{}, err
	}
	c := *d.Email
	if r.Messages != nil {
		own := r.Messages.Email[id]
		for _, f := range []struct{ own, into *string }{
			{&own.Subject, &c.Subject}, {&own.Heading, &c.Heading}, {&own.Body, &c.Body}, {&own.Button, &c.Button},
		} {
			if strings.TrimSpace(*f.own) != "" {
				*f.into = *f.own
			}
		}
	}
	return c, nil
}

// WhatsApp is how a message goes out on the tenant's WhatsApp: its sender, the approved
// template, and the template's variables as Twilio takes them ({"1": ..., "2": ...}).
// ok is false when the tenant hasn't chosen a template for the message.
func WhatsApp(r *tenant.Record, id string, v Values) (sender, template, variables string, ok bool, err error) {
	d, err := def(id, tenant.ChannelWhatsApp)
	if err != nil {
		return "", "", "", false, err
	}
	if r.Messages == nil || r.Messages.WhatsApp == nil || r.Messages.WhatsApp.Templates[id] == "" {
		return "", "", "", false, nil
	}
	all := withRecord(r, v)
	numbered := map[string]string{}
	for i, name := range d.WhatsAppVariables {
		numbered[strconv.Itoa(i+1)] = all[name]
	}
	raw, err := json.Marshal(numbered)
	if err != nil {
		return "", "", "", false, err
	}
	w := r.Messages.WhatsApp
	return w.MessagingServiceSID, w.Templates[id], string(raw), true, nil
}

// Samples are a message's variables filled with the catalogue's sample values, for previews.
func Samples(id string) Values {
	out := Values{}
	if d := tenant.Message(id); d != nil {
		for _, v := range d.Variables {
			out[v.Name] = v.Sample
		}
	}
	return out
}

// Segments is how an SMS is billed: how many parts it is sent in, and whether it needs
// Unicode, which shortens each part.
type Segments struct {
	Characters int  `json:"characters"`
	Parts      int  `json:"parts"`
	Unicode    bool `json:"unicode"`
}

// GSM 03.38: the alphabet a plain text message is sent in, and the characters that take two.
const (
	gsmBasic    = "@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà"
	gsmExtended = "^{}\\[~]|€\f"
)

// Count measures an SMS: 160 characters in one part, 153 a part beyond that; 70 and 67 when
// any character needs Unicode.
func Count(text string) Segments {
	units, unicode := 0, false
	for _, r := range text {
		switch {
		case strings.ContainsRune(gsmBasic, r):
			units++
		case strings.ContainsRune(gsmExtended, r):
			units += 2
		default:
			unicode = true
		}
	}
	s := Segments{Characters: units, Unicode: unicode}
	single, multi := 160, 153
	if unicode {
		s.Characters = len([]rune(text))
		single, multi = 70, 67
	}
	switch {
	case s.Characters == 0:
		s.Parts = 0
	case s.Characters <= single:
		s.Parts = 1
	default:
		s.Parts = (s.Characters + multi - 1) / multi
	}
	return s
}
