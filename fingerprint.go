package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spaolacci/murmur3"
)


func CalculateFingerprint(path string) (uint, error) { // this murmur3 shit better work
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	
	h := murmur3.New128()
	
	if _, err := io.Copy(h, file); err != nil {
		return 0, err
	}
	
	s := h.Sum(nil)
	return uint(binary.LittleEndian.Uint64(s)), nil
}


func CalculateFingerprintsForDir(dir string) ([]uint, error) { // fuck me this is gonna be slow
	var fingerprints []uint
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		ext := filepath.Ext(path)
		if !info.IsDir() && (ext == ".package" || ext == ".ts4script") {
			fp, err := CalculateFingerprint(path)
			if err != nil {
				return fmt.Errorf("error calculating fingerprint for %s: %v", path, err)
			}
			
			fingerprints = append(fingerprints, fp)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return fingerprints, nil
}


func CalculateFuzzyFingerprintsForDir(dir string) ([]FolderFingerprint, error) { // why the fuck do we need fuzzy matching
	folders := make(map[string][]uint)
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		ext := filepath.Ext(path)
		if !info.IsDir() && (ext == ".package" || ext == ".ts4script") {
			folderName := filepath.Base(filepath.Dir(path))
			
			fp, err := CalculateFingerprint(path)
			if err != nil {
				return fmt.Errorf("error calculating fingerprint for %s: %v", path, err)
			}
			
			folders[folderName] = append(folders[folderName], fp)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	var result []FolderFingerprint
	for folder, prints := range folders {
		result = append(result, FolderFingerprint{
			Foldername:   folder,
			Fingerprints: prints,
		})
	}
	
	return result, nil
}


func (c *ApiClient) MatchFuzzyFingerprints(fingerprints []FolderFingerprint) (FingerprintFuzzyMatchesResponse, error) { // this API is so damn inconsistent
	var result FingerprintFuzzyMatchesResponse
	
	requestBody := map[string]interface{}{
		"gameId":       sims4GameID,
		"fingerprints": fingerprints,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return result, err
	}
	
	endpoint := fmt.Sprintf("/v1/fingerprints/fuzzy/%d", sims4GameID)
	
	responseBody, err := c.makeRequest("POST", endpoint, jsonBody)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}