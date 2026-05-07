package fileMetadata

import "gorm.io/gorm"

type FileMetadata struct {
	gorm.Model
	Name      string    `gorm:"index"`
	CID       string    `gorm:"uniqueIndex"`
	EncryptionCID string
	OwnerID   uint      `gorm:"index"`
	Policy string 
}

type FileListDto struct{
	ID uint	`json:"id"`
	Name string	`json:"name"`
}

type FileDetailsDto struct{
	ID uint 	`json:"id"`
	Name string	`json:"name"`
	CID       string    `json:"cid"`
	EncryptionCID string `json:"encryptionCID"`
}

type FileUploadDto struct{
	Name string	`json:"name"`
	CID       string    `json:"cid"`
	EncryptionCID string `json:"encryptionCID"`
	Policy string `json:"policy"`
}