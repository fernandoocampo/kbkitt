package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInvalid(t *testing.T) {
	tests := []struct {
		name   string
		config *Configuration
		want   bool
	}{
		{
			name:   "nil configuration",
			config: nil,
			want:   true,
		},
		{
			name:   "nil server",
			config: &Configuration{},
			want:   true,
		},
		{
			name: "empty server url",
			config: &Configuration{
				Server: &Server{},
			},
			want: true,
		},
		{
			name: "missing file for sync path",
			config: &Configuration{
				Server:          &Server{URL: "http://localhost:8080"},
				DirForMediaPath: "/some/path",
			},
			want: true,
		},
		{
			name: "missing dir for media path",
			config: &Configuration{
				Server:          &Server{URL: "http://localhost:8080"},
				FileForSyncPath: "/some/path",
			},
			want: true,
		},
		{
			name: "valid configuration",
			config: &Configuration{
				Server:          &Server{URL: "http://localhost:8080"},
				FileForSyncPath: "/some/path",
				DirForMediaPath: "/some/dir",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.Invalid()
			if got != tt.want {
				t.Errorf("Invalid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveAppliesDefaultServerURL(t *testing.T) {
	tmpDir := t.TempDir()

	conf := &Configuration{
		Version:          "0.1.0",
		KBKittFolderPath: tmpDir,
		Server:           &Server{URL: ""},
	}

	err := Save(conf)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	if conf.Server.URL != defaultServerURL {
		t.Errorf("Save() Server.URL = %q, want %q", conf.Server.URL, defaultServerURL)
	}
}

func TestSaveAppliesDefaultServerURLWhenServerIsNil(t *testing.T) {
	tmpDir := t.TempDir()

	conf := &Configuration{
		Version:          "0.1.0",
		KBKittFolderPath: tmpDir,
		Server:           nil,
	}

	err := Save(conf)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	if conf.Server == nil {
		t.Fatal("Save() Server is nil, expected default server to be set")
	}

	if conf.Server.URL != defaultServerURL {
		t.Errorf("Save() Server.URL = %q, want %q", conf.Server.URL, defaultServerURL)
	}
}

func TestSavePreservesProvidedServerURL(t *testing.T) {
	tmpDir := t.TempDir()
	customURL := "http://myserver:9090"

	conf := &Configuration{
		Version:          "0.1.0",
		KBKittFolderPath: tmpDir,
		Server:           &Server{URL: customURL},
	}

	err := Save(conf)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	if conf.Server.URL != customURL {
		t.Errorf("Save() Server.URL = %q, want %q", conf.Server.URL, customURL)
	}
}

func TestSaveAppliesDefaultsForEmptyPaths(t *testing.T) {
	tmpDir := t.TempDir()

	conf := &Configuration{
		Version:          "0.1.0",
		KBKittFolderPath: tmpDir,
		Server:           &Server{URL: "http://localhost:8080"},
		FileForSyncPath:  "",
		DirForMediaPath:  "",
	}

	err := Save(conf)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	expectedSync := filepath.Join(tmpDir, syncFileName)
	if conf.FileForSyncPath != expectedSync {
		t.Errorf("Save() FileForSyncPath = %q, want %q", conf.FileForSyncPath, expectedSync)
	}

	expectedMedia := filepath.Join(tmpDir, mediaFolderName)
	if conf.DirForMediaPath != expectedMedia {
		t.Errorf("Save() DirForMediaPath = %q, want %q", conf.DirForMediaPath, expectedMedia)
	}
}

func TestSaveWritesConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	conf := &Configuration{
		Version:          "0.1.0",
		KBKittFolderPath: tmpDir,
		Server:           &Server{URL: ""},
		FileForSyncPath:  "",
		DirForMediaPath:  "",
	}

	err := Save(conf)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	configPath := filepath.Join(tmpDir, fileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Save() did not create config file at %q", configPath)
	}
}
