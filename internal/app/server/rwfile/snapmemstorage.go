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
	if !c.reader.Scan() {
		if c.reader.Err() != nil {
			return fmt.Errorf("error scan storage memento: %w", c.reader.Err())
		}
		return nil
	}

	data := c.reader.Bytes()

	err := json.Unmarshal(data, &memento)
	if err != nil {
		return fmt.Errorf("rwfile: func ReadSnapStorage() - %w", err)
	}

	return nil
}
