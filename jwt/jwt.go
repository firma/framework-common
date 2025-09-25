package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// defaultSigningKey 默认签名密钥（仅用于测试，生产环境必须通过参数传入）
	defaultSigningKey = "I6IkiJIUzI1NiGciOInR5cCIspX_TEST_KEY_ONLY"
	// ErrInvalidKey 无效密钥错误
	ErrInvalidKey = errors.New("invalid signing key")
	// ErrInvalidAlgorithm 无效算法错误
	ErrInvalidAlgorithm = errors.New("invalid signing algorithm")
)

// Create 生成token
func Create[T jwt.Claims](claims T, key string) (string, error) {
	// 验证密钥
	if key == "" {
		return "", ErrInvalidKey
	}
	// 密钥长度检查（HS256 建议至少 32 字节）
	if len(key) < 32 {
		return "", fmt.Errorf("key too short: minimum 32 characters required, got %d", len(key))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

// Parse 解析token
func Parse[T jwt.Claims](tokenString string, claims T, key string) (T, error) {
	// 验证密钥
	if key == "" {
		return claims, ErrInvalidKey
	}
	// 密钥长度检查
	if len(key) < 32 {
		return claims, fmt.Errorf("key too short: minimum 32 characters required, got %d", len(key))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 严格验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidAlgorithm
		}

		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidAlgorithm
		}
		return []byte(key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return claims, err
	}
	
	if claimsRes, ok := token.Claims.(T); ok && token.Valid {
		return claimsRes, nil
	}
	
	return claims, errors.New("invalid token or claims")
}

// CreateWithDefaultKey 使用默认密钥生成token（仅用于测试）
func CreateWithDefaultKey[T jwt.Claims](claims T) (string, error) {
	return Create(claims, defaultSigningKey)
}

// ParseWithDefaultKey 使用默认密钥解析token（仅用于测试）
func ParseWithDefaultKey[T jwt.Claims](tokenString string, claims T) (T, error) {
	return Parse(tokenString, claims, defaultSigningKey)
}