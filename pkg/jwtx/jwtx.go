// Package jwtx 提供统一的 JWT 令牌签发与解析，供网关与各 api 服务复用。
package jwtx

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

// UserIDKey 是写入/读取 JWT claim 时使用的键。
const UserIDKey = "userId"

// UserIDFromCtx 从 context 中解析 JWT claim userId，兼容 json.Number/float64/string 等类型。
func UserIDFromCtx(ctx context.Context) int64 {
	switch n := ctx.Value(UserIDKey).(type) {
	case json.Number:
		id, _ := n.Int64()
		return id
	case float64:
		return int64(n)
	case int64:
		return n
	case string:
		id, _ := strconv.ParseInt(n, 10, 64)
		return id
	default:
		return 0
	}
}

// BuildToken 用 HS256 签发令牌，claims 中带 userId。
// secret 为签名密钥，iat 为签发时间(秒)，seconds 为有效期(秒)。
func BuildToken(secret string, iat, seconds, userId int64) (string, error) {
	claims := jwt.MapClaims{
		"iat":    iat,
		"exp":    iat + seconds,
		"userId": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
