package license

import (
	"bytes"
	"encoding/json"
	"sort"
)

//
// ============================================================
// EncodePayloadCanonical
// ============================================================
//
// Canonical rules (MUST stay stable):
// 1. Go struct field order defines JSON key order
// 2. All slices are normalized (sorted) before encoding
// 3. No map types are allowed in Payload
// 4. No indentation, no extra whitespace
//
// ⚠️ Payload 是「被签名内容」
// ⚠️ 这里可以包含 TotalMemoryKB，但它【不参与指纹】
//
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
	// ---------- Hardware (record / signed only) ----------
	sort.Strings(p.Hardware.MACs)
	sort.Strings(p.Hardware.NvmeSerials)

	// ⚠️ TotalMemoryKB 保留，不参与 FP

	// ---------- Future-proof ----------
	// Add normalization here if new slices appear

	return p
}

//
// ============================================================
// EncodeHardwareFingerprintBindCanonical
// ============================================================
//
// Canonical rules (MUST stay stable):
// 1. Go struct field order defines JSON key order
// 2. All slices are normalized (sorted) before encoding
// 3. No map types are allowed
// 4. No indentation, no extra whitespace
//
// 🔐 ONLY used for HardwareFP calculation
// 🔐 MUST NOT include volatile fields (memory, cpu, etc.)
//
func EncodeHardwareFingerprintBindCanonical(h HardwareFingerprintBind) ([]byte, error) {
	normalized := normalizeHardwareFingerprintBind(h)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(normalized); err != nil {
		return nil, err
	}

	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func normalizeHardwareFingerprintBind(h HardwareFingerprintBind) HardwareFingerprintBind {
	sort.Strings(h.MACs)
	return h
}
