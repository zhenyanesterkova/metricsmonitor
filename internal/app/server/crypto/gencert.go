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
	existsPrivate, err := checkExists(pathToPrivate)
	if err != nil {
		return fmt.Errorf("failed check exists file with private key: %w", err)
	}
	existsPublic, err := checkExists(pathToPublic)
	if err != nil {
		return fmt.Errorf("failed check exists file with public key: %w", err)
	}
	if existsPrivate && existsPublic {
		return nil
	}

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

	filePrivate, err := os.OpenFile(pathToPrivate, os.O_CREATE, filePermission)
	if err != nil {
		return fmt.Errorf("failed open file for write private key: %w", err)
	}
	filePublic, err := os.OpenFile(pathToPublic, os.O_CREATE, filePermission)
	if err != nil {
		return fmt.Errorf("failed open file for write public key: %w", err)
	}

	_, err = filePrivate.Write(privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed write private key to file: %w", err)
	}
	_, err = filePublic.Write(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed write public key to file: %w", err)
	}

	err = filePrivate.Close()
	if err != nil {
		return fmt.Errorf("failed close private key file: %w", err)
	}
	err = filePublic.Close()
	if err != nil {
		return fmt.Errorf("failed close public key file: %w", err)
	}

	return nil
}

func checkExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed get info about file: %w", err)
	}

	return true, nil
}
