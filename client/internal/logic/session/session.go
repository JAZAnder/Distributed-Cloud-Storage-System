package session

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type session struct {
	jwt_SECRET          string
	coordinatorPassword string
	Coordinator         coordinatorConfig
	Node                nodeConfig
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
	UserName string `json:"username"`
	NodeURL  string `json:"coordinator_url"`
}

const configFileName = "config.json"

var once sync.Once
var currentSession session

func GetSession() session {
	once.Do(func() {
		createSession()
	})
	return currentSession
}

func createSession() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Coordinator (Master) URL: ")
	scanner.Scan()
	currentSession.Coordinator.CoordinatorURL = "https://" + strings.TrimSpace(scanner.Text())

	fmt.Print("Enter your Coordinator Username: ")
	scanner.Scan()
	currentSession.Coordinator.UserName = strings.TrimSpace(scanner.Text())

	return

}

func Connect() {

	demo, err := strconv.ParseBool(os.Getenv("DEMO"))
	if err != nil {
		demo = false
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your Coordinator Password: ")
	scanner.Scan()
	currentSession.coordinatorPassword = strings.TrimSpace(scanner.Text())

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
	fmt.Println(string(body))

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
