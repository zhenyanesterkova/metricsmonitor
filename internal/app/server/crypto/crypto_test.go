package crypto

import (
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateKeyPair(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-keys")
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		require.NoError(t, err)
	}()

	privatePath := filepath.Join(tempDir, "private.crt")
	pubPath := filepath.Join(tempDir, "pub.crt")

	err = GenerateKeyPair(privatePath, pubPath)
	require.NoError(t, err)

	t.Run("access rights and file existence", func(t *testing.T) {
		info, err := os.Stat(privatePath)
		require.NoError(t, err)
		//nolint:all //checking file access rights
		require.Equal(t, true, (info.Mode()&0600) == filePermission)

		info, err = os.Stat(pubPath)
		require.NoError(t, err)
		//nolint:all //checking file access rights
		require.Equal(t, true, (info.Mode()&0600) == filePermission)
	})

	t.Run("file contents", func(t *testing.T) {
		privateData, err := os.ReadFile(privatePath)
		require.NoError(t, err)
		publicData, err := os.ReadFile(pubPath)
		require.NoError(t, err)

		block, _ := pem.Decode(privateData)
		require.Equal(t, "RSA PRIVATE KEY", block.Type)

		block, _ = pem.Decode(publicData)
		require.Equal(t, "RSA PUBLIC KEY", block.Type)
	})

	t.Run("empty path", func(t *testing.T) {
		err := GenerateKeyPair("", "")
		require.Error(t, err)
	})
}

func Test_checkKeyFileExists(t *testing.T) {
	exists, err := checkKeyFileExists("./notexists")
	require.NoError(t, err)
	require.Equal(t, false, exists)

	tempDir, err := os.MkdirTemp("", "exists")
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		require.NoError(t, err)
	}()

	exists, err = checkKeyFileExists(tempDir)
	require.NoError(t, err)
	require.Equal(t, true, exists)
}
