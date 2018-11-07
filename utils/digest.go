package utils

import (
	"crypto/sha256"
	"fmt"
)

func NewSHA256Digest(data []byte) string {
	var sum string

	h := sha256.New()
	h.Write(data)
	sum = fmt.Sprintf("%x", h.Sum(nil))
	return sum
}
