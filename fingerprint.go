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

type FingerprintMatch struct {
	ID         int   `json:"id"`
	File       File  `json:"file"`
	LatestFiles []File `json:"latestFiles"`
}

type FingerprintMatchesResponse struct {
	Data struct {
		IsCacheBuilt            bool               `json:"isCacheBuilt"`
		ExactMatches            []FingerprintMatch `json:"exactMatches"`
		ExactFingerprints       []uint             `json:"exactFingerprints"`
		PartialMatches          []FingerprintMatch `json:"partialMatches"`
		PartialMatchFingerprints map[string][]uint  `json:"partialMatchFingerprints"`
		InstalledFingerprints   []uint             `json:"installedFingerprints"`
		UnmatchedFingerprints   []uint             `json:"unmatchedFingerprints"`
	} `json:"data"`
}

func CalculateFingerprint(path string) (uint, error) {
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

func CalculateFingerprintsForDir(dir string) ([]uint, error) {
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

func (c *ApiClient) MatchFingerprints(fingerprints []uint) (FingerprintMatchesResponse, error) {
	var result FingerprintMatchesResponse
	
	requestBody := map[string]interface{}{
		"fingerprints": fingerprints,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return result, err
	}
	
	endpoint := fmt.Sprintf("/v1/fingerprints/78062")
	
	responseBody, err := c.makeRequest("POST", endpoint, jsonBody)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}