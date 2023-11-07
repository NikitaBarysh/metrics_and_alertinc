package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

var Sign *Hasher

type Hasher struct {
	key []byte
}

func NewHasher(key []byte) *Hasher {
	return &Hasher{key: key}
}

func (h *Hasher) NewSign(buff []byte) ([]byte, error) {
	//fmt.Println("buff", buff)
	m := hmac.New(sha256.New, h.key)

	_, err := m.Write(buff)
	if err != nil {
		return nil, fmt.Errorf("service: hasher: NewSign: Write: %w", err)
	}
	return m.Sum(nil), nil
}

func (h *Hasher) CheckSign(data []byte, sign []byte) error {
	newSign, err := h.NewSign(data)
	if err != nil {
		return fmt.Errorf("service: hasher: CheckSign: NewSign: %w", err)
	}
	//fmt.Println(fmt.Sprintf("%x : %x", newSign, sign))
	if !hmac.Equal(newSign, sign) {
		return fmt.Errorf("sign not equal: %w", err)
	}

	return nil
}
