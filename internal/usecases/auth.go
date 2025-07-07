package usecases

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"jwt-service/internal/entities"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionRepo interface {
	CreateSession(session *entities.Session) error
	FindSessionByJTI(jti string) (*entities.Session, error)
	MarkSessionAsUsed(id string) error
}

type AuthUseCase struct {
	repo        SessionRepo
	jwtSecret   string
	refreshTTL  time.Duration
	emailSender EmailSender
}

type EmailSender interface {
	SendWarning(userGUID, newIP string) error
}

func NewAuthUseCase(repo SessionRepo, jwtSecret string, refreshTTL time.Duration, emailSender EmailSender) *AuthUseCase {
	return &AuthUseCase{
		repo:        repo,
		jwtSecret:   jwtSecret,
		refreshTTL:  refreshTTL,
		emailSender: emailSender,
	}
}

func (uc *AuthUseCase) GenerateTokens(userGUID, ip string) (*entities.TokenPair, error) {
	// Проверка правильности формата GUID
	if _, err := uuid.Parse(userGUID); err != nil {
		return nil, fmt.Errorf("invalid user GUID format: %w", err)
	}

	accessToken, jti, err := generateAccessToken(userGUID, ip, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	session := &entities.Session{
		ID:          generateUUID(),
		UserGUID:    userGUID,
		RefreshHash: string(refreshHash),
		JTI:         jti,
		IPAddress:   ip,
		ExpiresAt:   time.Now().Add(uc.refreshTTL),
	}

	if err := uc.repo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return &entities.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *AuthUseCase) RefreshTokens(oldAccessToken, refreshToken, clientIP string) (*entities.TokenPair, error) {
	claims, err := parseAccessToken(oldAccessToken, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("jti not found or not a string in token claims")
	}

	session, err := uc.repo.FindSessionByJTI(jti)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if session.Used {
		return nil, errors.New("refresh token already used")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(session.RefreshHash), []byte(refreshToken)); err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Если IP адреса не совпадают, отправить предупреждение
	if session.IPAddress != clientIP {
		_ = uc.emailSender.SendWarning(session.UserGUID, clientIP)
	}

	newAccess, newJTI, err := generateAccessToken(session.UserGUID, clientIP, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefresh, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	newRefreshHash, err := bcrypt.GenerateFromPassword([]byte(newRefresh), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash new refresh token: %w", err)
	}

	newSession := &entities.Session{
		ID:          generateUUID(),
		UserGUID:    session.UserGUID,
		RefreshHash: string(newRefreshHash),
		JTI:         newJTI,
		IPAddress:   clientIP,
		ExpiresAt:   time.Now().Add(uc.refreshTTL),
	}

	// Инвалидируем старую сессию и сохраняем новую
	if err := uc.repo.MarkSessionAsUsed(session.ID); err != nil {
		return nil, fmt.Errorf("failed to invalidate old session: %w", err)
	}

	if err := uc.repo.CreateSession(newSession); err != nil {
		return nil, fmt.Errorf("failed to save new session: %w", err)
	}

	return &entities.TokenPair{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func generateAccessToken(userGUID, ip, secret string) (string, string, error) {
	jti := generateUUID()
	claims := jwt.MapClaims{
		"user_guid": userGUID,
		"ip":        ip,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"jti":       jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(secret))
	return signedToken, jti, err
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func parseAccessToken(tokenString, secret string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("access token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || token == nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims")
	}

	return claims, nil
}

func generateUUID() string {
	return uuid.New().String()
}
