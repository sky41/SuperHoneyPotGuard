package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"superhoneypotguard/config"
	"superhoneypotguard/models"
)

func GenerateToken(claims *models.Claims) (string, error) {
	cfg := config.AppConfig
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":      claims.UserID,
		"username":    claims.Username,
		"roles":       claims.Roles,
		"permissions": claims.Permissions,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenString string) (*models.Claims, error) {
	cfg := config.AppConfig
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		roles := make([]string, 0)
		if rolesInterface, ok := claims["roles"].([]interface{}); ok {
			for _, r := range rolesInterface {
				if roleStr, ok := r.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}

		permissions := make([]string, 0)
		if permissionsInterface, ok := claims["permissions"].([]interface{}); ok {
			for _, p := range permissionsInterface {
				if permStr, ok := p.(string); ok {
					permissions = append(permissions, permStr)
				}
			}
		}

		userID := int(claims["userId"].(float64))
		username := claims["username"].(string)

		return &models.Claims{
			UserID:      userID,
			Username:    username,
			Roles:       roles,
			Permissions: permissions,
		}, nil
	}

	return nil, errors.New("invalid token")
}
