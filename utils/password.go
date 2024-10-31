package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

const (
	saltPassword    = "qkhPAGA13HocW3GAEWwb"
	defaultPassword = "3767141919"
)

func CreatePassword(str string) (password string) {
	// md5
	m := md5.New()
	m.Write([]byte(str))
	mByte := m.Sum(nil)

	// hmac
	h := hmac.New(sha256.New, []byte(saltPassword))
	h.Write(mByte)
	password = hex.EncodeToString(h.Sum(nil))

	return password
}

func CreatePasswordReturnSalt(str string) (password, salt string) {
	m := md5.New()
	m.Write([]byte(str))
	mByte := m.Sum(nil)

	sourceData := rand.NewSource(time.Now().UnixNano())
	r := rand.New(sourceData)
	salt = fmt.Sprintf("%06v", r.Int31n(1000000))
	// hmac
	h := hmac.New(sha256.New, []byte(salt))
	h.Write(mByte)
	password = hex.EncodeToString(h.Sum(nil))

	return password, salt
}

func CheckPasswordSalt(password, salt, checkPassword string) bool {
	m := md5.New()
	m.Write([]byte(password))
	mByte := m.Sum(nil)
	// hmac
	h := hmac.New(sha256.New, []byte(salt))
	h.Write(mByte)
	password = hex.EncodeToString(h.Sum(nil))

	if password == checkPassword {
		return true
	}
	return false
}

func ResetPassword() (password string) {
	m := md5.New()
	m.Write([]byte(defaultPassword))
	mStr := hex.EncodeToString(m.Sum(nil))

	password = CreatePassword(mStr)

	return
}

func GenerateLoginToken(id int64) (token string) {
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("%d%s", id, saltPassword)))
	token = hex.EncodeToString(m.Sum(nil))

	return
}

func CheckHMACSHA512(digest, secret, signature string) bool {
	secretKey := []byte(secret)
	data := []byte(digest)

	mac := hmac.New(sha512.New, secretKey)
	mac.Write(data)
	expectedMAC := mac.Sum(nil)
	//expHex := hex.EncodeToString(expectedMAC)
	//if expHex == signature {
	//	return true
	//} else {
	//	return false
	//}
	if sign, err := hex.DecodeString(signature); err != nil {
		return false
	} else {
		return hmac.Equal(expectedMAC, sign)
	}
}
