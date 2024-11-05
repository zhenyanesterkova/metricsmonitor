package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"time"
)

type envConfig struct {
	sConfig     ServerConfig
	lConfig     LoggerConfig
	rConfig     RestoreConfig
	flagsValues flags
}

type flags struct {
	address         string
	logLevel        string
	storeInt        int
	fileStoragePath string
	restore         bool
}

func newEnvConfig() *envConfig {
	return &envConfig{}
}

func (ec *envConfig) setFlagsVariables() {
	flag.StringVar(&ec.flagsValues.address, "a", DefaultServerAddress, "address and port to run server")
	flag.StringVar(&ec.flagsValues.logLevel, "l", DefaultLogLevel, "log level")
	flag.IntVar(&ec.flagsValues.storeInt, "i", DefaultStoreInterval, "store interval")
	flag.StringVar(&ec.flagsValues.fileStoragePath, "f", DefaultFileStoragePath, "file storage path")
	flag.BoolVar(&ec.flagsValues.restore, "r", DefaultRestore, "need restore")
	flag.Parse()
}

func (ec *envConfig) SetServerConfig() {

	envEndpoint := os.Getenv("ADDRESS")

	if envEndpoint != "" {
		ec.sConfig.Address = envEndpoint
		return
	}

	ec.sConfig.Address = ec.flagsValues.address
}

func (ec *envConfig) SetLoggerConfig() {

	envLogLevel := os.Getenv("LOG_LEVEL")

	if envLogLevel != "" {
		ec.lConfig.Level = envLogLevel
		return
	}

	ec.lConfig.Level = ec.flagsValues.logLevel
}

func (ec *envConfig) SetRestoreConfig() error {

	envStoreInt := os.Getenv("STORE_INTERVAL")

	var strDur string
	if envStoreInt != "" {
		strDur = envStoreInt + "s"
	} else {
		strDur = strconv.Itoa(ec.flagsValues.storeInt) + "s"
	}
	dur, err := time.ParseDuration(strDur)
	if err != nil {
		return errors.New("can not parse store interval as duration" + err.Error())
	}
	ec.rConfig.StoreInterval = dur

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")

	if envFileStoragePath != "" {
		ec.rConfig.FileStoragePath = envFileStoragePath
	} else {
		ec.rConfig.FileStoragePath = ec.flagsValues.fileStoragePath
	}

	envRestore := os.Getenv("RESTORE")

	if envRestore != "" {
		ec.rConfig.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			return errors.New("can not parse need store" + err.Error())
		}
	} else {
		ec.rConfig.Restore = ec.flagsValues.restore
	}

	return nil
}

func (ec *envConfig) Build() (Config, error) {
	ec.setFlagsVariables()
	ec.SetServerConfig()
	ec.SetLoggerConfig()
	err := ec.SetRestoreConfig()
	if err != nil {
		return Config{}, err
	}
	return Config{
		SConfig: ServerConfig{
			Address: ec.sConfig.Address,
		},
		LConfig: LoggerConfig{
			Level: ec.lConfig.Level,
		},
		RConfig: RestoreConfig{
			StoreInterval:   ec.rConfig.StoreInterval,
			FileStoragePath: ec.rConfig.FileStoragePath,
			Restore:         ec.rConfig.Restore,
		},
	}, nil
}
