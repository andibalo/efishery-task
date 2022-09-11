package middleware

import (
	"errors"
	"fetch-app/internal/model"
	"fetch-app/internal/response"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckJWTToken() gin.HandlerFunc {

	return func(c *gin.Context) {

		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			err := errors.New("JWT Token not found in Authorization header")
			log.Println(err)

			resp := response.NewResponse(response.Unauthorized, nil)

			resp.SetResponseMessage("JWT Token not found in Authorization header")

			c.JSON(http.StatusUnauthorized, resp)

			c.Abort()
			return
		}

		tokenString := strings.Split(bearerToken, " ")

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_SECRET")), nil
		})
		if err != nil {
			log.Println(err)
			err = errors.New("Failed to parse JWT Token/Invalid JWT Token")

			resp := response.NewResponse(response.ServerError, nil)

			resp.SetResponseMessage("Failed to parse JWT Token/Invalid JWT Token")

			c.JSON(http.StatusInternalServerError, resp)

			c.Abort()
			return
		}

		if !token.Valid {
			resp := response.NewResponse(response.Unauthorized, nil)

			resp.SetResponseMessage("JWT Token Invalid")

			c.JSON(http.StatusUnauthorized, resp)

			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckJWTTokenAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		var claim model.Claims

		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			err := errors.New("JWT Token not found in Authorization header")
			log.Println(err)

			resp := response.NewResponse(response.Unauthorized, nil)

			resp.SetResponseMessage("JWT Token not found in Authorization header")

			c.JSON(http.StatusUnauthorized, resp)

			c.Abort()
			return
		}

		tokenString := strings.Split(bearerToken, " ")

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_SECRET")), nil
		})
		if err != nil {
			log.Println(err)
			err = errors.New("Failed to parse JWT Token/Invalid JWT Token")

			resp := response.NewResponse(response.ServerError, nil)

			resp.SetResponseMessage("Failed to parse JWT Token/Invalid JWT Token")

			c.JSON(http.StatusInternalServerError, resp)

			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			claim.Role = fmt.Sprintf("%v", claims["role"])
		} else {
			log.Println(err)
			err = errors.New("Failed to parse private claims")

			resp := response.NewResponse(response.ServerError, nil)

			resp.SetResponseMessage("Failed to parse private claims")

			c.JSON(http.StatusInternalServerError, resp)

			c.Abort()
			return
		}

		if claim.Role != "admin" {
			err = errors.New("Invalid Role")
			log.Println(err)

			resp := response.NewResponse(response.Unauthorized, nil)

			resp.SetResponseMessage("Invalid Role")

			c.JSON(http.StatusUnauthorized, resp)

			c.Abort()

			return
		}

		c.Next()
	}
}
