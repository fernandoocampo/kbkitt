package filesystems

import (
	"fmt"
	"io"
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

func FileEmpty(filePath string) (bool, error) {
	stat, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to read file info: %w", err)
	}

	if stat.Size() == 0 {
		return true, nil
	}

	return false, nil
}

func SaveFile(filePath string, content []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	_, err = io.WriteString(f, string(content))
	if err != nil {
		return fmt.Errorf("unable to save file: %w", err)
	}

	return nil
}

func SaveOrAppendFile(filePath string, content []byte) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, drwxr_xr_x)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			return fmt.Errorf("unable to save file: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("unable to open file [%q]: %w", filePath, err)
	}

	defer file.Close()

	_, err = io.WriteString(file, string(content))
	if err != nil {
		return fmt.Errorf("unable to write into file: %w", err)
	}

	return nil
}

func TruncateFile(filePath string) error {
	err := os.Truncate(filePath, 0)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("unable to truncate file: %w", err)
	}

	return nil
}
