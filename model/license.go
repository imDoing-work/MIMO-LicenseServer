// model/license.go
package model

// LicenseType 用于区分本地或在线模式
type LicenseType string

const (
    LicenseTypeLocal  LicenseType = "local"
    LicenseTypeOnline LicenseType = "online"
)

// License 是客户端和服务端共享的许可证数据结构
type License struct {
    LicenseID       string      `json:"lid"`
    HWID            string      `json:"hwid"`
    ClientID        string      `json:"client_id"`
    Level           string      `json:"level"`
    Features        []string    `json:"features"`
    Expires         string      `json:"expires"`
    Type            LicenseType `json:"type"`
    MaxNvmeCapacity uint64      `json:"max_nvme_capacity"`
    
    // 签名字段必须在最后，因为它不参与签名内容计算
    Signature string `json:"signature"`
}