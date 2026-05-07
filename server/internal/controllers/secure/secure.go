package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fentec-project/gofe/abe"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/authenticator"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/database"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/keys"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/quickLog"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/responses"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/key"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/securityLog"

)

func getPublicPrams(w http.ResponseWriter, r *http.Request) {
	var config key.CryptoConfig

	db := database.GetDatabase()
	db.First(&config)

	encodedParams := base64.StdEncoding.EncodeToString(config.PublicParams)
	response := key.ConfigResponse{PublicParams: encodedParams}
	responses.RespondWithJSON(r, w, http.StatusOK, response)
}

func mintUserKey(w http.ResponseWriter, r *http.Request) {
	db := database.GetDatabase()
	currentUser, err := authenticator.Identify(*r)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusUnauthorized, err.Error())
	}

	symKey, err := keys.GetMSYDecryptionKey()
	if err != nil {
		responses.RespondWithError(r, w, http.StatusNotImplemented, "This is not a minting node.")
		return
	}
	quickLog.Log(currentUser.Username, "MSK_Decryption", "", r.RemoteAddr, "", "", "MSK decryption Key Has been accessed")

	var config key.CryptoConfig
	db.First(&config)

	block, err := aes.NewCipher(symKey)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusInternalServerError, "failed to create cipher: "+err.Error())
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusInternalServerError, "failed to create GCM: "+err.Error())
		return
	}

	nonceSize := gcm.NonceSize()
	if len(config.EncryptedMSY) < nonceSize {
		responses.RespondWithError(r, w, http.StatusInternalServerError, "encrypted MSK is too short or malformed")
		return
	}

	nonce, ciphertext := config.EncryptedMSY[:nonceSize], config.EncryptedMSY[nonceSize:]

	plaintextMsk, err := gcm.Open(nil, nonce, ciphertext, nil)
	var mskKey *abe.FAMESecKey
	json.Unmarshal(plaintextMsk, &mskKey)
	plaintextMsk = nil
	if err != nil {
		responses.RespondWithError(r, w, http.StatusInternalServerError, "decryption failed (invalid KEK or tampered data): "+err.Error())
		return
	}

	fame := abe.NewFAME()
	abeKey, _ := fame.GenerateAttribKeys([]string{"Role:ADMIN"}, mskKey)
	mskKey = nil
	symKey = nil

	serializedKey, _ := json.Marshal(abeKey)
	encodedKey := base64.StdEncoding.EncodeToString(serializedKey)

	db.Create(&securityLog.SecurityLog{
		Principal:  currentUser.Username,
		Action:     "MINT_KEY",
		Details:    fmt.Sprintf("Minted ABE key for attributes: ADMIN"),
		IPAddress:  r.RemoteAddr,
		ResourceID: strconv.Itoa(int(currentUser.ID)),
	})

	responses.RespondWithJSON(r, w, http.StatusOK, map[string]string{"private_key": encodedKey})
}
