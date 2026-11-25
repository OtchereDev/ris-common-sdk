package letterhead

import (
	"fmt"
	"strings"
)

// Page geometry, as pdf-gen prints: A4, with a one-inch band at the top for the header and
// at the bottom for the footer, and no side margins (the furniture pads itself by 40px).
const (
	pageWidth    = "210mm"
	pageHeight   = "297mm"
	headerHeight = "1in"
	footerHeight = "1in"
)

// SampleBody is report text for a preview: enough to show where content sits against the
// furniture, with nothing that could be mistaken for a real patient's.
const SampleBody = `<h2 style="margin:0 0 6mm;font-size:15pt;">CT Chest (sample report)</h2>
<p><strong>Patient:</strong> Sample Patient &nbsp; <strong>Accession:</strong> SAMPLE-0001</p>
<p><strong>Clinical information:</strong> Persistent cough for three weeks.</p>
<p><strong>Findings:</strong> The lungs are clear. No consolidation, effusion or
pneumothorax. The heart and mediastinum are within normal limits. No lymphadenopathy.</p>
<p><strong>Impression:</strong> Normal study.</p>`

// Preview is one page as pdf-gen prints it, as a standalone HTML document: the header in the
// top band, the footer in the bottom one, the watermark behind body. It is the same markup the
// printer is given; only the page box around it is drawn here, where the printer draws it
// from its margins.
func (l Letterhead) Preview(body string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, `<!doctype html><html><head><meta charset="utf-8"><title>%s</title><style>
html,body{margin:0;background:#e9ebe6;}
.page{position:relative;box-sizing:border-box;width:%s;height:%s;margin:0 auto;background:#fff;
  display:flex;flex-direction:column;overflow:hidden;isolation:isolate;}
.band{flex:none;display:flex;align-items:center;}
.header{height:%s;}
.footer{height:%s;}
.content{flex:1;position:relative;padding:0 40px;font-family:Helvetica,Arial,sans-serif;
  font-size:11pt;line-height:1.5;color:#191d17;}
</style></head><body><div class="page">`,
		htmlEscape(l.Name), pageWidth, pageHeight, headerHeight, footerHeight)
	fmt.Fprintf(&sb, `<div class="band header">%s</div>`, l.HeaderHTML())
	fmt.Fprintf(&sb, `<div class="content">%s%s</div>`, l.WatermarkHTML(), body)
	fmt.Fprintf(&sb, `<div class="band footer">%s</div>`, l.FooterHTML())
	sb.WriteString(`</div></body></html>`)
	return sb.String()
}

func htmlEscape(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}
