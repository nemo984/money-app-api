package util

import (
	"math/rand"
	"strings"
	"time"
)

const random = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(random)

	for i := 0; i < n; i++ {
		c := random[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
