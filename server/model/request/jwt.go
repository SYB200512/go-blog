package request

import (
	"server/model/appTypes"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

// access_token 的声明Claims
type JwtCustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

// refresh_token 的声明Claims
type JwtCustomRefreshClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

// BaseClaims 基础Claims
type BaseClaims struct {
	UserID uint
	UUID   uuid.UUID
	RoleID appTypes.RoleID
}
