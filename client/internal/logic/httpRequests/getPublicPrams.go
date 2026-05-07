package httpRequests

import (
	"encoding/base64"
	"encoding/json"

	"github.com/fentec-project/gofe/abe"

)

func GetPublicPrams() (*abe.FAMEPubKey, error) {

	body, err := CoordinatorRequests("GET", "/api/config", "")
	pubKey := new(abe.FAMEPubKey)

	if err != nil {
		return pubKey, err
	}

	var configResp configResponse
	err = json.Unmarshal(body, &configResp)
	if err != nil {
		return pubKey, err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(configResp.PublicParams)
	if err != nil {
		return pubKey, err
	}

	err = json.Unmarshal(pubKeyBytes, pubKey)

	return pubKey, err

}

type configResponse struct {
	PublicParams string `json:"public_params"`
}
