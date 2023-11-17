package hasher

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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
	s.rw.Header().Set("HashSHA256", hex.EncodeToString(sign))
	return s.rw.Write(b)
}

func (s *signRW) WriteHeader(status int) {
	s.rw.WriteHeader(status)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hash := r.Header.Get("HashSHA256")
		if hash == "" {
			w := newSignRW(rw)
			h.ServeHTTP(w, r)
			return
		}
		buff, _ := io.ReadAll(r.Body)
		sign, err := hex.DecodeString(hash)
		if err != nil {
			log.Println(fmt.Errorf("bad req sign: %w", err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := Sign.CheckSign(buff, sign); err != nil {
			log.Println(fmt.Errorf("CHeckSign: %w", err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		newBody := io.NopCloser(bytes.NewBuffer(buff))
		r.Body = newBody

		w := newSignRW(rw)
		h.ServeHTTP(w, r)
	})
}
