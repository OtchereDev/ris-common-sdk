package tenant

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
)

// The transactional messages a tenant's patients, doctors and staff get, and what each may
// say. Every message has product defaults, written for no tenant in particular; a tenant's
// record holds only what it changes. pkg/messages renders them.
//
// Campaigns, which tenants write themselves, are not these: they live in the system service.

// Channel is how a message is sent.
type Channel string

const (
	ChannelEmail    Channel = "email"
	ChannelSMS      Channel = "sms"
	ChannelWhatsApp Channel = "whatsapp"
)

// Messages is what a tenant changes about its messages. Empty means the defaults.
type Messages struct {
	// Email is copy, by message id, inside the one branded layout: never the layout itself.
	Email map[string]EmailCopy `json:"email,omitempty"`
	// SMS is the text, by message id.
	SMS map[string]string `json:"sms,omitempty"`
	// WhatsApp is the tenant's own sender and the approved templates it sends.
	WhatsApp *WhatsApp `json:"whatsapp,omitempty"`
}

// EmailCopy is the words of an email. A field left empty takes the default.
type EmailCopy struct {
	Subject string `json:"subject,omitempty"`
	Heading string `json:"heading,omitempty"`
	// Body is paragraphs separated by a blank line.
	Body   string `json:"body,omitempty"`
	Button string `json:"button,omitempty"`
}

// WhatsApp is a tenant's WhatsApp sender: a Twilio Messaging Service holding its own number,
// registered with Meta under its name, and the approved Content template for each message.
// Meta approves a template's wording, so the wording is chosen in Twilio, not here.
type WhatsApp struct {
	MessagingServiceSID string            `json:"messaging_service_sid,omitempty"`
	Templates           map[string]string `json:"templates,omitempty"`
}

// Variable is something a message can say, filled in when it is sent.
type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Sample is what previews show.
	Sample string `json:"sample"`
}

// MessageDef is one message: when it is sent, on which channels, what it may say, and its
// default words.
type MessageDef struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Channels    []Channel  `json:"channels"`
	Variables   []Variable `json:"variables"`
	Email       *EmailCopy `json:"email,omitempty"`
	SMS         string     `json:"sms,omitempty"`
	// WhatsAppVariables is what a WhatsApp template's numbered placeholders hold, in order:
	// {{1}} is the first. A template chosen for the message must use these and no others.
	WhatsAppVariables []string `json:"whatsapp_variables,omitempty"`
	// Code says the email shows {{code}} in a box, and Button that it links to {{link}}.
	Code bool `json:"code,omitempty"`
}

// Every message has the tenant's own name and contact line.
var common = []Variable{
	{"facility", "The tenant's name, as its emails are signed", "Quest Medical Imaging"},
	{"support_phone", "The tenant's phone number", "+233 30 000 0000"},
}

func vars(extra ...Variable) []Variable { return append(slices.Clone(common), extra...) }

var (
	vFirstName = Variable{"first_name", "The person's first name", "Ama"}
	vCode      = Variable{"code", "The one-time code", "482913"}
	vLink      = Variable{"link", "Where the button or link goes", "https://quest.spectraris.com/r/abc123"}
	vDoctor    = Variable{"doctor", "The referring doctor's name", "Kofi Mensah"}
	vProcedure = Variable{"procedure", "The examination", "CT Chest"}
	vPhone     = Variable{"phone", "The centre's phone number", "+233 30 000 0001"}
	vAddress   = Variable{"address", "The centre's address", "1 Example Street, Accra"}
	vEmail     = Variable{"email", "The person's email address", "ama@example.com"}
	vPassword  = Variable{"password", "Their first password", "pE7-xq2-Lm"}
	vNumber    = Variable{"whatsapp_number", "The WhatsApp number tried", "+233 24 000 0000"}
	vReason    = Variable{"reason", "Why it failed", "the number isn't on WhatsApp"}
)

// Catalogue is every message, in the order the Console lists them.
var Catalogue = []MessageDef{
	{
		ID: "report_ready", Name: "Report ready", Description: "A patient's report can be viewed.",
		Channels:  []Channel{ChannelEmail, ChannelWhatsApp},
		Variables: vars(vFirstName, vLink),
		Email: &EmailCopy{
			Subject: "Your {{facility}} report is ready",
			Heading: "Your report is ready",
			Body:    "Hello {{first_name}},\n\nThe report on your examination at {{facility}} is ready. Use the button below to view it.\n\nIf you have any questions, call us on {{support_phone}}.",
			Button:  "View report",
		},
		WhatsAppVariables: []string{"first_name", "link"},
	},
	{
		ID: "login_code", Name: "Sign-in code", Description: "A code to finish signing in.",
		Channels:  []Channel{ChannelEmail, ChannelSMS},
		Variables: vars(vFirstName, vCode),
		Code:      true,
		Email: &EmailCopy{
			Subject: "Your {{facility}} sign-in code",
			Heading: "Your sign-in code",
			Body:    "Hello {{first_name}},\n\nUse this code to finish signing in to {{facility}}. It expires in 10 minutes.\n\nIf you didn't try to sign in, you can ignore this email.",
		},
		SMS: "Hello {{first_name}}, your {{facility}} sign-in code is {{code}}. It expires in 10 minutes.",
	},
	{
		ID: "password_reset", Name: "Password reset code", Description: "A code to set a new password.",
		Channels:  []Channel{ChannelEmail, ChannelSMS},
		Variables: vars(vFirstName, vCode),
		Code:      true,
		Email: &EmailCopy{
			Subject: "Reset your {{facility}} password",
			Heading: "Reset your password",
			Body:    "Hello {{first_name}},\n\nUse this code to set a new password for your {{facility}} account. It expires in 10 minutes.\n\nIf you didn't ask to reset your password, you can ignore this email.",
		},
		SMS: "Hello {{first_name}}, your {{facility}} password reset code is {{code}}. It expires in 10 minutes.",
	},
	{
		ID: "welcome", Name: "Welcome", Description: "An account was made for a doctor or member of staff.",
		Channels:  []Channel{ChannelEmail},
		Variables: vars(vFirstName, vEmail, vPassword, vLink),
		Email: &EmailCopy{
			Subject: "Welcome to {{facility}}",
			Heading: "Welcome to {{facility}}",
			Body:    "Hello {{first_name}},\n\nAn account has been made for you. Sign in with {{email}} and this password: {{password}}\n\nYou'll be asked to choose your own password when you first sign in.",
			Button:  "Sign in",
		},
	},
	{
		ID: "appointment_requested", Name: "Appointment requested", Description: "A doctor requested an examination for a patient.",
		Channels:  []Channel{ChannelSMS},
		Variables: vars(vFirstName, vDoctor, vProcedure, vLink, vAddress, vPhone),
		SMS:       "Hi {{first_name}}, Dr. {{doctor}} has requested an appointment for {{procedure}}.\n\nSee your request form and get your report in our app: {{link}}\n\nVisit us at {{address}}, or call {{phone}} for directions.\n\n{{facility}}",
	},
	{
		ID: "report_delayed", Name: "Report delayed", Description: "A report is taking longer than usual.",
		Channels:  []Channel{ChannelSMS},
		Variables: vars(vFirstName, vProcedure, vPhone),
		SMS:       "Hi {{first_name}}, we're sorry your {{procedure}} report is taking longer than usual. We're reviewing it to make sure it's right, and we'll tell you as soon as it's ready. Questions? Call {{phone}}.\n{{facility}}",
	},
	{
		ID: "missed_appointment", Name: "Missed appointment", Description: "A requested examination hasn't happened yet.",
		Channels:  []Channel{ChannelSMS},
		Variables: vars(vFirstName, vDoctor, vProcedure, vPhone),
		SMS:       "Hi {{first_name}}, a reminder about the {{procedure}} Dr. {{doctor}} requested for you. Please visit us when you can. To reschedule or for help, call {{phone}}.\n{{facility}}",
	},
	{
		ID: "whatsapp_failed", Name: "WhatsApp not delivered", Description: "A report sent by WhatsApp didn't arrive.",
		Channels:  []Channel{ChannelSMS},
		Variables: vars(vFirstName, vNumber, vReason, vPhone),
		SMS:       "Hi {{first_name}}, we tried to send your report to WhatsApp on {{whatsapp_number}}, but it couldn't be delivered because {{reason}}. Please call {{facility}} on {{phone}} and we'll get it to you another way.",
	},
}

// Message is a message by id, or nil.
func Message(id string) *MessageDef {
	for i := range Catalogue {
		if Catalogue[i].ID == id {
			return &Catalogue[i]
		}
	}
	return nil
}

// Has says whether the message is sent on a channel.
func (m *MessageDef) Has(c Channel) bool { return slices.Contains(m.Channels, c) }

// Allows says whether the message may say a variable.
func (m *MessageDef) Allows(name string) bool {
	for _, v := range m.Variables {
		if v.Name == name {
			return true
		}
	}
	return false
}

var (
	placeholder  = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_]*)\s*\}\}`)
	messagingSID = regexp.MustCompile(`^MG[0-9a-f]{32}$`)
	contentSID   = regexp.MustCompile(`^HX[0-9a-f]{32}$`)
)

// Placeholders are the variable names a text uses, in order, each once.
func Placeholders(text string) []string {
	var out []string
	for _, m := range placeholder.FindAllStringSubmatch(text, -1) {
		if !slices.Contains(out, m[1]) {
			out = append(out, m[1])
		}
	}
	return out
}

// validate checks the messages against the catalogue: known messages on their channels,
// only the variables each may say, and Twilio identifiers that look like Twilio's.
func (m *Messages) validate() []error {
	var errs []error
	bad := func(field, format string, args ...any) {
		errs = append(errs, fmt.Errorf("messages.%s: "+format, append([]any{field}, args...)...))
	}
	checkText := func(field string, def *MessageDef, text string) {
		for _, name := range Placeholders(text) {
			if !def.Allows(name) {
				bad(field, "{{%s}} isn't something this message can say", name)
			}
		}
	}
	for _, id := range sortedKeys(m.Email) {
		def := Message(id)
		if def == nil || !def.Has(ChannelEmail) {
			bad("email."+id, "isn't an email message")
			continue
		}
		c := m.Email[id]
		for field, text := range map[string]string{"subject": c.Subject, "heading": c.Heading, "body": c.Body, "button": c.Button} {
			checkText("email."+id+"."+field, def, text)
		}
		if c.Button != "" && def.Email.Button == "" {
			bad("email."+id+".button", "this email has no button")
		}
	}
	for _, id := range sortedKeys(m.SMS) {
		def := Message(id)
		if def == nil || !def.Has(ChannelSMS) {
			bad("sms."+id, "isn't an SMS message")
			continue
		}
		checkText("sms."+id, def, m.SMS[id])
	}
	if w := m.WhatsApp; w != nil {
		if w.MessagingServiceSID != "" && !messagingSID.MatchString(w.MessagingServiceSID) {
			bad("whatsapp.messaging_service_sid", "%q isn't a Twilio Messaging Service SID (MG and 32 hex digits)", w.MessagingServiceSID)
		}
		for _, id := range sortedKeys(w.Templates) {
			if def := Message(id); def == nil || !def.Has(ChannelWhatsApp) {
				bad("whatsapp.templates."+id, "isn't a WhatsApp message")
			} else if !contentSID.MatchString(w.Templates[id]) {
				bad("whatsapp.templates."+id, "%q isn't a Twilio Content SID (HX and 32 hex digits)", w.Templates[id])
			}
		}
	}
	return errs
}

// forPublish is what the control plane adds: a tenant's WhatsApp templates are sent from its
// own sender, never another tenant's.
func (m *Messages) forPublish() []error {
	if w := m.WhatsApp; w != nil && len(w.Templates) > 0 && w.MessagingServiceSID == "" {
		return []error{fmt.Errorf("messages.whatsapp.messaging_service_sid: the tenant's own WhatsApp sender is required to send its templates")}
	}
	return nil
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
