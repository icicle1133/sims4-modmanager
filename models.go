package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ModInfo struct {
	Name        string    `json:"name"`
	InstallDate time.Time `json:"install_date"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
}

type AppSettings struct {
	ModsDirectory string `json:"mods_directory"`
	ApiKey        string `json:"api_key"`
}

var DefaultModsPath = filepath.Join(os.Getenv("HOME"), ".steam", "steam", "steamapps", "compatdata", "1222670", "pfx", "drive_c", "users", "steamuser", "Documents", "Electronic Arts", "The Sims 4", "Mods")

func LoadSettings() (AppSettings, error) {
	settings := AppSettings{
		ModsDirectory: DefaultModsPath,
	}
	
	env := loadEnvFile()
	if apiKey, ok := env["CURSEFORGE_API_KEY"]; ok && apiKey != "" {
		settings.ApiKey = apiKey
	}
	
	data, err := os.ReadFile("settings.json")
	if err != nil {
		if os.IsNotExist(err) {
			return settings, nil
		}
		return settings, err
	}
	
	err = json.Unmarshal(data, &settings)
	
	if settings.ApiKey != "" && (env["CURSEFORGE_API_KEY"] == "" || env["CURSEFORGE_API_KEY"] != settings.ApiKey) {
		updateEnvApiKey(settings.ApiKey)
	}
	
	return settings, err
}

func SaveSettings(settings AppSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile("settings.json", data, 0644)
}