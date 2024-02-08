// Package sender - Содержит бизнес логику
package sender

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NikitaBarysh/metrics_and_alertinc/internal/encrypt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/compress"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
)

type Sender struct {
	hash *hasher.Hasher
}

func NewSender(hash *hasher.Hasher) *Sender {
	return &Sender{
		hash: hash,
	}
}

// SendPostCompressJSON - отправка сжатых данных на сервер
func (s *Sender) SendPostCompressJSON(ctx context.Context, url string, storage entity.Metric) {
	data, err := json.Marshal(storage)
	if err != nil {
		panic(err)
	}

	buf, err := compress.Compress(data)
	if err != nil {
		panic(err)
	}
	buffer := buf.Bytes()

	request, err := http.NewRequest(http.MethodPost, url, buf)
	request = request.WithContext(ctx)
	if err != nil {
		panic(err)
	}
	if encrypt.MetricsEncryptor != nil {
		encryptBuf, err := encrypt.MetricsEncryptor.Encrypt(buffer)
		if err != nil {
			fmt.Println("err to encrypt")
		}
		buffer = encryptBuf
	}
	if s.hash != nil {
		hash, errSign := s.hash.NewSign(buffer)
		if errSign != nil {
			fmt.Println(fmt.Errorf("SendMetric: NewSign: %w", err))
		}
		request.Header.Set("HashSHA256", hex.EncodeToString(hash))
	}
	request.Header.Set(`Content-Type`, "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		service.Retry(func() error {
			retryClient := &http.Client{}
			res, err := retryClient.Do(request)
			if err != nil {
				fmt.Println("can't do retry request")
				return err
			}
			errBody := res.Body.Close()
			if errBody != nil {
				fmt.Println("can't close body in retry sender")
				return errBody
			}
			return err
		}, 0)
		fmt.Println(fmt.Errorf("usecase: sender: sendPostJSON: do request: %w", err))
		return
	}
	errBody := res.Body.Close()
	if errBody != nil {
		fmt.Println(fmt.Errorf("usecase: sender: sendPostJSON: close Body: %w", err))
		return
	}
}
