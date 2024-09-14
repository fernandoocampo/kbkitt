package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/filesystems"
	yaml "gopkg.in/yaml.v3"
)

/*
version: 0.1.0
server:
  url: http://localhost:8080
*/

type Configuration struct {
	Version             string  `yaml:"version"`
	FilepathForSyncPath string  `yaml:"fileForSyncPath"`
	Server              *Server `yaml:"server"`
}

type Server struct {
	URL string `yaml:"url"`
}

const (
	folderName = ".kbkitt"
	fileName   = "config.yaml"
)

func (c *Configuration) Invalid() bool {
	return c == nil ||
		c.Server == nil ||
		c.Server.URL == "" ||
		c.FilepathForSyncPath == ""
}

func LoadConfiguration() (*Configuration, error) {
	filePath, err := getKBKittConfigurationPath()
	if err != nil {
		return nil, fmt.Errorf("unable to load configuration: %w", err)
	}

	yamlFile, err := filesystems.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to load configuration: %w", err)
	}

	if yamlFile == nil {
		return nil, nil
	}

	var configuration Configuration

	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to load configuration: %w", err)
	}

	return &configuration, nil
}

func CheckAndCreateKBKittFolder() error {
	kbkittDir, err := getKBKittFolderPath()
	if err != nil {
		return fmt.Errorf("unable to get kbkitt folder path: %w", err)
	}

	exist, err := configurationFolderExist(kbkittDir)
	if err != nil {
		return fmt.Errorf("unable to check if kbkitt folder (%q) exist: %w", kbkittDir, err)
	}

	if exist {
		return nil
	}

	err = filesystems.MakeFolder(kbkittDir)
	if err != nil {
		return fmt.Errorf("unable to create kbkitt folder (%q): %w", kbkittDir, err)
	}

	return nil
}

func Save(newConf *Configuration) error {
	yamlFile, err := yaml.Marshal(newConf)
	if err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	path, err := getKBKittConfigurationPath()
	if err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	err = filesystems.SaveFile(path, yamlFile)
	if err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	return nil
}

func getKBKittConfigurationPath() (string, error) {
	kbKittDir, err := getKBKittFolderPath()
	if err != nil {
		return "", fmt.Errorf("unable to get kbkitt folder path")
	}

	return filepath.Join(kbKittDir, fileName), nil
}

func getKBKittFolderPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get kbkitt folder name: %w", err)
	}

	return filepath.Join(homeDir, folderName), nil
}

func configurationFolderExist(kbkittFolder string) (bool, error) {
	_, err := os.Stat(kbkittFolder)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("unable to check if configuration folder exists: %w", err)
	}

	return true, nil
}
