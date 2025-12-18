package repositories

import (
	"time"

	"github.com/ChinawatDc/011-go-api-auth-jwt/internal/models"
	"gorm.io/gorm"
)

type TokenRepo struct{ db *gorm.DB }

func NewTokenRepo(db *gorm.DB) *TokenRepo { return &TokenRepo{db: db} }

func (r *TokenRepo) SaveRefreshToken(t *models.RefreshToken) error {
	return r.db.Create(t).Error
}

func (r *TokenRepo) FindValidRefreshToken(userID uint, tokenID string, tokenHash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.db.Where(`
		user_id = ? AND token_id = ? AND token_hash = ?
		AND revoked_at IS NULL AND expires_at > ?
	`, userID, tokenID, tokenHash, time.Now()).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepo) RevokeByID(id uint) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked_at", &now).Error
}
