package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func GenerateRsaKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privily, _ := rsa.GenerateKey(rand.Reader, bits)
	return privily, &privily.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

func GeneSignString(s string, pri string) (string, error) {
	msgHash := sha1.New()
	_, err := msgHash.Write([]byte(s))
	if err != nil {
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)

	privateKey, _ := ParseRsaPrivateKeyFromPemStr(pri)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, msgHashSum)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func GeneSign2String(s string, pri string) (string, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(s))
	if err != nil {
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)

	block, _ := pem.Decode([]byte(pri))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, msgHashSum)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func SignVerify(s string, signature string, pubKey string) bool {
	msgHash := crypto.SHA1
	hashIn := msgHash.New()
	_, err := hashIn.Write([]byte(s))
	if err != nil {
		return false
	}
	msgHashSum := hashIn.Sum(nil)

	sb, _ := base64.StdEncoding.DecodeString(signature)

	rsaPubKey, _ := ParseRsaPublicKeyFromPemStr(pubKey)
	if err := rsa.VerifyPKCS1v15(rsaPubKey, msgHash, msgHashSum, sb); err != nil {
		return false
	}

	return true
}

func Sign2Verify(s string, signature string, pubKey string) bool {
	msgHash := crypto.SHA256
	hashIn := msgHash.New()
	_, err := hashIn.Write([]byte(s))
	if err != nil {
		return false
	}
	msgHashSum := hashIn.Sum(nil)

	sb, _ := base64.StdEncoding.DecodeString(signature)

	rsaPubKey, _ := ParseRsaPublicKeyFromPemStr(pubKey)
	if err := rsa.VerifyPKCS1v15(rsaPubKey, msgHash, msgHashSum, sb); err != nil {
		return false
	}

	return true
}

func Sign2VerifyStrKey(s string, signature string, pubKey string) bool {
	msgHash := crypto.SHA256
	hashIn := msgHash.New()
	_, err := hashIn.Write([]byte(s))
	if err != nil {
		return false
	}
	msgHashSum := hashIn.Sum(nil)

	sb, _ := base64.StdEncoding.DecodeString(signature)

	key, _ := base64.StdEncoding.DecodeString(pubKey)
	rsaPubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return false
	}
	if err := rsa.VerifyPKCS1v15(rsaPubKey.(*rsa.PublicKey), msgHash, msgHashSum, sb); err != nil {
		return false
	}

	return true
}

// 生成RSA私钥和公钥字符串
// bits 证书大小
// @return privateKeyStr publicKeyStr error
func GenerateRSAKey(bits int) (string, string, error) {
	var privateKeyStr, publicKeyStr string

	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return privateKeyStr, publicKeyStr, err
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}

	privateBuf := new(bytes.Buffer)
	pem.Encode(privateBuf, &privateBlock)
	privateKeyStr = privateBuf.String()

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return publicKeyStr, privateKeyStr, err
	}
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}

	publicBuf := new(bytes.Buffer)
	pem.Encode(publicBuf, &publicBlock)
	publicKeyStr = publicBuf.String()

	return privateKeyStr, publicKeyStr, nil
}
