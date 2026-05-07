package keys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/fentec-project/gofe/abe"

)

var encodedSymKey string

func GetMSYDecryptionKey() ([]byte, error) {
	
	if encodedSymKey == "" {
		encodedSymKey = os.Getenv("SYMKEY")
	}
	if encodedSymKey == "" {
		fmt.Println("No Symmetric Key Found")
		return nil, errors.New("Key Not Found - Node not able to decrypt MSK")
	}

	decodedKey, err := base64.StdEncoding.DecodeString(encodedSymKey)
	if err != nil {
		return nil, err
	}
	return decodedKey, nil
}

func SetMSYDecryptionKey(decodedKey []byte) string {
	encodedSymKey = base64.StdEncoding.EncodeToString(decodedKey)
	return encodedSymKey
}

func KeysToBytes(pubKey *abe.FAMEPubKey, msk *abe.FAMESecKey) ([]byte, []byte, error) {
	pubKeyBytes, err := json.Marshal(pubKey)
	if err != nil {
		return nil,nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	mskBytes, err := json.Marshal(msk)
	if err != nil {
		return nil,nil, fmt.Errorf("failed to marshal Master key: %w", err)
	}

	return pubKeyBytes, mskBytes, nil
}

func EncryptAtRest(msk []byte, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	encryptedMsk := gcm.Seal(nonce, nonce, msk, nil)

	return encryptedMsk, nil
}
