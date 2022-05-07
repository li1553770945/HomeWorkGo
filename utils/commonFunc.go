package utils

import (
	"math/rand"
	"time"
)

func GetToken() (token string) {
	rand.Seed(time.Now().UnixNano())
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < 10; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}

	return string(result) + time.Now().Format("20060102150405")
}
