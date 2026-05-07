package httpRequests

import (
	"encoding/json"
	"strconv"

	filemetadata "github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/objects/fileMetadata"

)

func downloadAndDecode(fileId int) (filemetadata.FileDetailsDto, error) {
	body, err := CoordinatorRequests("GET", "api/metadata/"+strconv.Itoa(fileId), "")
	var files filemetadata.FileDetailsDto
	if err != nil {
		return filemetadata.FileDetailsDto{}, err
	}
	json.Unmarshal(body, &files)

	//encryptedData, err := NodeDownloadRequest(files.CID)
	if err != nil {
		return files, err
	}
	//encryptedSymKey, err := NodeDownloadRequest(files.EncryptionCID)
	if err != nil {
		return files, err
	}

	return filemetadata.FileDetailsDto{}, err //REMOVE

}
