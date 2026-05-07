package diskOperations

import (
	"fmt"
	"os"
	"path/filepath"
)

func SaveFileToDownloads(fileName string, folder string, decryptedData []byte) error {
	downloadDir := "./downloads/"

	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create downloads directory: %w", err)
	}

	filePath := filepath.Join(downloadDir, fileName)

	// 4. Verification Step (Architectural Requirement) [1]
	// In your design, the client should recompute the hash of the data
	// here to ensure it matches the CID requested from the IPFS "Muscle".
	// if !verifyIntegrity(decryptedData, expectedCID) {
	//     return errors.New("data integrity check failed: CID mismatch")
	// }

	err = os.WriteFile(filePath, decryptedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file to disk: %w", err)
	}

	fmt.Printf("Successfully saved %s to %s\n", fileName, filePath)
	return nil
}
