package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
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

func (f *FileEngine) SetMetric(data []entity.Metric) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	file, err := os.OpenFile(f.storePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("service: file_engine: SetMetric: OpenFile: %w", err)
	}
	defer file.Close() //TODO
	for _, metricValue := range data {
		metricValueJSON, err := json.Marshal(metricValue)
		if err != nil {
			return fmt.Errorf("service: file_engine: SetMetric : Marshal err: %w", err)
		}
		_, err = file.Write(metricValueJSON)
		if err != nil {
			return fmt.Errorf("service: file_engine: SetMetric: write []byte: %w", err)
		}
		_, writeErr := file.WriteString("\n")
		if writeErr != nil {
			return fmt.Errorf("service: file_engine: SetMetric: error writeString: %w", writeErr)
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
		return nil, fmt.Errorf("service: file_engine: GetAllMetric: OpenFile: %w", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		if len(data) == 0 {
			return nil, errors.New("empty fail")
		}
		var memStorage entity.Metric
		err := json.Unmarshal(data, &memStorage)
		if err != nil {
			return nil, fmt.Errorf("service: file_engine: GetAllMetric: Unmarshal: %w", err)
		}
		metricMap[memStorage.ID] = memStorage
	}
	return metricMap, nil
}
