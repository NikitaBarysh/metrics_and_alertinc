package file_storage

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

func NewFileEngine(storePath string) (*FileEngine, error) {
	return &FileEngine{
		storePath: storePath,
	}, nil
}

func (f *FileEngine) SetMetrics(data []entity.Metric) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	file, err := os.OpenFile(f.storePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("service: file_engine: SetMetric: OpenFile: %w", err)
	}
	defer file.Close()
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

func (f *FileEngine) GetAllMetric() ([]entity.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	metricSlice := make([]entity.Metric, 0, 35)
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
		metricSlice = append(metricSlice, memStorage)
	}
	return metricSlice, nil
}

func (f *FileEngine) GetMetric(key string) (entity.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	metricSlice := make([]entity.Metric, 0, 35)

	file, err := os.OpenFile(f.storePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return entity.Metric{}, fmt.Errorf("service: file_engine: GetMetric: OpenFile: %w", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		if len(data) == 0 {
			return entity.Metric{}, err
		}
		var metric entity.Metric
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return entity.Metric{}, fmt.Errorf("service: file_engine: GetMetric: Unmarshal: %w", err)
		}
		metricSlice = append(metricSlice, metric)
	}
	for _, v := range metricSlice {
		if v.ID == key {
			return entity.Metric{ID: v.ID, MType: v.MType, Delta: v.Delta, Value: v.Value}, nil
		}
	}
	return entity.Metric{}, nil
}
