package rwfile

import (
	"encoding/json"
	"fmt"

	"github.com/zhenyanesterkova/metricsmonitor/internal/storage/memstorage"
)

func (c *FileWriter) WriteSnapStorage(memento memstorage.Memento) error {
	data, err := json.Marshal(&memento)
	if err != nil {
		return fmt.Errorf("rwfile: func WriteSnapStorage() - %w", err)
	}

	_, err = c.file.WriteAt(data, 0)
	if err != nil {
		return fmt.Errorf("rwfile: func WriteSnapStorage() - %w", err)
	}

	return nil
}

func (c *FileReader) ReadSnapStorage(memento *memstorage.Memento) error {
	err := c.decoder.Decode(memento)
	if err != nil {
		return fmt.Errorf("rwfile: func ReadSnapStorage() - %w", err)
	}
	return nil
}
