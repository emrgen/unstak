package compress

import (
	"bytes"
	"compress/gzip"
)

// GZip is a GZip implementation of the Compress interface.
type GZip struct {
}

// NewGZip returns a new GZip.
func NewGZip() GZip {
	return GZip{}
}

// Encode compresses data using GZip.
func (g GZip) Encode(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode decompresses data using GZip.
func (g GZip) Decode(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(gr)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
