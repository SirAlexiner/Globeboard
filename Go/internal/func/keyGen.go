package _func

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

var APIRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var UUIDRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func GenerateAPIKey(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = APIRunes[rand.Intn(len(APIRunes))]
	}
	return string(b)
}

func GenerateUUID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = UUIDRunes[rand.Intn(len(UUIDRunes))]
	}
	return string(b)
}
