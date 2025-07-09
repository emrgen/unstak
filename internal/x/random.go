package x

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)[:n]
}
