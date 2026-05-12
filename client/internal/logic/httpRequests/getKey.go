package httpRequests

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/fentec-project/gofe/abe"
	"golang.org/x/term"

)

func RequestAttrKey() (*abe.FAMEAttribKeys, error) {

	body, err := CoordinatorRequests("POST", "/api?key", "")

	if err != nil {
		return nil, fmt.Errorf("Error: The connected Coordinator lacks the MSK/KEK to mint keys.")
	}

	var keyResp struct {
		PrivateKey string `json:"private_key"`
	}
	err = json.Unmarshal(body, &keyResp)
	abeKeyBytes, _ := base64.StdEncoding.DecodeString(keyResp.PrivateKey)

	fmt.Print("Create a passphrase to protect your local key file: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))

	encryptedFileBlob, _ := locallyEncryptKey(abeKeyBytes, passphrase)
	os.WriteFile("user_identity.key", encryptedFileBlob, 0600)
	fmt.Println("New identity key successfully minted and saved locally.")

	var attrKey *abe.FAMEAttribKeys
	err = json.Unmarshal(abeKeyBytes, attrKey)

	return attrKey, err
}

func locallyEncryptKey(abeKeyBytes []byte, passphrase []byte) ([]byte, error) {
	key := sha256.Sum256(passphrase)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}


	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	encryptedBlob := gcm.Seal(nonce, nonce, abeKeyBytes, nil)

	return encryptedBlob, nil
}
