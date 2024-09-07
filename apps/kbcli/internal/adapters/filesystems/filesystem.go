package filesystems

import (
	"fmt"
	"os"
)

const drwxr_xr_x os.FileMode = 0755

func MakeFolder(folderPath string) error {
	// https://chmod-calculator.com
	err := os.Mkdir(folderPath, drwxr_xr_x)
	if err != nil {
		return fmt.Errorf("unable to make directory: %w", err)
	}

	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return file, nil
}
