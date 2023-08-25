package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/compress"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
	"net/http"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendPost(ctx context.Context, url string) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	request = request.WithContext(ctx)
	if err != nil {
		panic(err)
	}
	request.Header.Set(`Content-Type`, "text/plain")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
}

func (s *Sender) SendPostCompress(ctx context.Context, url string, storage repositories.MemStorageStruct) {
	data, err := json.Marshal(storage)
	if err != nil {
		panic(err)
	}
	compress.Compress(data)
	request, err := http.NewRequest(http.MethodPost, url, data)
	request = request.WithContext(ctx)
	if err != nil {
		panic(err)
	}
	request.Header.Set(`Content-Type`, "appllication/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
}
