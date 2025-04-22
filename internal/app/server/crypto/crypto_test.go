package crypto

import (
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateKeyPair(t *testing.T) {
	privateKeyDir := "example-private.crt"
	publicKeyDir := "example-public.crt"

	err := GenerateKeyPair(privateKeyDir, publicKeyDir)
	require.NoError(t, err)

	t.Run("access rights and file existence", func(t *testing.T) {
		info, err := os.Stat(privateKeyDir)
		require.NoError(t, err)
		//nolint:all //checking file access rights
		require.Equal(t, true, (info.Mode()&0600) == filePermission)

		info, err = os.Stat(publicKeyDir)
		require.NoError(t, err)
		//nolint:all //checking file access rights
		require.Equal(t, true, (info.Mode()&0600) == filePermission)
	})

	t.Run("file contents", func(t *testing.T) {
		privateData, err := os.ReadFile(privateKeyDir)
		require.NoError(t, err)
		publicData, err := os.ReadFile(publicKeyDir)
		require.NoError(t, err)

		block, _ := pem.Decode(privateData)
		require.Equal(t, "RSA PRIVATE KEY", block.Type)

		block, _ = pem.Decode(publicData)
		require.Equal(t, "RSA PUBLIC KEY", block.Type)
	})

	err = os.Remove(privateKeyDir)
	require.NoError(t, err)
	err = os.Remove(publicKeyDir)
	require.NoError(t, err)

	t.Run("empty path", func(t *testing.T) {
		err := GenerateKeyPair("", "")
		require.Error(t, err)
	})
}

func Test_checkKeyFileExists(t *testing.T) {
	exists, err := checkExists("./notexists")
	require.NoError(t, err)
	require.Equal(t, false, exists)

	tempDir, err := os.MkdirTemp("", "exists")
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		require.NoError(t, err)
	}()

	exists, err = checkExists(tempDir)
	require.NoError(t, err)
	require.Equal(t, true, exists)
}
