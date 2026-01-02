package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func compressAndSave(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()
	
	_, err = gzWriter.Write(data)
	return err
}

func loadAndDecompress(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	
	return io.ReadAll(gzReader)
}

func saveCompressedJson(v interface{}, filename string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	
	return compressAndSave(data, filename)
}

func loadCompressedJson(v interface{}, filename string) error {
	data, err := loadAndDecompress(filename)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, v)
}

func ensureDirectoryExists(path string) error {
	return os.MkdirAll(path, 0755)
}

func ensureModsDirectory(path string) error {
	err := ensureDirectoryExists(path)
	if err != nil {
		return err
	}
	
	readmePath := filepath.Join(path, "README.txt")
	
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		content := "This directory is managed by Sims 4 Mod Manager.\n" +
			"Please do not modify the files in this directory manually.\n"
		return os.WriteFile(readmePath, []byte(content), 0644)
	}
	
	return nil
}

func loadRecentMods() ([]ModInfo, error) {
	var mods []ModInfo
	
	filename := "recent_mods.dat"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return mods, nil
	}
	
	err := loadCompressedJson(&mods, filename)
	return mods, err
}