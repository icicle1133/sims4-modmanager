package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL     = "https://api.curseforge.com"
	sims4GameID = 78062
)

var apiKey string

type ApiClient struct {
	client *http.Client
	apiKey string
}

func NewApiClient(key string) *ApiClient {
	apiKey = key
	return &ApiClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: key,
	}
}

func (c *ApiClient) makeRequest(method, endpoint string, body []byte) ([]byte, error) {
	url := baseURL + endpoint
	
	fmt.Printf("Making %s request to: %s\n", method, url)
	
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err) // fuck this error handling
	}
	
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	
	if method == "POST" && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	fmt.Printf("Request headers:\n")
	for k, v := range req.Header {
		if k == "x-api-key" {
			fmt.Printf("  %s: %s...\n", k, v[0][:10])
		} else {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err) // shit, network error
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err) // can't even read the body wtf
	}
	
	fmt.Printf("Response status: %d %s\n", resp.StatusCode, resp.Status)
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error response body: %s\n", string(responseBody))
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode) // fucking API errors
	}
	
	previewLen := 100
	if len(responseBody) < previewLen {
		previewLen = len(responseBody)
	}
	fmt.Printf("Response preview: %s...\n", string(responseBody[:previewLen]))
	
	return responseBody, nil
}

func (c *ApiClient) SearchMods(searchFilter string, page int) (SearchModsResponse, error) {
	var result SearchModsResponse
	
	params := url.Values{}
	params.Add("gameId", strconv.Itoa(sims4GameID))
	if searchFilter != "" {
		params.Add("searchFilter", searchFilter)
	}
	params.Add("pageSize", "20")
	params.Add("index", strconv.Itoa((page-1)*20))
	params.Add("sortField", "2") // 2 = Popularity, this API is so damn unintuitive
	params.Add("sortOrder", "desc")
	
	endpoint := "/v1/mods/search?" + params.Encode()
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetMod(modId int) (GetModResponse, error) {
	var result GetModResponse
	
	endpoint := fmt.Sprintf("/v1/mods/%d", modId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetModDescription(modId int) (StringResponse, error) {
	var result StringResponse
	
	endpoint := fmt.Sprintf("/v1/mods/%d/description", modId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetFeaturedMods() (GetFeaturedModsResponse, error) {
	var result GetFeaturedModsResponse
	
	requestBody := map[string]interface{}{
		"gameId": sims4GameID,
		"excludedModIds": []int{},
	}
	
	body, err := json.Marshal(requestBody)
	if err != nil {
		return result, err
	}
	
	endpoint := "/v1/mods/featured"
	
	responseBody, err := c.makeRequest("POST", endpoint, body)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetModFiles(modId int) (GetModFilesResponse, error) {
	var result GetModFilesResponse
	
	endpoint := fmt.Sprintf("/v1/mods/%d/files", modId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	
	for i, file := range result.Data {
		fmt.Printf("File %d: ID=%d, Name=%s, DownloadURL=%s\n", 
			i, file.ID, file.FileName, file.DownloadURL)
	}
	
	return result, err
}

func (c *ApiClient) GetModFileDownloadURL(modId, fileId int) (StringResponse, error) {
	var result StringResponse
	
	endpoint := fmt.Sprintf("/v1/mods/%d/files/%d/download-url", modId, fileId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetModFileChangelog(modId, fileId int) (StringResponse, error) {
	var result StringResponse
	
	endpoint := fmt.Sprintf("/v1/mods/%d/files/%d/changelog", modId, fileId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetGames(index, pageSize int) (GetGamesResponse, error) {
	var result GetGamesResponse
	
	params := url.Values{}
	if index > 0 {
		params.Add("index", strconv.Itoa(index))
	}
	if pageSize > 0 && pageSize <= 50 {
		params.Add("pageSize", strconv.Itoa(pageSize))
	}
	
	endpoint := "/v1/games"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetGame(gameId int) (GetGameResponse, error) {
	var result GetGameResponse
	
	endpoint := fmt.Sprintf("/v1/games/%d", gameId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetGameVersions(gameId int) (GetVersionsResponse, error) {
	var result GetVersionsResponse
	
	endpoint := fmt.Sprintf("/v1/games/%d/versions", gameId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetGameVersionsV2(gameId int) (GetVersionsV2Response, error) {
	var result GetVersionsV2Response
	
	endpoint := fmt.Sprintf("/v2/games/%d/versions", gameId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetGameVersionTypes(gameId int) (GetVersionTypesResponse, error) {
	var result GetVersionTypesResponse
	
	endpoint := fmt.Sprintf("/v1/games/%d/version-types", gameId)
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetCategories(gameId int, classId int, classesOnly bool) (GetCategoriesResponse, error) {
	var result GetCategoriesResponse
	
	params := url.Values{}
	params.Add("gameId", strconv.Itoa(gameId))
	
	if classId > 0 {
		params.Add("classId", strconv.Itoa(classId))
	}
	
	if classesOnly {
		params.Add("classesOnly", "true")
	}
	
	endpoint := "/v1/categories?" + params.Encode()
	
	responseBody, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
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
	
	endpoint := fmt.Sprintf("/v1/fingerprints/%d", sims4GameID)
	
	responseBody, err := c.makeRequest("POST", endpoint, jsonBody)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) MatchFingerprintsGeneric(fingerprints []uint) (FingerprintMatchesResponse, error) {
	var result FingerprintMatchesResponse
	
	requestBody := map[string]interface{}{
		"fingerprints": fingerprints,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return result, err
	}
	
	endpoint := "/v1/fingerprints"
	
	responseBody, err := c.makeRequest("POST", endpoint, jsonBody)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetModsByIds(modIds []int) (GetModsResponse, error) {
	var result GetModsResponse
	
	requestBody := map[string]interface{}{
		"modIds": modIds,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return result, err
	}
	
	endpoint := "/v1/mods"
	
	responseBody, err := c.makeRequest("POST", endpoint, jsonBody)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}

func (c *ApiClient) GetFilesByIds(fileIds []int) (GetFilesResponse, error) {
	var result GetFilesResponse
	
	requestBody := map[string]interface{}{
		"fileIds": fileIds,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return result, err
	}
	
	endpoint := "/v1/mods/files"
	
	responseBody, err := c.makeRequest("POST", endpoint, jsonBody)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(responseBody, &result)
	return result, err
}