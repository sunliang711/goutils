package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenJwtToken 使用key作为签名秘钥来生成jwt token，并且在token里包含了自定义数据data
// 当配置文件中包含jwt.exp字段，并且该字段大于零时，在jwt token里会加入exp来表示token过期时间
func GenJwtToken(secret string, duration int, data map[string]interface{}) (string, error) {
	mapClaims := jwt.MapClaims{}

	if duration > 0 {
		mapClaims["exp"] = time.Now().Add(time.Second * time.Duration(duration))
	}

	for k, v := range data {
		mapClaims[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenStr, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseJwtToken 解析jwt token，返回*jwt.Token对象
func ParseJwtToken(token string, secret string) (*jwt.Token, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return nil, fmt.Errorf("token expired")
			default:
				return nil, fmt.Errorf("parse token error: %v", vErr.Error())
			}
		default:
			return nil, fmt.Errorf("invalid token")
		}
	}
	return t, nil
	// token2User[t.Raw] = t.Claims.(jwt.MapClaims)["phone"].(string)
}

// GetValueFromJwtToken 从ctx中解析jwt token(token存在于请求头的"jwt.header_name"中)
