package utils

import (
	"errors"
	"server/global"
	"server/model/request"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT  结构体
type JWT struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
}

// TokenErrorType 令牌错误类型
var (
	// TokenExpired 令牌过期错
	TokenExpired = errors.New(("token is expired"))
	// TokenInvalid 令牌无效错
	TokenInvalid = errors.New(("token is invalid"))
	//TokenMalformed 令牌格式错误
	TokenMalformed = errors.New(("thats's not even a token"))
	// TokenNotValidYet 令牌未生效错
	TokenNotValidYet = errors.New(("token is not valid yet"))
)

// NewJWT 创建 JWT 实例
func NewJWT() *JWT {
	return &JWT{
		AccessTokenSecret:  []byte(global.Config.Jwt.AccessTokenSecret),
		RefreshTokenSecret: []byte(global.Config.Jwt.RefreshTokenSecret),
	}
}

// 创建accessToken声明
func (j *JWT) CreateAccessTokenClaims(baseClaims request.BaseClaims) request.JwtCustomClaims {
	// 解析过期时间
	ep, _ := ParseDuration(global.Config.Jwt.AccessTokenExpiryTime)
	// 创建声明
	claims := request.JwtCustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Go Blog"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}
	return claims
}

// 创建accessToken
func (j *JWT) CreateAccessToken(claims request.JwtCustomClaims) (string, error) {
	//创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//返回token字符串
	return token.SignedString(j.AccessTokenSecret)
}

// 创建refreshToken声明
func (j *JWT) CreateRefreshTokenClaims(baseClaims request.BaseClaims) request.JwtCustomRefreshClaims {
	// 解析过期时间
	eq, _ := ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	// 创建声明
	claims := request.JwtCustomRefreshClaims{
		UserID: baseClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Go Blog"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(eq)),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}
	return claims
}

// 创建refreshToken
func (j *JWT) CreateRefreshToken(claims request.JwtCustomRefreshClaims) (string, error) {
	//创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//返回token字符串
	return token.SignedString(j.RefreshTokenSecret)
}

// 解析accessToken
func (j *JWT) ParseAccessToken(tokenString string) (*request.JwtCustomClaims, error) {
	claims, err := j.ParseToken(tokenString, &request.JwtCustomClaims{}, j.AccessTokenSecret)
	if err != nil {
		return nil, err
	}
	//检查类型是否转换成功
	if customClaims, ok := claims.(*request.JwtCustomClaims); ok {
		return customClaims, nil
	}
	return nil, TokenInvalid
}

// 解析refreshToken
func (j *JWT) ParseRefreshToken(tokenString string) (*request.JwtCustomRefreshClaims, error) {
	claims, err := j.ParseToken(tokenString, &request.JwtCustomRefreshClaims{}, j.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}
	//检查类型是否转换成功
	if refreshClaims, ok := claims.(*request.JwtCustomRefreshClaims); ok {
		return refreshClaims, nil
	}
	return nil, TokenInvalid
}

// 解析token
func (j *JWT) ParseToken(tokenString string, claims jwt.Claims, secretKey interface{}) (interface{}, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformed

			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpired

			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYet
			default:
				return nil, TokenInvalid
			}
		}
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, TokenInvalid
}
