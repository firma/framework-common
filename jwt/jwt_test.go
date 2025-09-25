package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

type ClaimsItem struct {
	UserId    int64   `json:"user_id"`
	TenantId  int64   `json:"tenant_id"`
	Platform  string  `json:"platform"`
	RoleIds   []int64 `json:"role_ids"`
	IsSupper  bool    `json:"is_supper"`
	IsRefresh bool    `json:"is_refresh"`
	jwt.RegisteredClaims
}

// 测试正常的 token 创建和解析
func TestCreateAndParse(t *testing.T) {
	// 使用足够长的密钥（至少32字符）
	testKey := "jvc29RZEj1oYK1KVHjntB8fZEj1oYK1KVHjZEj1oYK1KVHjZEj1oYK1KVHjOj2dhpuXZyCOKwU"

	claims := ClaimsItem{
		UserId:    1,
		TenantId:  1,
		Platform:  "windows",
		RoleIds:   []int64{1, 2, 3},
		IsSupper:  true,
		IsRefresh: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建 token
	token, err := Create(claims, testKey)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// 解析 token
	data := &ClaimsItem{}
	parsed, err := Parse(token, data, testKey)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// 验证解析结果
	if parsed.UserId != claims.UserId {
		t.Errorf("UserId mismatch: expected %d, got %d", claims.UserId, parsed.UserId)
	}
	if parsed.TenantId != claims.TenantId {
		t.Errorf("TenantId mismatch: expected %d, got %d", claims.TenantId, parsed.TenantId)
	}
	if parsed.Platform != claims.Platform {
		t.Errorf("Platform mismatch: expected %s, got %s", claims.Platform, parsed.Platform)
	}
}

// 测试空密钥错误
func TestCreateWithEmptyKey(t *testing.T) {
	claims := ClaimsItem{
		UserId:           1,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	_, err := Create(claims, "")
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey, got: %v", err)
	}
}

// 测试密钥长度不足错误
func TestCreateWithShortKey(t *testing.T) {
	claims := ClaimsItem{
		UserId:           1,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	shortKey := "short"
	_, err := Create(claims, shortKey)
	if err == nil {
		t.Error("Expected error for short key, got nil")
	}
}

// 测试使用不同密钥解析 token
func TestParseWithWrongKey(t *testing.T) {
	correctKey := "jvc29Rn_tB8fZEj1oYK1KVHjOj2dhpuXZyCOKwU"
	wrongKey := "wrong_key_with_enough_length_1234567890"

	claims := ClaimsItem{
		UserId:           1,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	token, err := Create(claims, correctKey)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	data := &ClaimsItem{}
	_, err = Parse(token, data, wrongKey)
	if err == nil {
		t.Error("Expected error when parsing with wrong key, got nil")
	}
}

// 测试过期的 token
func TestExpiredToken(t *testing.T) {
	testKey := "jvc29Rn_tB8fZEj1oYK1KVHjOj2dhpuXZyCOKwU"

	claims := ClaimsItem{
		UserId: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 已过期
		},
	}

	token, err := Create(claims, testKey)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	data := &ClaimsItem{}
	_, err = Parse(token, data, testKey)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

// 测试默认密钥函数（仅用于测试环境）
func TestDefaultKeyFunctions(t *testing.T) {
	claims := ClaimsItem{
		UserId:           1,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// 使用默认密钥创建
	token, err := CreateWithDefaultKey(claims)
	if err != nil {
		t.Fatalf("Failed to create token with default key: %v", err)
	}

	// 使用默认密钥解析
	data := &ClaimsItem{}
	parsed, err := ParseWithDefaultKey(token, data)
	if err != nil {
		t.Fatalf("Failed to parse token with default key: %v", err)
	}

	if parsed.UserId != claims.UserId {
		t.Errorf("UserId mismatch: expected %d, got %d", claims.UserId, parsed.UserId)
	}
}
