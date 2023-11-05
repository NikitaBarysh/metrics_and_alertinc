package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"os"
	"sync"
)

type FileEngine struct {
	storePath string
	mu        sync.Mutex
}

func NewFileEngine(storePath string) (*FileEngine, error) {
	return &FileEngine{
		storePath: storePath,
	}, nil
}

func (f *FileEngine) SetMetrics(data map[string]entity.Metric) error {
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

func (f *FileEngine) GetAllMetric() (map[string]entity.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	metricMap := make(map[string]entity.Metric)
	file, err := os.OpenFile(f.storePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return metricMap, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var memStorage entity.Metric
		err := json.Unmarshal(data, &memStorage)
		if err != nil {
			return metricMap, fmt.Errorf("write error: %w", err)
		}
		metricMap[memStorage.ID] = memStorage
	}
	return metricMap, nil
}
