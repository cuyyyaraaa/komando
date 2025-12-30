package auth

import (
	"time"

	"komando/internal/middlewares"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo          *Repo
	jwtSecret     string
	jwtExpiresMin int
}

func NewService(repo *Repo, jwtSecret string, jwtExpiresMin int) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, jwtExpiresMin: jwtExpiresMin}
}

func (s *Service) Login(nip, password string) (*LoginResponse, error) {
	u, err := s.repo.GetUserByNIP(nip)
	if err != nil {
		return nil, err
	}
	if !u.IsActive {
		return nil, bcrypt.ErrMismatchedHashAndPassword // treat as login failed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	now := time.Now()
	exp := now.Add(time.Duration(s.jwtExpiresMin) * time.Minute)

	claims := &middlewares.AuthClaims{
		UserID:     u.UserID,
		RoleCode:   u.RoleCode,
		RegionalID: u.RegionalID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	resp := &LoginResponse{
		AccessToken: signed,
		User: UserDTO{
			UserID:     u.UserID,
			NIP:        u.NIP,
			FullName:   u.FullName,
			Email:      u.Email,
			RoleCode:   u.RoleCode,
			RegionalID: u.RegionalID,
			IsActive:   u.IsActive,
		},
	}
	return resp, nil
}

func (s *Service) Me(userID int64) (*UserDTO, error) {
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	dto := &UserDTO{
		UserID:     u.UserID,
		NIP:        u.NIP,
		FullName:   u.FullName,
		Email:      u.Email,
		RoleCode:   u.RoleCode,
		RegionalID: u.RegionalID,
		IsActive:   u.IsActive,
	}
	return dto, nil
}
