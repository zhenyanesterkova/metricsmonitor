package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetTestFlags(c *Config) (*flags, error) {
	adress := ""
	flag.StringVar(
		&adress,
		"a",
		adress,
		"address and port to run server",
	)
	err := flag.CommandLine.Set("a", "www.testfromflag")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -a: %w", err)
	}

	config := ""
	flag.StringVar(
		&config,
		"c",
		config,
		"config name",
	)
	err = flag.CommandLine.Set("c", "config.json")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -c: %w", err)
	}

	logLevel := ""
	flag.StringVar(
		&logLevel,
		"l",
		logLevel,
		"log level",
	)
	err = flag.CommandLine.Set("l", "levelfromflag")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -l: %w", err)
	}

	var tempDur int
	flag.IntVar(
		&tempDur,
		"i",
		300,
		"store interval",
	)
	err = flag.CommandLine.Set("i", "500")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -i: %w", err)
	}

	fileStoragePath := ""
	flag.StringVar(
		&fileStoragePath,
		"f",
		fileStoragePath,
		"file storage path",
	)
	err = flag.CommandLine.Set("f", "fromflagsstorage.txt")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -f: %w", err)
	}

	restore := false
	flag.BoolVar(
		&restore,
		"r",
		restore,
		"need restore",
	)
	err = flag.CommandLine.Set("r", "true")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -r: %w", err)
	}

	dsn := ""
	flag.StringVar(
		&dsn,
		"d",
		dsn,
		"database dsn",
	)
	err = flag.CommandLine.Set("d", "postgres://testfromflag")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -d: %w", err)
	}

	hashKey := ""
	flag.StringVar(
		&hashKey,
		"k",
		hashKey,
		"hash key",
	)
	err = flag.CommandLine.Set("k", "testfromflag")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -k: %w", err)
	}

	needGenKeys := false
	flag.BoolVar(
		&needGenKeys,
		"need-gen",
		needGenKeys,
		"need to generate a private and public key for asymmetric encryption",
	)
	err = flag.CommandLine.Set("need-gen", "false")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -need-gen: %w", err)
	}

	cryptoKey := ""
	flag.StringVar(
		&cryptoKey,
		"crypto-key",
		cryptoKey,
		"path to the file with the private key",
	)
	err = flag.CommandLine.Set("crypto-key", "testfromflag")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -crypto-key: %w", err)
	}

	cryptoPublicKey := ""
	flag.StringVar(
		&cryptoPublicKey,
		"crypto-pub-key",
		cryptoPublicKey,
		"path to the file with the private key",
	)
	err = flag.CommandLine.Set("crypto-pub-key", "testfromflag")
	if err != nil {
		return nil, fmt.Errorf("failed set flag -crypto-pub-key: %w", err)
	}

	return &flags{
		adress:          adress,
		config:          config,
		logLevel:        logLevel,
		cryptoKey:       cryptoKey,
		cryptoPublicKey: cryptoPublicKey,
		fileStoragePath: fileStoragePath,
		tempDur:         tempDur,
		restore:         restore,
		dsn:             dsn,
		hashKey:         &hashKey,
		needGenKeys:     needGenKeys,
	}, nil
}

func TestConfig(t *testing.T) {
	defaultCfg := &Config{
		SConfig: ServerConfig{
			Address:              DefaultServerAddress,
			CryptoPrivateKeyPath: DefaultCryptoPrivateKeyPath,
			CryptoPublicKeyPath:  DefaultCryptoPublicKeyPath,
			ConfigsFileName:      DefaultConfigsFileName,
		},
		LConfig: LoggerConfig{
			Level: "info",
		},
		DBConfig: DataBaseConfig{
			PostgresConfig: &PostgresConfig{
				DSN: "",
			},
			FileStorageConfig: &FileStorageConfig{
				FileStoragePath: "storage.txt",
				StoreInterval:   300 * time.Second,
				Restore:         true,
			},
		},
		RetryConfig: RetryConfig{
			MinDelay:   time.Second,
			MaxDelay:   5 * time.Second,
			MaxAttempt: 3,
		},
	}

	var cfg *Config
	t.Run("New()", func(t *testing.T) {
		cfg = New()
		require.Equal(t, defaultCfg, cfg)
	})

	flgs, err := SetTestFlags(cfg)
	require.NoError(t, err)

	t.Run("setFlagsVariables()", func(t *testing.T) {
		hashKey := "testfromflag"
		wantCfg := &Config{
			SConfig: ServerConfig{
				Address:              "www.testfromflag",
				HashKey:              &hashKey,
				CryptoPrivateKeyPath: "testfromflag",
				CryptoPublicKeyPath:  "testfromflag",
				ConfigsFileName:      "config.json",
			},
			LConfig: LoggerConfig{
				Level: "levelfromflag",
			},
			DBConfig: DataBaseConfig{
				FileStorageConfig: &FileStorageConfig{
					FileStoragePath: "fromflagsstorage.txt",
					StoreInterval:   500 * time.Second,
					Restore:         true,
				},
				PostgresConfig: &PostgresConfig{
					DSN: "postgres://testfromflag",
				},
			},
			RetryConfig: RetryConfig{
				MinDelay:   time.Second,
				MaxDelay:   5 * time.Second,
				MaxAttempt: 3,
			},
		}

		err := cfg.setFlagsVariables(flgs)
		require.NoError(t, err)

		assert.EqualValues(t, wantCfg, cfg)
	})

	t.Run("setEnvServerConfig()", func(t *testing.T) {
		key := "fromenv"
		err := os.Setenv("ADDRESS", "www.fromenv.ru")
		require.NoError(t, err)

		err = os.Setenv("KEY", key)
		require.NoError(t, err)

		err = cfg.setEnvServerConfig()
		require.NoError(t, err)
		require.Equal(
			t,
			ServerConfig{
				Address:              "www.fromenv.ru",
				HashKey:              &key,
				CryptoPrivateKeyPath: "testfromflag",
				CryptoPublicKeyPath:  "testfromflag",
				ConfigsFileName:      "config.json",
			},
			cfg.SConfig,
		)
	})

	t.Run("setEnvLoggerConfig()", func(t *testing.T) {
		err := os.Setenv("LOG_LEVEL", "levelenv")
		require.NoError(t, err)

		cfg.setEnvLoggerConfig()
		require.Equal(
			t,
			LoggerConfig{
				Level: "levelenv",
			},
			cfg.LConfig,
		)
	})
	t.Run("setDBConfig()", func(t *testing.T) {
		err := os.Setenv("FILE_STORAGE_PATH", "fromenv.txt")
		require.NoError(t, err)
		err = os.Setenv("STORE_INTERVAL", "500")
		require.NoError(t, err)
		err = os.Setenv("RESTORE", "true")
		require.NoError(t, err)
		err = os.Setenv("DATABASE_DSN", "fromenv")
		require.NoError(t, err)

		err = cfg.setDBConfig()
		require.NoError(t, err)
		require.Equal(
			t,
			DataBaseConfig{
				FileStorageConfig: &FileStorageConfig{
					FileStoragePath: "fromenv.txt",
					Restore:         true,
					StoreInterval:   500000000000,
				},
				PostgresConfig: &PostgresConfig{
					DSN: "fromenv",
				},
			},
			cfg.DBConfig,
		)
	})
	t.Run("setDBConfig()_failed_parse_store_interval", func(t *testing.T) {
		err = os.Setenv("STORE_INTERVAL", "notnumber")
		require.NoError(t, err)

		err = cfg.setDBConfig()
		require.Error(t, err)
	})
	t.Run("setDBConfig()_failed_parse_restore", func(t *testing.T) {
		err = os.Setenv("RESTORE", "notbool")
		require.NoError(t, err)

		err = cfg.setDBConfig()
		require.Error(t, err)
	})
}

func TestIsFlagPassed(t *testing.T) {
	t.Run("flag -f is set in storage.txt", func(t *testing.T) {
		res := isFlagPassed("f")
		assert.Equal(t, true, res)
	})
	t.Run("flag -unknownflag is not set", func(t *testing.T) {
		res := isFlagPassed("unknownflag")
		assert.Equal(t, false, res)
	})
}

func Test_FileConfig(t *testing.T) {
	hashkey := "fromjson"
	wantCfg := &Config{
		SConfig: ServerConfig{
			Address:              "fromjson",
			HashKey:              &hashkey,
			CryptoPrivateKeyPath: "fromjson",
			CryptoPublicKeyPath:  "fromjson",
			ConfigsFileName:      "fromjson",
			NeedGenKeys:          true,
		},
		LConfig: LoggerConfig{
			Level: "fromjson",
		},
		DBConfig: DataBaseConfig{
			FileStorageConfig: &FileStorageConfig{
				FileStoragePath: "fromjson",
				StoreInterval:   500,
				Restore:         true,
			},
			PostgresConfig: &PostgresConfig{
				DSN: "fromjson",
			},
		},
		RetryConfig: RetryConfig{
			MinDelay:   time.Second,
			MaxDelay:   5 * time.Second,
			MaxAttempt: 3,
		},
	}

	err := os.Setenv("CONFIG", "test_config.json")
	require.NoError(t, err)

	cfg := New()
	err = cfg.fileBuild()
	require.NoError(t, err)

	require.Equal(t, wantCfg, cfg)
}
