package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/compress"
	"net/http"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendPost(ctx context.Context, url string, storage entity.Metric) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	request = request.WithContext(ctx)
	if err != nil {
		panic(err)
	}
	request.Header.Set(`Content-Type`, "text/plain")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(fmt.Errorf("useCase: sender: sendPost: do request: %w", err))
		return
	}
	res.Body.Close()
}

func (s *Sender) SendPostCompressJSON(ctx context.Context, url string, storage entity.Metric) {
	data, err := json.Marshal(storage)
	if err != nil {
		panic(err)
	}
	buf, err := compress.Compress(data)
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest(http.MethodPost, url, buf)
	request = request.WithContext(ctx)
	if err != nil {
		panic(err)
	}
	request.Header.Set(`Content-Type`, "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(fmt.Errorf("useCase: sender: sendPostJSON: do request: %w", err))
		return
	}
	res.Body.Close()
}
