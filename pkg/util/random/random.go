package random

import (
	"math/rand"
	// "crypto/rand"

	"unsafe"
)

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func AlphaNumeric(n int) string {
	b := make([]byte, n)
	cache, remain := rand.Int63(), rs6LetterIdxMax // nolint:gosec
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), rs6LetterIdxMax // nolint:gosec
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
