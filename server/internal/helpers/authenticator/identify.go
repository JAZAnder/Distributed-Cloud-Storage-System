package authenticator

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/database"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/securityLog"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/user"
)

func Identify(r http.Request) (user.User, error) {
	var err error
	secret := os.Getenv("JWT_SECRET")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {

		return user.User{}, errors.New("Missing Authorization Header")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &user.JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return user.User{}, errors.New("Missing Authorization Header")
	}

	db := database.GetDatabase()
	var userObj user.User
	if err := db.First(&userObj, claims.UserID).Error; err != nil {
		return user.User{}, errors.New("Missing Authorization Header")
	}

	logEntry := securityLog.SecurityLog{
		Principal:  userObj.Username,
		IPAddress:  r.RemoteAddr,
		Action:     "AUTH_CHECK",
		ResourceID: strconv.FormatUint(uint64(userObj.ID), 10),
		Details:    fmt.Sprintf("User %s has been authenticated", userObj.Username),
	}
	db.Create(&logEntry)

	return userObj, nil
}
