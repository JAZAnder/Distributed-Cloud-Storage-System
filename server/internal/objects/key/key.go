package key

type CryptoConfig struct {
    ID              uint   `gorm:"primaryKey"`
    PublicParams    []byte 
    EncryptedMSY []byte 
}

type ConfigResponse struct {
	PublicParams string `json:"public_params"`
}