// https://gist.github.com/arxdsilva/8caeca47b126a290c4562a25464895e8
package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateNewToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
