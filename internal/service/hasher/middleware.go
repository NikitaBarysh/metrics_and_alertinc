package hasher

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type signRW struct {
	rw   http.ResponseWriter
	Hash Hasher
}

func newSignRW(rw http.ResponseWriter) *signRW {
	return &signRW{
		rw: rw,
	}
}

func (s *signRW) Header() http.Header {
	return s.rw.Header()
}

func (s *signRW) Write(b []byte) (int, error) {
	sign, err := s.Hash.NewSign(b)
	if err != nil {
		return 0, fmt.Errorf("hasher: Write: NewSign: %w", err)
	}
	//fmt.Println("Write ", hex.EncodeToString(sign))
	s.rw.Header().Set("HashSHA256", hex.EncodeToString(sign))
	return s.rw.Write(b)
}

func (s *signRW) WriteHeader(status int) {
	s.rw.WriteHeader(status)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if hash := r.Header.Get("HashSHA256"); hash != "" {
			buff, _ := io.ReadAll(r.Body)
			fmt.Println(hash)
			sign, err := hex.DecodeString(hash)
			if err != nil {
				fmt.Println(fmt.Errorf("bad req sign: %w", err))
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			//fmt.Println(fmt.Sprintf("%x : %x", , sign))
			if err := Sign.CheckSign(buff, sign); err != nil {
				fmt.Println(fmt.Errorf("CHeckSign: %w", err))
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			newBody := io.NopCloser(bytes.NewBuffer(buff))
			r.Body = newBody
		}

		w := newSignRW(rw)

		h.ServeHTTP(w, r)
	})
}
