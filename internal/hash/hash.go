package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Compute(body []byte, key string) string {
	if key == "" {
		return ""
	}
	h := sha256.New()
	h.Write(body)
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
