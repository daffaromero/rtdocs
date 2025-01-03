package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

type TokenGenerator interface {
	GenerateToken(ID, username string) (*Token, error)
	ValidateToken(token string) (jwt.Claims, error)
}

type tokenGenerator struct {
	secretKey            string
	accessTokenDuration  string
	refreshTokenDuration string
}

func NewTokenGenerator(secretKey, accessTokenDuration, refreshTokenDuration string) TokenGenerator {
	return &tokenGenerator{
		secretKey:            secretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (t *tokenGenerator) GenerateToken(ID, username string) (*Token, error) {
	durationAccess, err := time.ParseDuration(t.accessTokenDuration)
	if err != nil {
		return nil, err
	}
	durationRefresh, err := time.ParseDuration(t.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"user_id":  ID,
		"username": username,
		"exp":      time.Now().Add(durationAccess).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(t.secretKey))
	if err != nil {
		return nil, err
	}

	claims["exp"] = time.Now().Add(durationRefresh).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString([]byte(t.secretKey))
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (t *tokenGenerator) ValidateToken(token string) (jwt.Claims, error) {
	fields := strings.Split(token, " ")
	if len(fields) != 2 {
		return nil, errors.New("invalid token")
	}

	jwtToken, err := jwt.Parse(fields[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	return jwtToken.Claims, nil
}
