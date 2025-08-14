package encoding

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Decoder interface {
	Decode([]byte) ([]byte, error)
}

type PlainDecoder struct {
}

func (d *PlainDecoder) Decode(b []byte) ([]byte, error) {
	return b, nil
}

type GZIPDecoder struct {
}

func (d *GZIPDecoder) Decode(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	return buf, err
}
