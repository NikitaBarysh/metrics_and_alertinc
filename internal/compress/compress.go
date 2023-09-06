package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Compress(data []byte) (*bytes.Buffer, error) {
	var b bytes.Buffer

	rw := gzip.NewWriter(&b)

	_, err := rw.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %w", err)
	}

	err = rw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed close compress writer: %w", err)
	}

	return &b, nil
}
