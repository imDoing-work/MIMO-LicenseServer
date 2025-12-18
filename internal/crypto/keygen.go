package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateRSAKeyPair(bits int, privPath, pubPath string) error {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privFile, err := os.Create(privPath)
	if err != nil {
		return err
	}
	defer privFile.Close()

	pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	pubFile, err := os.Create(pubPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pem.Encode(pubFile, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})

	return nil
}
