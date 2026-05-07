package menus

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/decryption"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/diskOperations"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/httpRequests"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/objects/fileMetadata"

)

func downloadAndDecode(fileId int) error {
	body, err := httpRequests.CoordinatorRequests("GET", "api/metadata/"+strconv.Itoa(fileId), "")
	if err != nil {
		return err
	}
	var file fileMetadata.FileDetailsDto
	json.Unmarshal(body, &file)

	fileData, err := httpRequests.NodeDownloadRequest(file.CID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	encryptedSymKey, err := httpRequests.NodeDownloadRequest(file.EncryptionCID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	symKey, err := decryption.UnencryptKey(encryptedSymKey)
	if err != nil {
		return err
	}
	var fileToWrite []byte
	fileToWrite, err = decryption.DecodeData(fileData, []byte(symKey)) 
	if err != nil {
		return err
	}
	err = diskOperations.SaveFileToDownloads(file.Name,file.Name, fileToWrite)
	if err != nil {
		return err
	}
	return nil
}
