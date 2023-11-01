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

//func (f *FileEngine) SetMetrics(data []entity.Metric) error {
//	f.mu.Lock()
//	defer f.mu.Unlock()
//	fmt.Println("11")
//	fmt.Println(data)
//	file, err := os.OpenFile(f.storePath, os.O_CREATE|os.O_WRONLY, 0666)
//	if err != nil {
//		return fmt.Errorf("service: file_engine: SetMetric: OpenFile: %w", err)
//	}
//	metricSlice := make([]entity.Metric, 0, 35)
//	fmt.Println(file)
//	fmt.Println("22")
//	//if err != nil {
//	//	return fmt.Errorf("file_engine: SetMetric: getCounterMetric: %w", err)
//	//}
//	scanner := bufio.NewScanner(file)
//	defer file.Close()
//
//	for scanner.Scan() {
//		data := scanner.Bytes()
//		var metric entity.Metric
//		err := json.Unmarshal(data, &metric)
//		if err != nil {
//			fmt.Println(err)
//		}
//		metricSlice = append(metricSlice, metric)
//	}
//	var metricDelta int64
//	for _, v := range metricSlice {
//		if v.ID == "PollCount" {
//			metricDelta = v.Delta
//		}
//	}
//
//	fmt.Println("44")
//	for _, metricValue := range data {
//		fmt.Println("55")
//		delta := metricDelta + metricValue.Delta
//		metric := entity.Metric{
//			ID:    metricValue.ID,
//			MType: metricValue.MType,
//			Delta: delta,
//			Value: metricValue.Value}
//		fmt.Println(metricValue)
//		metricValueJSON, err := json.Marshal(metric)
//		fmt.Println(metricValueJSON)
//		if err != nil {
//			return fmt.Errorf("service: file_engine: SetMetric : Marshal err: %w", err)
//		}
//		fmt.Println("66")
//		_, err = file.Write(metricValueJSON)
//		fmt.Println("77")
//		if err != nil {
//			return fmt.Errorf("service: file_engine: SetMetric: write []byte: %w", err)
//		}
//		fmt.Println("88")
//		_, writeErr := file.WriteString("\n")
//		fmt.Println("99")
//		if writeErr != nil {
//			return fmt.Errorf("service: file_engine: SetMetric: error writeString: %w", writeErr)
//		}
//		fmt.Println("10 10")
//	}
//	fmt.Println("11 11")
//	return nil
//}

//func (f *FileEngine) GetAllMetric() ([]entity.Metric, error) {
//	f.mu.Lock()
//	defer f.mu.Unlock()
//	metricSlice := make([]entity.Metric, 0, 35)
//	file, err := os.OpenFile(f.storePath, os.O_RDONLY|os.O_CREATE, 0666)
//	if err != nil {
//		return nil, fmt.Errorf("service: file_engine: GetAllMetric: OpenFile: %w", err)
//	}
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		data := scanner.Bytes()
//		if len(data) == 0 {
//			return nil, errors.New("empty fail")
//		}
//		var mem_storage entity.Metric
//		err := json.Unmarshal(data, &mem_storage)
//		if err != nil {
//			return nil, fmt.Errorf("service: file_engine: GetAllMetric: Unmarshal: %w", err)
//		}
//		metricSlice = append(metricSlice, mem_storage)
//	}
//	return metricSlice, nil
//}

//func (f *FileEngine) GetMetric(key string) (entity.Metric, error) {
//	f.mu.Lock()
//	defer f.mu.Unlock()
//
//	fmt.Println("get metric 1")
//
//	metricSlice := make([]entity.Metric, 0, 35)
//
//	fmt.Println(metricSlice)
//	fmt.Println("get metric 2")
//
//	file, err := os.OpenFile(f.storePath, os.O_RDONLY|os.O_CREATE, 0666)
//	fmt.Println("get metric 3")
//	if err != nil {
//		return entity.Metric{}, fmt.Errorf("service: file_engine: GetMetric: OpenFile: %w", err)
//	}
//
//	fmt.Println("get metric 4")
//	scanner := bufio.NewScanner(file)
//	fmt.Println("get metric 5")
//	for scanner.Scan() {
//		fmt.Println("get metric 6")
//
//		data := scanner.Bytes()
//
//		fmt.Println(data)
//		fmt.Println("get metric 7")
//		//if len(data) == 0 {
//		//	return entity.Metric{}, models.ErrNotFound
//		//}
//		fmt.Println("get metric 8")
//		var metric entity.Metric
//		err := json.Unmarshal(data, &metric)
//		if err != nil {
//			for _, v := range metricSlice {
//				fmt.Println(v)
//				if v.ID == key {
//					fmt.Println(v.ID, key)
//					metric := entity.Metric{ID: v.ID, MType: v.MType, Delta: v.Delta, Value: v.Value}
//					return metric, nil
//
//				}
//			}
//			return entity.Metric{}, fmt.Errorf("service: file_engine: GetMetric: Unmarshal: %w", err)
//		}
//		metricSlice = append(metricSlice, metric)
//		fmt.Println(metricSlice)
//	}
//	//for _, v := range metricSlice {
//	//	fmt.Println(v)
//	//	if v.ID == key {
//	//		fmt.Println(v.ID, key)
//	//		return entity.Metric{ID: v.ID, MType: v.MType, Delta: v.Delta, Value: v.Value}, nil
//	//	}
//	//}
//	return entity.Metric{}, nil
//}
