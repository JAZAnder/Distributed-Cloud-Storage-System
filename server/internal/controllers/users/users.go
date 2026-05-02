package users

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/database"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/quickLog"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/responses"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/securityLog"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/user"

)

func login(w http.ResponseWriter, r *http.Request) {
	db := database.GetDatabase()

	logEntry := securityLog.SecurityLog{
		Action:    "LoginAttempt",
		IPAddress: r.RemoteAddr,
	}

	u := user.UserLoginDto{
		Username: r.PostFormValue("userName"),
		Password: r.PostFormValue("password"),
	}

	var userObj user.User
	if err := db.Where("username = ?", u.Username).First(&userObj).Error; err != nil {
		logEntry.Details = "Failed: " + u.Username + " not found"
		db.Create(&logEntry)

		responses.RespondWithError(r, w, http.StatusUnauthorized, "Username or Password Incorrect")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userObj.PasswordHash), []byte(u.Password)); err != nil {
		logEntry.ResourceID = strconv.FormatUint(uint64(userObj.ID), 10)
		logEntry.Details = "Failed: incorrect password for " + userObj.Username
		db.Create(&logEntry)

		responses.RespondWithError(r, w, http.StatusUnauthorized, "Username or Password Incorrect")

		return
	}

	logEntry.Principal = userObj.Username
	logEntry.ResourceID = strconv.FormatUint(uint64(userObj.ID), 10)

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Println("WARNING JWT_SECRET ENV NOT SET. Security may be Compromised")
		secret = "development_secret"
	}

	var jwtKey = []byte(secret)
	expirationTime := time.Now().Add(time.Hour)

	claim := user.JWTClaim{
		UserID: userObj.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logEntry.Details = "Failed: " + err.Error()
		db.Create(&logEntry)
		responses.RespondWithError(r, w, http.StatusInternalServerError, err.Error())
		return
	} else {
		quickLog.Log(userObj.Username,"ClaimGeneration",strconv.FormatUint(uint64(userObj.ID), 10),r.RemoteAddr,"","","Expiration: "+claim.ExpiresAt.String())
	}
	logEntry.Details = "Success"
	db.Create(&logEntry)

	responses.RespondWithJSONNoLog(w, http.StatusOK, map[string]string{"Claim": tokenString})

}

func CreateUser(requestor string, newUser user.UserCreateDto) (user.User, error) {
	logEntry := securityLog.SecurityLog{
		Principal: requestor,
		Action:    "CreateUser",
	}
	userObj := user.User{
		Username: newUser.Username,
	}
	db := database.GetDatabase()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		logEntry.Details = "Failed:" + err.Error()
		db.Create(&logEntry)
		return userObj, err
	}

	userObj.PasswordHash = string(hashedPassword)

	result := db.Create(&userObj)
	logEntry.ResourceID = strconv.FormatUint(uint64(userObj.ID), 10)

	if result.Error != nil {
		logEntry.Details = "Failed:" + err.Error()
		db.Create(&logEntry)
		return userObj, result.Error
	}

	db.Create(&logEntry)

	return userObj, nil
}
