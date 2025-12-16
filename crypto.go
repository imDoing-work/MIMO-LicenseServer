package main

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "fmt"
    "os"
    
    "mimo-license/model"
    "encoding/json"
)

// --- 密钥加载辅助函数 ---

// loadPrivateKey 从 PEM 文件加载私钥
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read private key file: %w", err)
    }
    block, _ := pem.Decode(data)
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }
    key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %w", err)
    }
    return key, nil
}

// loadPublicKey 从 PEM 文件加载公钥
func loadPublicKey(path string) (*rsa.PublicKey, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read public key file: %w", err)
    }
    block, _ := pem.Decode(data)
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }
    key, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %w", err)
    }
    rsaKey, ok := key.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("key is not RSA public key")
    }
    return rsaKey, nil
}

// --- 签名/验证核心逻辑 ---

// getCanonicalBytes 返回用于签名的规范化 JSON 字节数组
func getCanonicalBytes(lic model.License) ([]byte, error) {
    // 1. 清空签名 (签名自身不参与签名计算)
    tempLic := lic
    tempLic.Signature = "" 
    
    // 2. 将结构体转换为 map，再转回 JSON，以确保键名按字典序排序
    data, err := json.Marshal(tempLic)
    if err != nil {
        return nil, err
    }
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return nil, err
    }
    // 再次 Marshal，Go 会自动按字典序排序键，实现规范化
    canonicalBytes, err := json.Marshal(raw) 
    if err != nil {
        return nil, err
    }
    return canonicalBytes, nil
}

// SignLicense 使用 RSA 私钥对许可证进行签名
func SignLicense(lic model.License, privateKeyPath string) (string, error) {
    // 1. 加载私钥
    privateKey, err := loadPrivateKey(privateKeyPath)
    if err != nil {
        return "", err
    }

    // 2. 获取规范化字节
    canonicalBytes, err := getCanonicalBytes(lic)
    if err != nil {
        return "", fmt.Errorf("failed to get canonical bytes: %w", err)
    }

    // 3. 计算 SHA256 哈希
    hashed := sha256.Sum256(canonicalBytes)
    
    // 4. 使用私钥签名 (RSASSA-PKCS1-V1_5)
    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
    if err != nil {
        return "", fmt.Errorf("failed to sign payload: %w", err)
    }

    // 5. 编码为 Base64 字符串
    return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyLicenseSignature 使用 RSA 公钥验证许可证签名
func VerifyLicenseSignature(lic model.License, publicKeyPath string) (bool, error) {
    // 1. 加载公钥
    publicKey, err := loadPublicKey(publicKeyPath)
    if err != nil {
        return false, err
    }

    // 2. 获取规范化字节
    canonicalBytes, err := getCanonicalBytes(lic)
    if err != nil {
        return false, fmt.Errorf("failed to get canonical bytes: %w", err)
    }

    // 3. 提取签名
    signatureBytes, err := base64.StdEncoding.DecodeString(lic.Signature)
    if err != nil {
        return false, fmt.Errorf("failed to decode signature: %w", err)
    }
    
    // 4. 计算 Payload 的 SHA256 哈希
    hashed := sha256.Sum256(canonicalBytes)

    // 5. 使用公钥验证签名
    // VerifyPKCS1v15 只在验证失败时返回 error，成功返回 nil
    err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signatureBytes)
    if err != nil {
        return false, fmt.Errorf("signature verification failed: %w", err)
    }

    return true, nil
}