package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	secretKey := []byte(os.Getenv("APP_SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (uint, error) {
	secretKey := []byte(os.Getenv("APP_SECRET"))

	jwtString := tokenString
	if strings.HasPrefix(jwtString, "Bearer ") {
		jwtString = strings.Split(tokenString, "Bearer ")[1]
	}

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad_signature"+" : %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["userID"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid")
}