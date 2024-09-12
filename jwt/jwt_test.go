package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
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

func TestCreate(t *testing.T) {
	claims := ClaimsItem{
		UserId:           1,
		TenantId:         1,
		Platform:         "windows",
		RoleIds:          []int64{1, 2, 3},
		IsSupper:         true,
		IsRefresh:        false,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token, err := Create(claims)
	if err != nil {
		t.Error(err)
	}
	data := &ClaimsItem{}
	parse, err := Parse(token, data)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(parse)
}
