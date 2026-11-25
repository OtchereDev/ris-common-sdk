package tenant

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// A tenant's artwork lives in platform's asset bucket, which anyone can read. The record holds
// each image's key there, never the image itself: tenants/<tenant>/<slot>-<sha256 prefix>.<ext>.
// The hash makes a key name exactly one file, so a new logo is a new key and an address built
// from a key can be cached forever.
//
// Records published before this held the artwork inline, as a data: URI or bare base64.
// Validate still accepts those, so a service never refuses a record it already has;
// ValidateForPublish refuses new ones.

// ImageSlot names a place in the profile that holds an image.
type ImageSlot string

const (
	ImageLogoFull          ImageSlot = "logo.full"
	ImageLogoMark          ImageSlot = "logo.mark"
	ImageLogoFavicon       ImageSlot = "logo.favicon"
	ImageDocumentLogo      ImageSlot = "document.logo"
	ImageDocumentWatermark ImageSlot = "document.watermark"
)

// ImageSlots is every slot, in the order the Console shows them.
var ImageSlots = []ImageSlot{
	ImageLogoFull, ImageLogoMark, ImageLogoFavicon, ImageDocumentLogo, ImageDocumentWatermark,
}

// ImageTypes are the file types an image may be, by the extension its key ends in.
var ImageTypes = map[string]string{
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"webp": "image/webp",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
}

// assetHashLen is how much of the SHA-256 a key carries: 16 hex digits, 64 bits.
const assetHashLen = 16

var assetKeyPattern = regexp.MustCompile(`^tenants/([a-z][a-z0-9]{1,30})/([a-z]+\.[a-z]+)-([0-9a-f]{16})\.([a-z]+)$`)

// Valid reports whether s is one of the slots.
func (s ImageSlot) Valid() bool {
	for _, slot := range ImageSlots {
		if s == slot {
			return true
		}
	}
	return false
}

// AssetKey is the key an image is stored under. sum is the file's SHA-256 and ext one of
// ImageTypes.
func AssetKey(tenantKey string, slot ImageSlot, sum []byte, ext string) string {
	return fmt.Sprintf("tenants/%s/%s-%s.%s", tenantKey, slot, hex.EncodeToString(sum)[:assetHashLen], ext)
}

// ParseAssetKey checks that key is an asset key for the given tenant and slot.
func ParseAssetKey(tenantKey string, slot ImageSlot, key string) error {
	m := assetKeyPattern.FindStringSubmatch(key)
	switch {
	case m == nil:
		return fmt.Errorf("%q is not an image key (tenants/<tenant>/<slot>-<hash>.<ext>)", key)
	case m[1] != tenantKey:
		return fmt.Errorf("%q is another tenant's image", key)
	case m[2] != string(slot):
		return fmt.Errorf("%q is the image for %s, not %s", key, m[2], slot)
	case ImageTypes[m[4]] == "":
		return fmt.Errorf("%q: .%s is not an image type platform stores", key, m[4])
	}
	return nil
}

// IsInlineImage reports whether v is artwork carried in the record itself, as records did
// before images moved to the asset bucket.
func IsInlineImage(v string) bool {
	return strings.HasPrefix(v, "data:") || (len(v) > 200 && !strings.HasPrefix(v, "tenants/"))
}

// Image is what a slot holds: a key, an inline image in an older record, or "".
func (p *Profile) Image(slot ImageSlot) string {
	switch slot {
	case ImageLogoFull:
		return p.Logo.Full
	case ImageLogoMark:
		return p.Logo.Mark
	case ImageLogoFavicon:
		return p.Logo.Favicon
	case ImageDocumentLogo:
		return deref(p.Document.Logo)
	case ImageDocumentWatermark:
		return deref(p.Document.Watermark)
	}
	return ""
}

// SetImage puts v in a slot; "" empties it.
func (p *Profile) SetImage(slot ImageSlot, v string) {
	var ptr *string
	if v != "" {
		ptr = &v
	}
	switch slot {
	case ImageLogoFull:
		p.Logo.Full = v
	case ImageLogoMark:
		p.Logo.Mark = v
	case ImageLogoFavicon:
		p.Logo.Favicon = v
	case ImageDocumentLogo:
		p.Document.Logo = ptr
	case ImageDocumentWatermark:
		p.Document.Watermark = ptr
	}
}

// validateImages is part of Validate: each slot is empty, a key for this tenant and slot, or
// (in a record published before keys) an inline image.
func (p *Profile) validateImages(tenantKey string) []error {
	var errs []error
	for _, slot := range ImageSlots {
		v := p.Image(slot)
		if v == "" || IsInlineImage(v) {
			continue
		}
		if err := ParseAssetKey(tenantKey, slot, v); err != nil {
			errs = append(errs, fmt.Errorf("profile.%s: %w", slot, err))
		}
	}
	return errs
}

// inlineImages is part of ValidateForPublish: a new record carries keys only.
func (p *Profile) inlineImages() []error {
	var errs []error
	for _, slot := range ImageSlots {
		if IsInlineImage(p.Image(slot)) {
			errs = append(errs, fmt.Errorf("profile.%s: is an image inside the record; upload it in the Console, which stores it and puts its key here", slot))
		}
	}
	return errs
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
