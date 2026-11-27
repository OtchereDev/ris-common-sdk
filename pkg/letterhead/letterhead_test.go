package letterhead_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/OtchereDev/ris-common-sdk/pkg/letterhead"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant/tenanttest"
)

// fakeImages turns a key into a recognisable data URI, and records what it was asked for.
func fakeImages(asked *[]string) letterhead.Images {
	return func(_ context.Context, v string) (string, error) {
		*asked = append(*asked, v)
		return "data:image/png;base64,FOR-" + v, nil
	}
}

func ptr[T any](v T) *T { return &v }

func TestFromProfile(t *testing.T) {
	p := tenanttest.Record().Profile
	p.Name = "Demo RIS"
	p.Site = "www.demo.example"
	p.Logo = tenant.Logo{Full: "tenants/demo/logo.full-1111111111111111.png", Mark: "tenants/demo/logo.mark-2222222222222222.png"}
	p.Locale.CurrencyPrefix = ""

	var asked []string
	l, err := letterhead.FromProfile(context.Background(), &p, fakeImages(&asked))
	if err != nil {
		t.Fatal(err)
	}
	// The profile's one identity serves the page where the document overrides nothing.
	if l.Name != "Demo RIS" || l.Logo != "data:image/png;base64,FOR-"+p.Logo.Full {
		t.Errorf("name %q, logo %q", l.Name, l.Logo)
	}
	// The mark is not a watermark: one is drawn only when the profile sets it.
	if l.Watermark != "" {
		t.Errorf("a watermark without one set: %q", l.Watermark)
	}
	if l.WatermarkOpacity != letterhead.DefaultWatermarkOpacity || l.FooterColor != letterhead.DefaultFooterColor {
		t.Errorf("defaults: %v %q", l.WatermarkOpacity, l.FooterColor)
	}
	if l.Currency != "GHS" {
		t.Errorf("currency from the locale: %q", l.Currency)
	}

	// The document's own name, logo, watermark and colours win.
	p.Document = tenant.Document{
		Name:             ptr("Demo Diagnostics"),
		Logo:             ptr("tenants/demo/document.logo-3333333333333333.png"),
		Watermark:        ptr("tenants/demo/document.watermark-4444444444444444.png"),
		WatermarkOpacity: ptr(0.1),
		FooterColor:      "#0b76a8",
	}
	l, _ = letterhead.FromProfile(context.Background(), &p, fakeImages(&asked))
	if l.Name != "Demo Diagnostics" || !strings.HasSuffix(l.Logo, "document.logo-3333333333333333.png") ||
		!strings.HasSuffix(l.Watermark, "document.watermark-4444444444444444.png") || l.WatermarkOpacity != 0.1 || l.FooterColor != "#0b76a8" {
		t.Errorf("document overrides: %+v", l)
	}

	// No artwork at all asks for nothing and prints no header.
	asked = nil
	bare := tenanttest.Record().Profile
	l, _ = letterhead.FromProfile(context.Background(), &bare, fakeImages(&asked))
	if len(asked) != 0 || l.HeaderHTML() != "" || l.WatermarkHTML() != "" {
		t.Errorf("no artwork: asked %v, header %q", asked, l.HeaderHTML())
	}
}

func TestFooter(t *testing.T) {
	l := letterhead.Letterhead{
		Site:        "www.pluslabgh.com",
		FooterColor: "#0b76a8",
		Contact: [][]tenant.ContactItem{
			{{Icon: "location", Text: "Near Bank of Africa, Comm. 21 Annex"}, {Icon: "phone", Text: "0303 319509"}},
			{{Text: "Lab <&> Co"}},
		},
	}
	f := l.FooterHTML()
	for _, want := range []string{"color:#0b76a8", "Near Bank of Africa", "<svg", "width:8px", "Lab &lt;&amp;&gt; Co", "www.pluslabgh.com"} {
		if !strings.Contains(f, want) {
			t.Errorf("footer is missing %q:\n%s", want, f)
		}
	}
	if strings.Count(f, "<svg") != 2 {
		t.Errorf("an item without an icon got one:\n%s", f)
	}
	if got := l.AddressLines(); len(got) != 2 || got[0] != "Near Bank of Africa, Comm. 21 Annex  0303 319509" {
		t.Errorf("address lines: %q", got)
	}

	// A footer brought whole is printed untouched.
	raw := `<div style="font-size: 12px"><p style="color: #e76f0f">No. 12 Opoku Adjei Avenue</p></div>`
	l.FooterRaw = raw
	if l.FooterHTML() != raw {
		t.Errorf("raw footer changed: %s", l.FooterHTML())
	}
	if (letterhead.Letterhead{}).FooterHTML() != "" {
		t.Error("an empty letterhead has a footer")
	}
}

func TestPreview(t *testing.T) {
	l := letterhead.Letterhead{
		Name: `Quest "Medical"`, Logo: "data:image/png;base64,LOGO", Watermark: "data:image/png;base64,MARK",
		WatermarkOpacity: 0.06, Contact: [][]tenant.ContactItem{{{Text: "Accra"}}},
	}
	page := l.Preview(letterhead.SampleBody)
	for _, want := range []string{"<!doctype html>", "width:210mm", "height:1in", "base64,LOGO", "max-height:0.8in", "base64,MARK", "opacity:0.06", "Accra", "sample report", "Quest &quot;Medical&quot;"} {
		if !strings.Contains(page, want) {
			t.Errorf("preview is missing %q", want)
		}
	}
}

func TestFromBucket(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if r.URL.Path != "/tenants/demo/logo.full-1111111111111111.svg" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("<svg/>"))
	}))
	defer srv.Close()
	images := letterhead.FromBucket(srv.URL+"/", srv.Client())
	ctx := context.Background()

	for range 2 {
		uri, err := images(ctx, "tenants/demo/logo.full-1111111111111111.svg")
		if err != nil || uri != "data:image/svg+xml;base64,PHN2Zy8+" {
			t.Fatalf("%q, %v", uri, err)
		}
	}
	if hits.Load() != 1 {
		t.Errorf("fetched %d times; a key names one file, so once is enough", hits.Load())
	}
	if _, err := images(ctx, "tenants/demo/logo.mark-2222222222222222.png"); err == nil {
		t.Error("a missing image should be an error")
	}
	if uri, _ := images(ctx, "iVBORw0KGgo"); uri != "data:image/png;base64,iVBORw0KGgo" {
		t.Errorf("bare base64 from an older record: %q", uri)
	}
}
