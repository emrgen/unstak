package compress

import (
	"fmt"
	lz4 "github.com/pierrec/lz4/v4"
)

type Lz4Compress struct{}

func NewLz4Compress() Lz4Compress {
	return Lz4Compress{}
}

func (l Lz4Compress) Encode(data []byte) ([]byte, error) {
	buf := make([]byte, lz4.CompressBlockBound(len(data)))

	var c lz4.Compressor
	n, err := c.CompressBlock(data, buf)
	if err != nil {
		return nil, err
	}
	if n >= len(data) {
		fmt.Printf("`%s` is not compressible", data)
	}
	buf = buf[:n] // compressed data

	return buf, nil
}

func (l Lz4Compress) Decode(data []byte) ([]byte, error) {
	out := make([]byte, 10*len(data))
	n, err := lz4.UncompressBlock(data, out)
	if err != nil {
		fmt.Println(err)
	}
	out = out[:n] // uncompressed data

	return out, nil
}
