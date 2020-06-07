package randutil

import (
	"math/rand"
	"time"
)

var (
	initialized = false
	letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	if !initialized {
		rand.Seed(time.Now().UnixNano())
	}
}

// RandStringRunes : Generate a random string with fixed length
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
