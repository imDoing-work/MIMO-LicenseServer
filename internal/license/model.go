package license

// ============================================================
// 🔐 仅用于【硬件指纹计算】的绑定字段
// ============================================================
type HardwareFingerprintBind struct {
	BoardUUID   string   `json:"board_uuid"`
	MACs        []string `json:"mac_addresses"`
}

// ============================================================
// 📦 用于【记录 / 展示 / 签名】的完整硬件信息
// ============================================================
type HardwareBind struct {
	// --- 唯一性 / 身份类 ---
	BoardUUID   string   `json:"board_uuid"`
	MACs        []string `json:"mac_addresses"`
	NvmeSerials []string `json:"nvme_serials"`

	// --- 仅记录，不参与指纹 ---
	TotalMemoryKB uint64 `json:"total_memory_kb"`
}

// ============================================================
// 🎛 Feature 开关
// ============================================================
type Features struct {
	SuperBlock bool `json:"SuperBlock"`
}

// ============================================================
// 📄 License Payload（被签名的主体）
// ============================================================
type Payload struct {
	LicenseUUID string `json:"license_uuid"`
	Product     string `json:"product"`
	Edition     string `json:"edition"`

	IssuedAt string `json:"issued_at"`
	ExpireAt string `json:"expire_at"`

	// 📦 完整硬件信息（可展示、可签名）
	Hardware HardwareBind `json:"hardware_bind"`

	// 🔐 硬件指纹（只由 HardwareFingerprintBind 算出）
	HardwareFP string `json:"hardware_fp"`

	Features Features `json:"features"`
}

// ============================================================
// 🧾 最终 License 文件
// ============================================================
type License struct {
	Payload   Payload `json:"payload"`
	Signature string  `json:"signature"`
}
