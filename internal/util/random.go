package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numeric = "0123456789"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomFloat generates a random float between min and max
func RandomFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomNumericString generates a random numeric string of length n
func RandomNumericString(n int) string {
	var sb strings.Builder
	k := len(numeric)

	for i := 0; i < n; i++ {
		c := numeric[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
