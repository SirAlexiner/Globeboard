package _func

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

var Runes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func GenerateAPIKey(n int) string {
	a := make([]rune, n)
	b := make([]rune, n)
	for i := range a {
		a[i] = Runes[rand.Intn(len(Runes))]
		b[i] = Runes[rand.Intn(len(Runes))]
	}
	return "sk-" + string(a) + "tLtRjPv" + string(b)
}

func GenerateUID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = Runes[rand.Intn(len(Runes))]
	}
	return string(b)
}
