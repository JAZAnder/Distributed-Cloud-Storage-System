package httpRequests

import (
	"encoding/json"

)

func GetDirectory() ([]FileListDto,error)  {
	body, err := CoordinatorRequests("GET", "/api/metadata", "")

	if err != nil {
		return []FileListDto{}, err
	}

	var files []FileListDto
	json.Unmarshal(body, &files)

	return files, nil

}

type FileListDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}