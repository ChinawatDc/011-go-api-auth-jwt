package services

import (
	"errors"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/models"
	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/repositories"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
)

type AuthService struct {
	users  *repositories.UserRepo
	tokens *repositories.TokenRepo
	jwt    *JWTService
}

func NewAuthService(users *repositories.UserRepo, tokens *repositories.TokenRepo, jwt *JWTService) *AuthService {
	return &AuthService{users: users, tokens: tokens, jwt: jwt}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	_, err := s.users.FindByEmail(email)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &models.User{Email: email, PasswordHash: hash}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(email, password string) (access string, refresh string, err error) {
	u, err := s.users.FindByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	if err := ComparePassword(u.PasswordHash, password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	access, err = s.jwt.NewAccessToken(u.ID)
	if err != nil {
		return "", "", err
	}

	rt, err := s.jwt.NewRefreshToken(u.ID)
	if err != nil {
		return "", "", err
	}

	if err := s.tokens.SaveRefreshToken(&models.RefreshToken{
		UserID:    u.ID,
		TokenID:   rt.TokenID,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
	}); err != nil {
		return "", "", err
	}

	return access, rt.Token, nil
}

func (s *AuthService) Refresh(refreshToken string) (newAccess string, err error) {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return "", ErrInvalidRefresh
	}

	userIDf, ok := claims["sub"].(float64) // jwt MapClaims decode as float64
	if !ok {
		return "", ErrInvalidRefresh
	}
	userID := uint(userIDf)

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return "", ErrInvalidRefresh
	}

	// ตรวจ tokenHash ว่าตรงกับที่เก็บใน DB และยังไม่ revoked/expired
	tokenHash := sha256Hex(refreshToken)
	t, err := s.tokens.FindValidRefreshToken(userID, jti, tokenHash)
	if err != nil {
		return "", ErrInvalidRefresh
	}

	// ออก access ใหม่
	newAccess, err = s.jwt.NewAccessToken(userID)
	if err != nil {
		return "", err
	}

	_ = t // ใช้ได้ถ้าต้องการ rotate token ต่อ (advanced)
	return newAccess, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	claims, err := s.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return ErrInvalidRefresh
	}

	userIDf, ok := claims["sub"].(float64)
	if !ok {
		return ErrInvalidRefresh
	}
	userID := uint(userIDf)

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return ErrInvalidRefresh
	}

	tokenHash := sha256Hex(refreshToken)
	t, err := s.tokens.FindValidRefreshToken(userID, jti, tokenHash)
	if err != nil {
		return ErrInvalidRefresh
	}

	return s.tokens.RevokeByID(t.ID)
}
