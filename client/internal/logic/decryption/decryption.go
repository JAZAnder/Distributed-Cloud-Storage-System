package decryption

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fentec-project/gofe/abe"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/httpRequests"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"

)

func UnencryptKey(encryptedSymKey []byte) (string, error) {
	session.GetPrivateKey()
	pubKey, err := httpRequests.GetPublicPrams()
	var userPrivateKey *abe.FAMEAttribKeys

	if err != nil {
		return "", err
	}

	privateKey, err := session.GetPrivateKey()
	if privateKey == nil {
		fmt.Print("No identity key found. Attempt to generate new key from Coordinator? (y/n): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if strings.TrimSpace(strings.ToLower(input)) == "y" {
			fmt.Println("Requesting new identity key from Coordinator...")
			userPrivateKey, err = httpRequests.RequestAttrKey()
			if err != nil {
				return "", err
			}
		}
	}
	var abeCipherText *abe.FAMECipher
	json.Unmarshal(encryptedSymKey, abeCipherText)
	fame := abe.NewFAME()

	decryptedEncryptionKey, err := fame.Decrypt(abeCipherText, userPrivateKey, pubKey)
	if err != nil {
		return "", fmt.Errorf("Access Denied: Attributes do not match policy")
	}
	return decryptedEncryptionKey, nil
}

func DecodeData(fileData []byte, encryptionKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create cipher: " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("Failed to create GCM: "+err.Error())
	}
	nonceSize := gcm.NonceSize()
	if len(fileData) < nonceSize {
		return nil, fmt.Errorf("encrypted fileData is too short or malformed")
	}
	nonce, cipherText := fileData[:nonceSize], fileData[nonceSize:]
	data, err := gcm.Open(nil, nonce, cipherText, nil)

	return data, err

}
