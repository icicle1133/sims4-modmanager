package main

import (
	"bufio"
	"os"
	"strings"
)

func loadEnvFile() map[string]string {
	env := make(map[string]string)
	
	file, err := os.Open(".env")
	if err != nil {
		return env
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		value = strings.Trim(value, "\"'")
		
		env[key] = value
	}
	
	return env
}

func updateEnvApiKey(apiKey string) error {
	env := loadEnvFile()
	env["CURSEFORGE_API_KEY"] = apiKey
	
	content := "# CurseForge API Key\nCURSEFORGE_API_KEY=" + apiKey + "\n"
	
	return os.WriteFile(".env", []byte(content), 0644)
}