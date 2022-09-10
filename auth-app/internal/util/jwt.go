package util

import (
	"auth-app/internal/model"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user model.User) (tokenString string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":      user.Name,
		"phone":     user.Phone,
		"role":      user.Role,
		"timestamp": user.Timestampz,
	})

	tokenString, err = token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		err = errors.New("Failed to sign JWT token")
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (claim model.Claims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println(err)
		err = errors.New("Failed to parse JWT Token")
		return claim, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claim.Name = fmt.Sprintf("%v", claims["name"])
		claim.Phone = fmt.Sprintf("%v", claims["phone"])
		claim.Role = fmt.Sprintf("%v", claims["role"])
		claim.Timestampz = fmt.Sprintf("%v", claims["timestamp"])
	} else {
		log.Println(err)
		err = errors.New("Failed to parse private claims")
		return claim, err
	}

	return claim, nil

}
