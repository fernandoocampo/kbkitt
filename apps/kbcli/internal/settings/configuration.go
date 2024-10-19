package settings

import (
	"context"
	"fmt"
	"log/slog"
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

type Storage interface {
	InitializeDB(ctx context.Context) error
}

type Configuration struct {
	Version          string  `yaml:"version"`
	KBKittFolderPath string  `yaml:"-"`
	FileForSyncPath  string  `yaml:"fileForSyncPath"`
	DirForMediaPath  string  `yaml:"dirForMediaPath"`
	Server           *Server `yaml:"server"`
}

type Server struct {
	URL string `yaml:"url"`
}

const (
	folderName      = ".kbkitt"
	mediaFolderName = "media"
	fileName        = "config.yaml"
	syncFileName    = "sync.yaml"
	dbName          = "kbkitt.db"
)

func (c *Configuration) Invalid() bool {
	return c == nil ||
		c.Server == nil ||
		c.Server.URL == "" ||
		c.FileForSyncPath == "" ||
		c.DirForMediaPath == ""
}

func LoadConfiguration() (*Configuration, error) {
	kbKittFolderPath, err := getKBKittFolderPath()
	if err != nil {
		return nil, fmt.Errorf("unable to get kbkitt folder path")
	}

	yamlFile, err := filesystems.ReadFile(getKBKittConfigurationPath(kbKittFolderPath))
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

	configuration.KBKittFolderPath = kbKittFolderPath

	return &configuration, nil
}

func CheckAndCreateKBKittFolder() error {
	kbkittDir, err := getKBKittFolderPath()
	if err != nil {
		return fmt.Errorf("unable to get kbkitt folder path: %w", err)
	}

	exist, err := filesystems.FolderExist(kbkittDir)
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
	err := newConf.setKBKittFolderPath()
	if err != nil {
		return fmt.Errorf("unable to get kbkitt folder path: %w", err)
	}

	if newConf.DirForMediaPath == "" {
		slog.Info("using default media dir", slog.String("file", newConf.getDefaultMediaDir()))
		newConf.DirForMediaPath = newConf.getDefaultMediaDir()
	}

	if newConf.FileForSyncPath == "" {
		slog.Info("using default sync file", slog.String("file", newConf.getDefaultSyncFilePath()))
		newConf.FileForSyncPath = newConf.getDefaultSyncFilePath()
	}

	yamlFile, err := yaml.Marshal(newConf)
	if err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	err = filesystems.SaveFile(newConf.getKBKittConfigurationPath(), yamlFile)
	if err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	return nil
}

func CreateDatabaseIfNotExist(ctx context.Context, newConf *Configuration, storage Storage) error {
	fmt.Println("file db", newConf.GetDBPath())
	fileExist, err := filesystems.FileExists(newConf.GetDBPath())
	if err != nil {
		return fmt.Errorf("unable to create database: %w", err)
	}

	if fileExist {
		return nil
	}

	err = storage.InitializeDB(ctx)
	if err != nil {
		return fmt.Errorf("unable to create database: %w", err)
	}

	return nil
}

func (c Configuration) getKBKittConfigurationPath() string {
	return filepath.Join(c.KBKittFolderPath, fileName)
}

func getKBKittConfigurationPath(kbKittFolderPath string) string {
	return filepath.Join(kbKittFolderPath, fileName)
}

func getKBKittFolderPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get kbkitt folder name: %w", err)
	}

	return filepath.Join(homeDir, folderName), nil
}

func (c Configuration) GetDBPath() string {
	return filepath.Join(c.KBKittFolderPath, dbName)
}

func (c Configuration) getDefaultMediaDir() string {
	return filepath.Join(c.KBKittFolderPath, mediaFolderName)
}

func (c Configuration) getDefaultSyncFilePath() string {
	return filepath.Join(c.KBKittFolderPath, syncFileName)
}

func (c *Configuration) setKBKittFolderPath() error {
	if c.KBKittFolderPath != "" {
		return nil
	}

	kbKittFolderPath, err := getKBKittFolderPath()
	if err != nil {
		return fmt.Errorf("unable to get kbkitt folder path")
	}

	c.KBKittFolderPath = kbKittFolderPath

	return nil
}
