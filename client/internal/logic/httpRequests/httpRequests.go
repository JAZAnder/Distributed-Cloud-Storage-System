package httpRequests

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"

)

func CoordinatorRequests(method, path, data string) {
	demo, err := strconv.ParseBool(os.Getenv("DEMO"))
	if err != nil {
		demo = false
	}
	currentSession := session.GetSession()
	url := currentSession.Coordinator.CoordinatorURL + path
	payload := strings.NewReader(data)
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
	req.Header.Add("Authorization", "Bearer "+ currentSession.GetToken())
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
