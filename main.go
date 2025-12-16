// main.go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"mimo-license/model"
)

const (
    PrivateKeyPath = "./keys/private.pem" // 服务器私钥
    PublicKeyPath  = "./keys/public.pem"  // 客户端公钥
    OutputFilePath = "signed_license.json"
)

func main() {
    if err := run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}

func run() error {
    // 1. 演示生成一个 License Payload (模拟从数据库读取的授权记录)
    // 注意：HWID 应是客户机器的实际硬件指纹
    licensePayload := model.License{
        LicenseID: "LID-DEMO-001",
        HWID: "4c4c4544-0046-5110-8041-b6c04f564331", // 假设的目标 HWID
        ClientID: "MIMO-ENTERPRISE-01",
        Level: "Premium",
        Features: []string{"raid:level_6", "compression:enabled", "support:24x7"},
        Expires: time.Now().AddDate(1, 0, 0).Format("2006-01-02"), // 一年后过期
        Type: model.LicenseTypeOnline,
        MaxNvmeCapacity: 256 * 1024 * 1024 * 1024 * 1024, // 256TB
    }
    fmt.Println("--- 1. License Payload Constructed ---")
    
    // 2. 签名
    signature, err := SignLicense(licensePayload, PrivateKeyPath)
    if err != nil {
        return fmt.Errorf("failed to sign license: %w", err)
    }
    licensePayload.Signature = signature
    fmt.Printf("--- 2. License Signed Successfully ---\nSignature: %s...\n", signature[:10])

    // 3. 写入文件 (这是客户收到的文件)
    signedData, err := json.MarshalIndent(licensePayload, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal signed license: %w", err)
    }
    if err := os.WriteFile(OutputFilePath, signedData, 0644); err != nil {
        return fmt.Errorf("failed to save signed license: %w", err)
    }
    fmt.Printf("--- 3. Signed License Saved to %s ---\n", OutputFilePath)

    // 4. 演示验证 (客户端校验逻辑)
    fmt.Println("\n--- 4. Client Side Verification Demo ---")
    isValid, err := VerifyLicenseSignature(licensePayload, PublicKeyPath)
    if err != nil {
        fmt.Printf("Verification failed: %v\n", err)
    } else if isValid {
        fmt.Println("✅ Verification successful! The license is authentic and untampered.")
    }

    // 5. 演示篡改失败
    fmt.Println("\n--- 5. Tamper Test Demo ---")
    licensePayload.Expires = "2099-12-31" // 尝试延长有效期
    tamperCheck, err := VerifyLicenseSignature(licensePayload, PublicKeyPath)
    if err == nil && tamperCheck {
        fmt.Printf("❌ Tamper check FAILED! Still valid: %v\n", tamperCheck)
    } else {
        fmt.Println("✅ Tamper check successful! Verification failed after modification.")
    }

    return nil
}

// --------------------------------------------------------------------------------
// ⚠️ 注意：您需要手动添加一个 Key 生成函数或使用外部工具生成 keys/private.pem 和 keys/public.pem
// 示例 Key 生成代码 (可放在 crypto.go 或独立文件运行一次)
func generateKeys() {
    privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
    
    // 保存私钥 (PKCS#1 PEM format)
    privateKeyPEM := &pem.Block{
        Type: "RSA PRIVATE KEY", 
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    }
    os.WriteFile(PrivateKeyPath, pem.EncodeToMemory(privateKeyPEM), 0600)
    
    // 保存公钥 (PKIX PEM format)
    publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
    publicKeyPEM := &pem.Block{
        Type: "PUBLIC KEY",
        Bytes: publicKeyBytes,
    }
    os.WriteFile(PublicKeyPath, pem.EncodeToMemory(publicKeyPEM), 0644)
}
// --------------------------------------------------------------------------------