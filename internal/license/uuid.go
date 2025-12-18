package license

import (
	"crypto/rand"
	"fmt"
)

// NewUUIDv4 生成 RFC4122 标准的 UUID v4
// 不依赖任何第三方包
func NewUUIDv4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// RFC 4122 variant
	b[8] = b[8]&0xbf | 0x80
	// Version 4
	b[6] = b[6]&0x0f | 0x40

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}
