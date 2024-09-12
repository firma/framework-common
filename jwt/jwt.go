package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	// signingKey
	signingKey = "I6IkiJIUzI1NiGciOInR5cCIspX"
)

// Create 生成token
// Create 生成token
func Create[T jwt.Claims](claims T) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

// Parse 解析token
func Parse[T jwt.Claims](tokenString string, claims T) (T, error) {

	token, err := jwt.ParseWithClaims(
		tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		}, jwt.WithLeeway(1*time.Second),
	)
	if err != nil {
		fmt.Println(err)
		return claims, err
	} else if claims, ok := token.Claims.(T); ok {
		fmt.Println(claims)
		return claims, nil
	} else {
		fmt.Println("unknown claims type, cannot proceed")
		return claims, err
	}
}
