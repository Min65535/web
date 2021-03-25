package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLNMOPQRSTUVWXYZ-_"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
