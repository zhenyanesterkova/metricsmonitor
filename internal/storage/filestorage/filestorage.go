package filestorage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/config"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"
	"github.com/zhenyanesterkova/metricsmonitor/internal/app/server/rwfile"
	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

type FileStorage struct {
	*memstorage.MemStorage
	w   *rwfile.FileWriter
	r   *rwfile.FileReader
	log logger.LogrusLogger
}

func New(conf config.RestoreConfig, storeLog logger.LogrusLogger) (*FileStorage, error) {
	fileWriter, err := rwfile.NewFileWriter(conf.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("file writer error: %w", err)
	}
	fileReader, err := rwfile.NewFileReader(conf.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("file reader error: %w", err)
	}

	fileStore := FileStorage{
		MemStorage: memstorage.New(),
		w:          fileWriter,
		r:          fileReader,
		log:        storeLog,
	}

	if conf.Restore {
		err = fileStore.readStorage()
		if err != nil {
			return nil, fmt.Errorf("read storage error: %w", err)
		}
	}

	go func() error {
		err := fileStore.startSaveStorage(conf.StoreInterval)
		if err != nil {
			fileStore.log.LogrusLog.Errorf("error save storage: %v", err)
			return fmt.Errorf("error save storage: %w", err)
		}
		return nil
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
	if !fs.r.Reader.Scan() {
		if fs.r.Reader.Err() != nil {
			return fmt.Errorf("error read storage from file: %w", fs.r.Reader.Err())
		}
		return nil
	}

	data := fs.r.Reader.Bytes()

	memento := fs.MemStorage.CreateMemento()
	err := json.Unmarshal(data, memento)
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

	_, err = fs.w.File.WriteAt(data, 0)
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
