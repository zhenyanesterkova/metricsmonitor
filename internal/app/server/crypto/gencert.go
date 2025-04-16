package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	filePermission = 0o600
)

func GenerateKeyPair(pathToPrivate, pathToPublic string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed generate private key: %w", err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile(pathToPrivate, privateKeyPEM, filePermission)
	if err != nil {
		return fmt.Errorf("failed write private key in file: %w", err)
	}
	err = os.WriteFile(pathToPublic, publicKeyPEM, filePermission)
	if err != nil {
		return fmt.Errorf("failed write public key: %w", err)
	}

	return nil
}
