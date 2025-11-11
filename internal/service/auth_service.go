package service

import (
	"context"
	"errors"
	"time"

	"github.com/LouisFernando1204/kai-backend.git/domain"
	"github.com/LouisFernando1204/kai-backend.git/dto"
	"github.com/LouisFernando1204/kai-backend.git/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository domain.UserRepository
	config         *config.Config
}

func NewAuthService(userRepo domain.UserRepository, conf *config.Config) domain.AuthService {
	return &authService{
		userRepository: userRepo,
		config:         conf,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	existingUser, err := s.userRepository.FindByEmail(ctx, req.Email)

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &domain.User{
		ID:             primitive.NewObjectID(),
		FullName:       req.FullName,
		Email:          req.Email,
		Password:       string(hashedPassword),
		PhoneNumber:    req.PhoneNumber,
		IdentityNumber: req.IdentityNumber,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.userRepository.Save(ctx, newUser)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.config.Jwt.Exp)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Jwt.Key))
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{Token: tokenString}, nil
}
