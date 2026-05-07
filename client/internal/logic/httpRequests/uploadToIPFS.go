package httpRequests

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"

)

func UploadToIPFS(content []byte) (string, error) {
	currentSession := session.GetSession()
	resp, err := http.Post(currentSession.Coordinator.CoordinatorURL+"/api/v0/add", "application/octet-stream", bytes.NewReader(content))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ipfsResp struct {
		Hash string `json:"Hash"`
	}
	json.NewDecoder(resp.Body).Decode(&ipfsResp)
	return ipfsResp.Hash, nil
}
