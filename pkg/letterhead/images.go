package letterhead

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/OtchereDev/ris-common-sdk/pkg/tenant"
)

// FromBucket resolves images from the asset bucket, read at baseURL (PLATFORM_ASSETS_URL,
// e.g. https://assets.spectraris.com). A key names exactly one file, so each is fetched once
// and kept.
func FromBucket(baseURL string, client *http.Client) Images {
	if client == nil {
		client = http.DefaultClient
	}
	var mu sync.Mutex
	cache := map[string]string{}
	return func(ctx context.Context, value string) (string, error) {
		if !strings.HasPrefix(value, "tenants/") {
			return Inline(value), nil
		}
		mu.Lock()
		uri, ok := cache[value]
		mu.Unlock()
		if ok {
			return uri, nil
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/"+value, nil)
		if err != nil {
			return "", err
		}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%s: %s", value, resp.Status)
		}
		data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		if err != nil {
			return "", err
		}
		uri = DataURI(value, data)
		mu.Lock()
		cache[value] = uri
		mu.Unlock()
		return uri, nil
	}
}

// DataURI encodes an image stored under key, its type taken from the key's extension.
func DataURI(key string, data []byte) string {
	mediaType := tenant.ImageTypes[strings.TrimPrefix(path.Ext(key), ".")]
	if mediaType == "" {
		mediaType = "image/png"
	}
	return "data:" + mediaType + ";base64," + base64.StdEncoding.EncodeToString(data)
}

// Inline is artwork carried in a record made before keys, as a data: URI: as it is if it is
// one already, otherwise bare base64 taken to be a PNG.
func Inline(value string) string {
	if strings.HasPrefix(value, "data:") {
		return value
	}
	return "data:image/png;base64," + value
}
