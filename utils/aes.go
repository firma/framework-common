package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/codec"
	"io"
)

// 加密字符串
func GcmEncrypt(key, plaintext string) (string, error) {
	if len(key) != 32 && len(key) != 24 && len(key) != 16 {
		return "", errors.New("the length of key is error")
	}

	if len(plaintext) < 1 {
		return "", errors.New("plaintext is null")
	}

	keyByte := []byte(key)
	plainByte := []byte(plaintext)

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	seal := aesGcm.Seal(nonce, nonce, plainByte, nil)
	return base64.URLEncoding.EncodeToString(seal), nil
}

// 解密字符串
func GcmDecrypt(key, cipherText string) (string, error) {
	if len(key) != 32 && len(key) != 24 && len(key) != 16 {
		return "", errors.New("the length of key is error")
	}

	if len(cipherText) < 1 {
		return "", errors.New("cipherText is null")
	}

	cipherByte, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(cipherByte) < 12 {
		return "", errors.New("cipherByte is error")
	}

	nonce, cipherByte := cipherByte[:12], cipherByte[12:]

	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainByte, err := aesGcm.Open(nil, nonce, cipherByte, nil)
	if err != nil {
		return "", err
	}

	return string(plainByte), nil
}

// 生成32位md5字串
func GetAesKey(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))

}

func DesECBEncryptBase64(data []byte, key string) []byte {

	dataString, err := DesECBEncrypt(data, []byte(key))
	if err != nil {
		return nil
	}

	return []byte(base64.StdEncoding.EncodeToString(dataString))
}

func DesECBEncrypt(data, key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, nil
	}
	//NewCipher创建一个新的加密块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()
	data = Pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}

	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//Encrypt加密第一个块，将其结果保存到dst
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesECBBase64Decrypt(data string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	mode := codec.NewECBDecrypter(block)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}
