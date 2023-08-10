package auth

import (
	"os"
	"time"

	"github.com/egiferdians/micro-auth/models"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Authenticated struct
type Authenticated struct {
	User         *models.User `json:"user"`
	RefreshToken string      `json:"refresh_token"`
	AccessToken  string      `json:"access_token"`
}

// VerifyPassword compare hashed password with password string
func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken generate JWT token
func GenerateToken(id uuid.UUID) (string, string, error) {
	// AccessToken
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", "", err
	}
	// RefreshToken
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = id.String()
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refresh, err := refreshToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", "", err
	}
	return token, refresh, nil
}
