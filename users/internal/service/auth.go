package service

import (
	"awesomeProject1/internal/entity"
	"awesomeProject1/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	repo      repository.UserRepository
	secretkey []byte
	TokenTTL  time.Duration
}

func generatePasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func NewAuthService(repo repository.UserRepository, secretkey []byte, ttl time.Duration) *AuthService {
	return &AuthService{repo: repo,
		secretkey: secretkey,
		TokenTTL:  ttl,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, username, password, email string) (int64, error) {
	cryptedPassword := generatePasswordHash(password)
	return s.repo.Create(ctx, &entity.User{
		Name:         username,
		Email:        email,
		PasswordHash: cryptedPassword,
	})
}

func (s *AuthService) GenerateToken(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found, unable to generate token")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("passwords do not match")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(s.TokenTTL).Unix(),
		"sub":   user.ID,
	})
	tokenString, err := token.SignedString(s.secretkey)
	if err != nil {
		return "", errors.New("unable to generate token")
	}
	return tokenString, nil
}

func (s *AuthService) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretkey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["sub"].(float64)), nil
	}
	return 0, errors.New("invalid token")
}
