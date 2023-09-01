package restorer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"os"
	"sync"
)

type FileEngine struct {
	storePath string
	mu        sync.Mutex
}

func NewFileEngine(storePath string) *FileEngine {
	return &FileEngine{
		storePath: storePath,
	}
}

func (f *FileEngine) WriteFile(data map[string]repositories.MemStorageStruct) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	file, err := os.OpenFile(f.storePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("err: %w", err)
	}
	defer file.Close()
	for _, metricValue := range data {
		metricValueJSON, err := json.Marshal(metricValue)
		if err != nil {
			return fmt.Errorf("error Marshal: %w", err)
		}
		_, err = file.Write(metricValueJSON)
		file.WriteString("\n")
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return nil
}

func (f *FileEngine) ReadFile() (map[string]repositories.MemStorageStruct, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	metricMap := make(map[string]repositories.MemStorageStruct)
	file, err := os.OpenFile(f.storePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var memStorage repositories.MemStorageStruct
		err := json.Unmarshal(data, &memStorage)
		if err != nil {
			return nil, fmt.Errorf("write error: %w", err)
		}
		metricMap[memStorage.ID] = memStorage
	}
	return metricMap, nil
}
