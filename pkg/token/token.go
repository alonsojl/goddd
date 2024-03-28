package token

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const prefix = "Bearer "

// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(kid string, key []byte, method jwt.SigningMethod) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwt.SigningMethod
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	// claims := jwt.MapClaims{
	//  "foo": "bar",
	// 	"aud": data.Client.GetID(),
	// 	"exp": data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
	// }

	claims := JWTAccessClaims{
		"bar",
		jwt.RegisteredClaims{
			Audience:  []string{data.Client.GetID()},
			Subject:   data.UserID,
			ExpiresAt: jwt.NewNumericDate(data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn())),
			// IssuedAt:  jwt.NewNumericDate(data.CreateAt),
			// NotBefore: jwt.NewNumericDate(data.CreateAt),
		},
	}

	token := jwt.NewWithClaims(a.SignedMethod, claims)
	if a.SignedKeyID != "" {
		token.Header["kid"] = a.SignedKeyID
	}

	access, err := token.SignedString(a.SignedKey)
	if err != nil {
		return "", "", err
	}

	var refresh string
	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}

func (a *JWTAccessGenerate) Validate(accessToken string) (string, error) {
	var tokenStr string

	if accessToken == "" || !strings.HasPrefix(accessToken, prefix) {
		return tokenStr, errors.ErrInvalidAccessToken
	}

	tokenStr = accessToken[len(prefix):]
	token, err := jwt.ParseWithClaims(tokenStr, &JWTAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.SignedKey, nil
	})

	if err != nil {
		return tokenStr, err
	}

	if _, ok := token.Claims.(*JWTAccessClaims); !ok && !token.Valid {
		return tokenStr, errors.ErrInvalidAccessToken
	}

	return tokenStr, nil
}
