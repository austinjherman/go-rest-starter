package util

import (
	"fmt"
	"math/rand"
	"time"
)

// StringRandom generates a random string with length
func StringRandom(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}