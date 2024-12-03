package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/retry"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/rwfile"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

type FileStorage struct {
	*memstorage.MemStorage
	w       *rwfile.FileWriter
	r       *rwfile.FileReader
	log     logger.LogrusLogger
	retrier retry.Retrier
}

func New(conf config.DataBaseConfig, storeLog logger.LogrusLogger, retCfg config.RetryConfig) (*FileStorage, error) {
	b := retry.NewBackoff(retCfg.Min, retCfg.Max, retCfg.MaxAttempt, nil)
	fileStore := FileStorage{
		MemStorage: memstorage.New(),
		log:        storeLog,
		retrier:    retry.NewRetrier(b, nil),
	}

	var fileWriter *rwfile.FileWriter
	err := fileStore.retrier.Run(context.TODO(), func() error {
		var err error
		fileWriter, err = rwfile.NewFileWriter(conf.FileStoragePath)
		if err != nil {
			return fmt.Errorf("failed create writer: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("file writer error: %w", err)
	}
	fileStore.w = fileWriter

	var fileReader *rwfile.FileReader
	err = fileStore.retrier.Run(context.TODO(), func() error {
		var err error
		fileReader, err = rwfile.NewFileReader(conf.FileStoragePath)
		if err != nil {
			return fmt.Errorf("failed create reader: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("file reader error: %w", err)
	}
	fileStore.r = fileReader

	if conf.Restore {
		err = fileStore.readStorage()
		if err != nil {
			return nil, fmt.Errorf("read storage error: %w", err)
		}
	}

	go func() {
		err := fileStore.startSaveStorage(conf.StoreInterval)
		if err != nil {
			fileStore.log.LogrusLog.Errorf("error save storage: %v", err)
		}
	}()

	return &fileStore, nil
}

func (fs *FileStorage) Close() error {
	err := fs.writeStorage()
	if err != nil {
		fs.log.LogrusLog.Errorf("can not write storage to file: %v", err)
		return fmt.Errorf("can not write storage to file: %w", err)
	}

	err = fs.w.Close()
	if err != nil {
		return fmt.Errorf("can not close file writer: %w", err)
	}

	err = fs.r.Close()
	if err != nil {
		return fmt.Errorf("can not close file reader: %w", err)
	}

	return nil
}

func (fs *FileStorage) readStorage() error {
	err := fs.retrier.Run(context.TODO(), func() error {
		if !fs.r.Reader.Scan() {
			if fs.r.Reader.Err() != nil {
				return fmt.Errorf("failed read from file: %w", fs.r.Reader.Err())
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error read storage from file: %w", fs.r.Reader.Err())
	}

	data := fs.r.Reader.Bytes()

	memento := fs.MemStorage.CreateMemento()
	err = json.Unmarshal(data, memento)
	if err != nil {
		return fmt.Errorf("rwfile: func ReadSnapStorage() - %w", err)
	}

	fs.MemStorage.RestoreMemento(memento)
	return nil
}

func (fs *FileStorage) writeStorage() error {
	fs.log.LogrusLog.Info("write storage to file...")
	data, err := json.Marshal(fs.MemStorage.CreateMemento())
	if err != nil {
		return fmt.Errorf("marshal storage error: %w", err)
	}

	err = fs.retrier.Run(context.TODO(), func() error {
		_, err := fs.w.File.WriteAt(data, 0)
		if err != nil {
			return fmt.Errorf("failed write to file: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("write storage error: %w", err)
	}

	return nil
}

func (fs *FileStorage) startSaveStorage(storeInterval time.Duration) error {
	ticker := time.NewTicker(storeInterval)
	for range ticker.C {
		fs.log.LogrusLog.Info("starting storage copying...")
		err := fs.writeStorage()
		if err != nil {
			fs.log.LogrusLog.Errorf("can not write storage to file: %v", err)
			return fmt.Errorf("can not write storage to file: %w", err)
		}
		fs.log.LogrusLog.Info("end storage copying...")
	}
	return nil
}

func (fs *FileStorage) Ping() (bool, error) {
	if fs.MemStorage == nil {
		return false, errors.New("filestorage is not initialized")
	}
	pingMemStorage, err := fs.MemStorage.Ping()
	if err != nil || !pingMemStorage {
		return false, fmt.Errorf("memstorage in filestorage is not initialized: %w", err)
	}
	if fs.r == nil || fs.w == nil {
		return false, fmt.Errorf("reader or writer in filestorage is not initialized: %w", err)
	}
	return true, nil
}
