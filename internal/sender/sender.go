package sender

import (
	"context"
	"fmt"
	"net/http"
)

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendPost(ctx context.Context, url string) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	request.WithContext(ctx)
	request.Header.Set(`Content-Type`, "text/plain")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
}
