package compress

import (
	"bytes"
	"github.com/andybalholm/brotli"
)

type Brotli struct {
}

func NewBrotli() Brotli {
	return Brotli{}
}

func (b Brotli) Encode(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	bw := brotli.NewWriter(&buf)
	_, err := bw.Write(data)
	if err != nil {
		return nil, err
	}
	err = bw.Close()

	return buf.Bytes(), nil
}

func (b Brotli) Decode(data []byte) ([]byte, error) {
	br := brotli.NewReader(bytes.NewReader(data))
	var buf bytes.Buffer
	_, err := buf.ReadFrom(br)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
