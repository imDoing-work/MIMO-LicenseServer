package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"mimo-license/internal/crypto"
	"mimo-license/internal/license"
)

func main() {
	// --------------------------------------------------
	// Load private key
	// --------------------------------------------------
	privKey, err := crypto.LoadPrivateKey("keys/private.pem")
	if err != nil {
		panic(err)
	}

	// --------------------------------------------------
	// Generate license UUID
	// --------------------------------------------------
	licenseUUID, err := license.NewUUIDv4()
	if err != nil {
		panic(err)
	}

	// --------------------------------------------------
	// Build payload (WITHOUT HardwareFP)
	// --------------------------------------------------
	payload := license.Payload{
		LicenseUUID: licenseUUID,
		Product:     "MIMO",
		Edition:     "Enterprise",

		IssuedAt: time.Now().UTC().Format(time.RFC3339),
		ExpireAt: time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339),

		Hardware: license.HardwareBind{
			BoardUUID: "4c4c4544-004a-5410-8058-b2c04f435732",

			MACs: []string{
				"50:6b:4b:01:7f:4e",
				"50:6b:4b:0d:24:82",
				"b8:2a:72:db:f3:a4",
				"b8:2a:72:db:f3:a5",
				"b8:2a:72:db:f3:a6",
				"b8:2a:72:db:f3:a7",
			},

			NvmeSerials: []string{
				"035AMGD9Q4001341",
				"035AMGD9Q4021724",
				"035HAXD9Q4019571",
				"035HAXD9Q4020228",
			},

			TotalMemoryKB: 394791660,
		},

		Features: license.Features{
			SuperBlock: true,
		},
	}

	// ==================================================
	// ① Canonical encode HardwareBind ONLY
	// ==================================================
	fpBind := license.HardwareFingerprintBind{
		BoardUUID:   payload.Hardware.BoardUUID,
		MACs:        payload.Hardware.MACs,
	}
	hwCanonical, err := license.EncodeHardwareFingerprintBindCanonical(fpBind)
	if err != nil {
		panic(err)
	}

	// ==================================================
	// ② sha256 + base64 → HardwareFP
	// ==================================================
	fp := sha256.Sum256(hwCanonical)
	payload.HardwareFP = base64.StdEncoding.EncodeToString(fp[:])

	// ==================================================
	// ③ Canonical encode FULL payload (WITH FP)
	// ==================================================
	canonicalPayload, err := license.EncodePayloadCanonical(payload)
	if err != nil {
		panic(err)
	}

	// ==================================================
	// ④ Sign payload
	// ==================================================
	signature, err := crypto.SignPayload(privKey, canonicalPayload)
	if err != nil {
		panic(err)
	}

	// --------------------------------------------------
	// Build final license
	// --------------------------------------------------
	lic := license.License{
		Payload:   payload,
		Signature: signature,
	}

	// --------------------------------------------------
	// Write license.json
	// --------------------------------------------------
	out, err := json.MarshalIndent(lic, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("license.json", out, 0644); err != nil {
		panic(err)
	}

	fmt.Println("License issued successfully")
	fmt.Println("License UUID:", payload.LicenseUUID)
}
