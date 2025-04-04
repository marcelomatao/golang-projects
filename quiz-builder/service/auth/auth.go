package auth

import (
	"errors"
	"quiz-builder/config"
	"quiz-builder/entity"
	"quiz-builder/logger"

	jwt "github.com/dgrijalva/jwt-go"
)

type IAuthService interface {
	Authorize(tokenString string) (bool, error)
	CreateToken(entity.User) string
}

// AuthService service relating to authorization
type AuthService struct {
	config config.Config
	logger logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	config config.Config,
	logger logger.Logger,
) AuthService {
	return AuthService{
		config: config,
		logger: logger,
	}
}

// Authorize authorizes the generated token
func (s AuthService) Authorize(tokenString string) (bool, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims["uuid"].(string), nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, "", errors.New("token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, "", errors.New("token expired")
		}
	}
	return false, "", errors.New("couldn't handle token")
}

// CreateToken creates jwt auth token
func (s AuthService) CreateToken(user entity.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid":  user.UUID,
		"name":  user.Name,
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))

	if err != nil {
		s.logger.Error("JWT validation failed: ", err)
	}

	return tokenString
}
