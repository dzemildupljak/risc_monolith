package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomString generate a string of random characters of given length
func GenerateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx := rand.Int63() % int64(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}
	return sb.String()
}

func ValidateExpirationTime(expTime time.Time) bool {
	currTime := time.Now().Local()
	fmt.Println("expTime", expTime)
	fmt.Println("currTime", currTime)
	fmt.Println("expTime.Sub(currTime).Seconds()", expTime.Sub(currTime).Seconds())
	return expTime.Sub(currTime).Seconds() > 0
}
