package session

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/fentec-project/gofe/abe"
	"golang.org/x/term"

)

type session struct {
	jwt_SECRET          string
	coordinatorPassword string
	Coordinator         coordinatorConfig
	Node                nodeConfig
	PrivateKey          *abe.FAMEAttribKeys
	PublicKey           *abe.FAMEPubKey
}

type configurationFile struct {
	coordinators []coordinatorConfig
	nodes        []nodeConfig
}

type coordinatorConfig struct {
	UserName       string `json:"username"`
	CoordinatorURL string `json:"coordinator_url"`
}

type nodeConfig struct {
	UploadNodeURL       string `json:"upload_node_url"`
	DownloadloadNodeURL string `json:"Downloadload_node_url"`
}

const configFileName = "config.json"

var once sync.Once
var currentSession session

type ConnectionResponse struct {
	Claim string `json:"Claim"`
}

func (cs *session) GetToken() string {
	return cs.jwt_SECRET
}

func GetSession() session {
	once.Do(func() {
		createSession()
	})
	return currentSession
}

func createSession() {
	scanner := bufio.NewScanner(os.Stdin)
	if os.Getenv("COORDINATOR") == "" {
		fmt.Print("\nEnter the Coordinator (Master) URL: ")
		scanner.Scan()

		currentSession.Coordinator.CoordinatorURL = "https://" + strings.TrimSpace(scanner.Text())

	} else {
		currentSession.Coordinator.CoordinatorURL = "https://" + os.Getenv("COORDINATOR")
	}

	if os.Getenv("USERNAME") == "" {
		fmt.Print("\nEnter your Coordinator Username: ")
		scanner.Scan()

		currentSession.Coordinator.UserName = "https://" + strings.TrimSpace(scanner.Text())

	} else {
		currentSession.Coordinator.UserName = os.Getenv("COORDINATORUSERNAME")
	}

	if os.Getenv("DOWNLOADNODE") == "" {
		fmt.Print("\nEnter the Download IPFS Node URL: ")
		scanner.Scan()

		currentSession.Coordinator.UserName = "https://" + strings.TrimSpace(scanner.Text())

	} else {
		currentSession.Node.DownloadloadNodeURL = os.Getenv("DOWNLOADNODE")
	}

	if os.Getenv("UPLOADNODE") == "" {
		fmt.Print("\nEnter the Upload IPFS Node URL: ")
		scanner.Scan()

		currentSession.Coordinator.UserName = "https://" + strings.TrimSpace(scanner.Text())

	} else {
		currentSession.Node.UploadNodeURL = os.Getenv("UPLOADNODE")
	}

	return

}

func Connect() {

	demo, err := strconv.ParseBool(os.Getenv("DEMO"))
	if err != nil {
		demo = false
	}
	if demo {
		fmt.Println("\n\n- - THE CLIENT IS IN DEMO MODE!! -- TLS Certs will NOT be Verified - - \n\n")
	}

	fmt.Print("Enter your Coordinator Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln("\nError reading password:", err)
		return
	}

	currentSession.coordinatorPassword = string(bytePassword)
	fmt.Println("\n")

	url := currentSession.Coordinator.CoordinatorURL + "/api/login"
	method := "POST"

	payload := strings.NewReader("userName=" + currentSession.Coordinator.UserName + "&password=" + currentSession.coordinatorPassword)

	client := &http.Client{}
	if demo {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}

	}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response ConnectionResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalln("Error parsing JSON: %v", err)
	}

	jwtToken := response.Claim
	if jwtToken == "" {
		log.Fatalln("Login Failed - Exiting Program\n")
	}

	currentSession.jwt_SECRET = jwtToken
	fmt.Println("Extracted Token:", jwtToken)
	currentSession.coordinatorPassword = " "

}

func loadSession() error {
	var configFile configurationFile
	scanner := bufio.NewScanner(os.Stdin)

	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println("Searching for Config File...")
		fileData, err := os.ReadFile(configFileName)
		if err == nil {
			json.Unmarshal(fileData, &configFile)
			fmt.Println("Configuration File Found...")
			for i, coordinator := range configFile.coordinators {
				fmt.Println("Configuration: " + strconv.Itoa(i))
				fmt.Println("UserName: " + coordinator.UserName)
				fmt.Println("Configuration: " + coordinator.CoordinatorURL)
			}

			fmt.Print("Would you like to load an existing following configuration? (y/n): ")
			scanner.Scan()
			choice := strings.ToLower(strings.TrimSpace(scanner.Text()))

			if choice == "y" || choice == "yes" {
				for i := 0; i < 5; i++ {
					fmt.Print("Which would you like to load? (1/2/3): ")
					scanner.Scan()
					choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
					if err == nil {
						currentSession.Coordinator = configFile.coordinators[choice]
						i = 6
					} else {
						fmt.Println(err)
						err = nil
					}
				}

			}

		}
	}
	if currentSession.Coordinator.CoordinatorURL == "" {
		//Get Information and Offer to Save
	}
	return nil
}

func GetPrivateKey() (*abe.FAMEAttribKeys, error) {
	const keyPath = "user_identity.key"

	if currentSession.PrivateKey != nil {
		return currentSession.PrivateKey, nil
	} else {
		if _, err := os.Stat(keyPath); err == nil {
			fmt.Println("Found existing identity key. Loading...")
			return loadAndDecryptKey(keyPath)
		}

	}
	return nil, fmt.Errorf("No AttrKey Available")
}

func loadAndDecryptKey(keyPath string) (*abe.FAMEAttribKeys, error) {
	fileData, _ := os.ReadFile(keyPath)

	fmt.Print("Enter passphrase to unlock your private key: ")
	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	key := sha256.Sum256(passphrase)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(fileData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := fileData[:nonceSize], fileData[nonceSize:]
	decryptedJson, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("incorrect passphrase: %w", err)
	}

	attribKeys := new(abe.FAMEAttribKeys)
	err = json.Unmarshal(decryptedJson, attribKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABE keys: %w", err)
	}
	passphrase = nil

	return attribKeys, nil
}
