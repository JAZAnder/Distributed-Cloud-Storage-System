package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/authenticator"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/database"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/responses"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/fileMetadata"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/securityLog"

)

func getMyMetadata(w http.ResponseWriter, r *http.Request) {
	db := database.GetDatabase()
	var files []fileMetadata.FileMetadata
	var fileInfo []fileMetadata.FileListDto

	currentUser, err := authenticator.Identify(*r)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusUnauthorized, err.Error())
	}

	err = db.Where("owner_id = ?", currentUser.ID).Find(&files).Error
	if err != nil {
		responses.RespondWithError(r, w, http.StatusInternalServerError, err.Error())
		return
	}

	logEntry := securityLog.SecurityLog{
		Principal:  currentUser.Username,
		Action:     "LIST_FILES",
		ResourceID: "ALL_FILES",
		Details:    fmt.Sprintf("User %s listed the directory", currentUser.Username),
		IPAddress:  r.RemoteAddr,
	}
	db.Create(&logEntry)

	for _, file := range files {
		fileInfo = append(fileInfo, fileMetadata.FileListDto{
			ID:   file.ID,
			Name: file.Name,
		})
	}

	responses.RespondWithJSON(r, w, http.StatusOK, fileInfo)
}

func getMetadataById(w http.ResponseWriter, r *http.Request) {
	db := database.GetDatabase()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responses.RespondWithError(r, w, http.StatusBadRequest, "ID must int")
		return
	}

	currentUser, err := authenticator.Identify(*r)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusUnauthorized, err.Error())
		return
	}

	file := fileMetadata.FileMetadata{OwnerID: currentUser.ID, Model: gorm.Model{ID: uint(id)}}
	err = db.Where(&file).First(&file).Error
	if err != nil {
		responses.RespondWithError(r, w, http.StatusNotFound, err.Error())
		return
	}

	logEntry := securityLog.SecurityLog{
		Principal:  currentUser.Username,
		Action:     "GET_FILES",
		ResourceID: strconv.Itoa(int(file.ID)),
		Details:    fmt.Sprintf("User %s requested %s", currentUser.Username, file.Name),
		IPAddress:  r.RemoteAddr,
	}
	db.Create(&logEntry)

	fileDto := fileMetadata.FileDetailsDto{
		ID:            file.ID,
		Name:          file.Name,
		CID:           file.CID,
		EncryptionCID: file.EncryptionCID,
	}

	responses.RespondWithJSON(r, w, http.StatusOK, fileDto)
	return
}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	var input fileMetadata.FileUploadDto
	db := database.GetDatabase()

	currentUser, err := authenticator.Identify(*r)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusUnauthorized, err.Error())
		return
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		responses.RespondWithError(r, w, http.StatusBadRequest, "Malformed metadata JSON")
		return
	}

	var metadataObj = fileMetadata.FileMetadata{
		Name:          input.Name,
		CID:           input.CID,
		EncryptionCID: input.EncryptionCID,
		Policy:        input.Policy,
	}
	metadataObj.OwnerID = currentUser.ID

	result := db.Create(&metadataObj)
	secLog := securityLog.SecurityLog{
		Principal:  currentUser.Username,
		Action:     "UPLOAD_FILE",
		ResourceID: strconv.Itoa(int(metadataObj.ID)),
		IPAddress:  r.RemoteAddr,
		Details:    currentUser.Username + "Uploaded Metadata for the file: " + metadataObj.Name,
	}
	db.Create(&secLog)
	if result.Error != nil {
		responses.RespondWithError(r, w, http.StatusInternalServerError, "Database persistence failed")
		return
	}

	responses.RespondWithJSON(r, w, http.StatusOK, fileMetadata.FileDetailsDto{
		ID: metadataObj.ID,
		Name: metadataObj.Name,
		CID: metadataObj.CID,
		EncryptionCID: metadataObj.EncryptionCID,
	})
	return
}
