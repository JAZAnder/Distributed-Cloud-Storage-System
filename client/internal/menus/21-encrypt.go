package menus

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fentec-project/gofe/abe"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/httpRequests"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/objects/fileMetadata"

)

func encryptAndUpload(fileName string) error {
	filePath := filepath.Join("./downloads", fileName)
	data, _ := os.ReadFile(filePath)

	sessionKey := make([]byte, 32)
	rand.Read(sessionKey)

	block, _ := aes.NewCipher(sessionKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	fileCiphertext := gcm.Seal(nonce, nonce, data, nil)
	pubKey, err := httpRequests.GetPublicPrams()
	if err != nil {
		return err
	}

	fame := abe.NewFAME()
	policy := "(Role:ADMIN)"
	msp, err := abe.BooleanToMSP(policy, false)
	abeCiphertext, _ := fame.Encrypt(string(sessionKey), msp, pubKey)

	serializedAbeKey, _ := json.Marshal(abeCiphertext)

	cid, _ := httpRequests.UploadToIPFS(fileCiphertext)
	fmt.Printf("File uploaded to IPFS. CID: %s\n", cid)
	encryptionCid, _ := httpRequests.UploadToIPFS(serializedAbeKey)
	fmt.Printf("Keys uploaded to IPFS. CID: %s\n", cid)

	metadata := fileMetadata.FileUploadDto{
		Name:          fileName,
		CID:           cid,
		EncryptionCID: encryptionCid,
		Policy:        policy,
	}

	uploadToCoordinator(metadata)
	if err != nil {
		return err
	}

	return nil
}

func uploadToCoordinator(newFile fileMetadata.FileUploadDto) error {
	cs := session.GetSession()
	body, _ := json.Marshal(newFile)
	resp, err := http.Post(cs.Coordinator.CoordinatorURL+"/api/metadata", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to reach coordinator")
	}
	return nil
}
