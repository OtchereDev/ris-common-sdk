package tenant_test

import (
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
	"github.com/OtchereDev/ris-common-sdk/pkg/tenant/tenanttest"
)

func TestAssetKey(t *testing.T) {
	sum := sha256.Sum256([]byte("a logo"))
	key := tenant.AssetKey("demo", tenant.ImageLogoFull, sum[:], "png")
	if !strings.HasPrefix(key, "tenants/demo/logo.full-") || !strings.HasSuffix(key, ".png") || len(key) != len("tenants/demo/logo.full-")+16+4 {
		t.Fatalf("key: %s", key)
	}
	if err := tenant.ParseAssetKey("demo", tenant.ImageLogoFull, key); err != nil {
		t.Fatalf("its own key: %v", err)
	}
	for _, c := range []struct {
		tenant, key, want string
		slot              tenant.ImageSlot
	}{
		{"quest", key, "another tenant's", tenant.ImageLogoFull},
		{"demo", key, "not document.logo", tenant.ImageDocumentLogo},
		{"demo", "https://cdn.example.com/logo.png", "not an image key", tenant.ImageLogoFull},
		{"demo", strings.Replace(key, ".png", ".exe", 1), ".exe is not", tenant.ImageLogoFull},
	} {
		err := tenant.ParseAssetKey(c.tenant, c.slot, c.key)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s as %s for %s: %v, want %q", c.key, c.slot, c.tenant, err, c.want)
		}
	}
}

func TestImagesInARecord(t *testing.T) {
	sum := sha256.Sum256([]byte("a logo"))
	key := tenant.AssetKey("demo", tenant.ImageDocumentLogo, sum[:], "png")
	inline := "data:image/png;base64," + strings.Repeat("iVBORw0KGgo/", 30)

	withImage := func(v string) *tenant.Record {
		r := tenanttest.Record()
		r.Profile.SetImage(tenant.ImageDocumentLogo, v)
		return &r
	}

	// A key: valid, and publishable.
	if err := withImage(key).ValidateForPublish(nil); err != nil {
		t.Errorf("a record with a key: %v", err)
	}
	// Inline artwork: a service still reads it, but it can't be published again.
	r := withImage(inline)
	if err := r.Validate(); err != nil {
		t.Errorf("a service must still read an older record with an inline image: %v", err)
	}
	if err := r.ValidateForPublish(nil); err == nil || !strings.Contains(err.Error(), "profile.document.logo: is an image inside the record") {
		t.Errorf("publishing an inline image: %v", err)
	}
	// Anything else in a slot is a mistake.
	if err := withImage("logo.png").Validate(); err == nil || !strings.Contains(err.Error(), "profile.document.logo:") {
		t.Errorf("a slot holding neither: %v", err)
	}
	// Emptying a slot clears it.
	r.Profile.SetImage(tenant.ImageDocumentLogo, "")
	if r.Profile.Document.Logo != nil || r.Profile.Image(tenant.ImageDocumentLogo) != "" {
		t.Errorf("emptied slot: %v", r.Profile.Document.Logo)
	}
}
