package quickLog

import (
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/helpers/database"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/objects/securityLog"
)

func Log(principal, action, resourceID, ipAddress, fromValue, toValue, details string) {
	db := database.GetDatabase()
	logEntry := securityLog.SecurityLog{
		Principal:  principal,
		Action:     action,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		FromValue:  fromValue,
		ToValue:    toValue,
		Details:    details,
	}
	db.Create(&logEntry)
}
