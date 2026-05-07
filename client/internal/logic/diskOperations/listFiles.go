package diskOperations

import (
	"fmt"
	"os"
	"strconv"

)

func ListFilesForUpload() ([]string, error) {
	downloadDir := "./downloads"
	var folders []string

	entries, err := os.ReadDir(downloadDir)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) == 0 {
		return []string{}, fmt.Errorf("no files found in %s", downloadDir)
	}

	fmt.Println("Available files for upload:")
	for i, entry := range entries {
		if !entry.IsDir() {
			fmt.Println("\t[" + strconv.Itoa(i) + "]\t" + entry.Name())
			folders = append(folders, entry.Name())
		}
	}

	return folders, nil
}
