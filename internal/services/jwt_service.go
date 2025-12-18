package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/config"
)

type JWTService struct {
	cfg config.Config
}

func NewJWTService(cfg config.Config) *JWTService {
	return &JWTService{cfg: cfg}
}

func (s *JWTService) NewAccessToken(userID uint) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": s.cfg.JWTIssuer,
		"sub": userID,
		"typ": "access",
		"iat": now.Unix(),
		"exp": now.Add(time.Minute * time.Duration(s.cfg.AccessMinutes)).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.AccessSecret))
}

type RefreshTokenResult struct {
	Token     string
	TokenID   string
	TokenHash string
	ExpiresAt time.Time
}

func (s *JWTService) NewRefreshToken(userID uint) (RefreshTokenResult, error) {
	now := time.Now()
	jti := uuid.NewString()
	expires := now.Add(time.Hour * 24 * time.Duration(s.cfg.RefreshDays))

	claims := jwt.MapClaims{
		"iss": s.cfg.JWTIssuer,
		"sub": userID,
		"jti": jti,
		"typ": "refresh",
		"iat": now.Unix(),
		"exp": expires.Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := t.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return RefreshTokenResult{}, err
	}

	return RefreshTokenResult{
		Token:     raw,
		TokenID:   jti,
		TokenHash: sha256Hex(raw),
		ExpiresAt: expires,
	}, nil
}

func (s *JWTService) ParseAccess(token string) (jwt.MapClaims, error) {
	return parse(token, s.cfg.AccessSecret, "access")
}

func (s *JWTService) ParseRefresh(token string) (jwt.MapClaims, error) {
	return parse(token, s.cfg.RefreshSecret, "refresh")
}

func parse(token, secret, tokenType string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if claims["typ"] != tokenType {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
