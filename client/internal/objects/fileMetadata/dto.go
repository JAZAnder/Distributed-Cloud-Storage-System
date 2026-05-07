package fileMetadata

type FileDetailsDto struct{
	ID uint 	`json:"id"`
	Name string	`json:"name"`
	CID       string    `json:"cid"`
	EncryptionCID string `json:"encryptionCID"`
}

type FileListDto struct{
	ID uint	`json:"id"`
	Name string	`json:"name"`
}

type FileUploadDto struct{
	Name string	`json:"name"`
	CID       string    `json:"cid"`
	EncryptionCID string `json:"encryptionCID"`
	Policy string `json:"policy"`
}