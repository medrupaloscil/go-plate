package services

import (
	"math/rand"
	"reflect"
)

func GetStructType(myvar interface{}) string {
	return reflect.TypeOf(myvar).String()
}

// TODO: Potential weakness, to verify
func RandStringBytes(n int, special bool) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	if special {
		letterBytes = letterBytes + "?/!.&"
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ArrayContainsUint(s []uint, e uint) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
