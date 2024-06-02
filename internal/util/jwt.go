package jwt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt/v5"
)

// retrieve JWT key from .env file
var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY")) //TODO privateKey="", fix

// generate JWT token
func GenerateJWT(userId int) (string, error) {
	log.Println(string(privateKey))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(context echo.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func ValidateRole(context echo.Context, id uint) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == id {
		return nil
	}
	return errors.New("invalid token role")
}

// validate Customer role
func ValidateCustomerRoleJWT(context echo.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 2 || userRole == 1 {
		return nil
	}
	return errors.New("invalid author token provided")
}

// fetch user details from the token
func GetUserIdFromContext(context echo.Context) int {
	token, _ := getToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := int(claims["id"].(float64))

	return userId
}

// check token validity
func getToken(context echo.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

// extract token from request Authorization header
func getTokenFromRequest(context echo.Context) string {
	bearerToken := context.Request().Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
