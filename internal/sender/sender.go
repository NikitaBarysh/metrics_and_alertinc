package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/compress"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/storage/repositories"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendPost(ctx context.Context, url string, storage repositories.MemStorageStruct) {
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

func (s *Sender) SendPostCompressJSON(ctx context.Context, url string, storage repositories.MemStorageStruct) {
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
	request.Header.Set(`Content-Type`, "appllication/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
}
