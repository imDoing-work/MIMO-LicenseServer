package license

import (
	"bytes"
	"encoding/json"
	"sort"
)

// EncodePayloadCanonical
//
// Canonical rules (MUST stay stable):
// 1. Go struct field order defines JSON key order
// 2. All slices are normalized (sorted) before encoding
// 3. No map types are allowed in Payload
// 4. No indentation, no extra whitespace
//
// This guarantees identical byte output across
// issuing and verification sides.
func EncodePayloadCanonical(p Payload) ([]byte, error) {
	normalized := normalizePayload(p)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	// ❗ no SetIndent — whitespace must be deterministic

	if err := enc.Encode(normalized); err != nil {
		return nil, err
	}

	// Encoder always appends '\n'
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func normalizePayload(p Payload) Payload {
	// ---------- Hardware ----------
	sort.Strings(p.Hardware.MACs)
	sort.Strings(p.Hardware.NvmeVendorOUIs)

	// ---------- Future-proof ----------
	// If you later add more slices, normalize them here
	//
	// Example:
	// sort.Strings(p.SomeFutureField)

	return p
}
