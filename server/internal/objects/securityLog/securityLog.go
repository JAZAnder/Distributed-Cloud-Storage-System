package securityLog

import "gorm.io/gorm"

type SecurityLog struct {
	gorm.Model               
	Principal string   `gorm:"index"` 
	Action      string `gorm:"index"` 
	ResourceID  string `gorm:"index"`
	IPAddress string `gorm:"index"`
	FromValue   string
	ToValue   string             
	Details     string
}
