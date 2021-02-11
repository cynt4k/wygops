package randutil

import (
	"math/rand"
	"time"

	"github.com/cynt4k/wygops/pkg/util/random"
)

var (
	initialized = false
	letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rnd         = rand.New(random.CryptoSource{}) // nolint:gosec
)

func init() { // nolint:gochecknoinits
	if !initialized {
		rnd.Seed(time.Now().UnixNano())
	}
}

// RandStringRunes : Generate a random string with fixed length
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rnd.Intn(len(letterRunes))]
	}
	return string(b)
}
