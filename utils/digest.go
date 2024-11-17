package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

func Base64String(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Byte(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func Base32String(bytes []byte) string {
	return base32.StdEncoding.EncodeToString(bytes)
}

func RandBase32String() string {
	key := make([]byte, 10)
	rand.Read(key)
	return base32.StdEncoding.EncodeToString(key)
}

func Hmac512String(bytes []byte, key string) (string, error) {
	mac := hmac.New(sha512.New, []byte(key))
	_, err := mac.Write(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
