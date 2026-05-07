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

func CoordinatorRequests(method, path, data string) ([]byte, error) {
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
		return []byte{}, err
	}
	req.Header.Add("Authorization", "Bearer "+currentSession.GetToken())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	//fmt.Println(string(body))
	return body, nil

}

func NodeDownloadRequest(cid string) ([]byte, error) {
	demo, err := strconv.ParseBool(os.Getenv("DEMO"))
	if err != nil {
		demo = false
	}
	currentSession := session.GetSession()

	url := currentSession.Node.DownloadloadNodeURL + "/ipfs/" + cid
	client := &http.Client{}
	if demo {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	//fmt.Println(string(body))
	return body, nil

}
