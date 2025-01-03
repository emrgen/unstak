package compress

type Compress interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}
