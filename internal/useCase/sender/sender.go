package sender

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/entity"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/interface/compress"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service"
	"github.com/NikitaBarysh/metrics_and_alertinc/internal/service/hasher"
	"net/http"
)

type Sender struct {
	hash *hasher.Hasher
}

func NewSender(hash *hasher.Hasher) *Sender {
	return &Sender{
		hash: hash,
	}
}

//func (s *Sender) SendPost(ctx context.Context, url string, storage entity.Metric) {
//	request, err := http.NewRequest(http.MethodPost, url, nil)
//	request = request.WithContext(ctx)
//	if err != nil {
//		panic(err)
//	}
//	request.Header.Set(`Content-Type`, "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(request)
//	if err != nil {
//		service.Retry(func() error {
//			retryClient := &http.Client{}
//			res, err := retryClient.Do(request)
//			if err != nil {
//				fmt.Println("can't do retry request")
//				return err
//			}
//			errBody := res.Body.Close()
//			if errBody != nil {
//				fmt.Println("can't close body in retry sender")
//				return errBody
//			}
//			return err
//		}, 0)
//		fmt.Println(fmt.Errorf("useCase: sender: sendPost: do request: %w", err))
//		return
//	}
//	err = res.Body.Close()
//	if err != nil {
//		fmt.Println("body not closed", err)
//	}
//}

func (s *Sender) SendPostCompressJSON(ctx context.Context, url string, storage entity.Metric) {
	data, err := json.Marshal(storage)
	if err != nil {
		panic(err)
	}
	//fmt.Println("data", data)
	buf, err := compress.Compress(data)
	if err != nil {
		panic(err)
	}
	//fmt.Println("buf sender", hex.EncodeToString(buf.Bytes()))
	//hash, errSign := s.hash.NewSign(buf.Bytes())
	//fmt.Println("hash sender", hex.EncodeToString(hash))

	request, err := http.NewRequest(http.MethodPost, url, buf)
	request = request.WithContext(ctx)
	if err != nil {
		panic(err)
	}
	if s.hash != nil {
		hash, errSign := s.hash.NewSign(buf.Bytes())
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
		fmt.Println(fmt.Errorf("useCase: sender: sendPostJSON: do request: %w", err))
		return
	}
	errBody := res.Body.Close()
	if errBody != nil {
		fmt.Println(fmt.Errorf("useCase: sender: sendPostJSON: close Body: %w", err))
		return
	}
}
