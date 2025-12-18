package license

type HardwareBind struct {
	// --- 唯一性 / 身份类 ---
	BoardUUID string   `json:"board_uuid"`
	MACs      []string `json:"mac_addresses"`
	

	// --- 设备厂商 / 类型约束 ---
	NvmeVendorOUIs []string `json:"nvme_vendors"`

	// --- 数量 / 规模类约束（比对型） ---
	TotalNvmeCap  uint64 `json:"total_nvme_cap"`   // bytes or MB，需统一单位
	NvmeCount     int    `json:"nvme_count"`
	TotalMemoryKB uint64 `json:"total_memory_kb"`
}


type Features struct {
	SuperBlock     bool `json:"SuperBlock"`
}

type Payload struct {
	LicenseUUID string `json:"license_uuid"`
	Product     string `json:"product"`
	Edition     string `json:"edition"`

	IssuedAt string `json:"issued_at"`
	ExpireAt string `json:"expire_at"`

	Hardware HardwareBind `json:"hardware_bind"`
	Features Features     `json:"features"`
}

type License struct {
	Payload   Payload `json:"payload"`
	Signature string  `json:"signature"`
}
