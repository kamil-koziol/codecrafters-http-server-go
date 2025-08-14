package encoding

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

type Encoder interface {
	Encode([]byte) ([]byte, error)
}

type PlainEncoder struct {
}

func (e *PlainEncoder) Encode(b []byte) ([]byte, error) {
	return b, nil
}

type GZIPEncoder struct {
}

func (e *GZIPEncoder) Encode(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(b)
	if err != nil {
		return nil, fmt.Errorf("unable to gzip write: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
