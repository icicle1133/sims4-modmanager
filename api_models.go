package main

import (
	"time"
)

// Common models
type Pagination struct {
	Index        int `json:"index"`
	PageSize     int `json:"pageSize"`
	ResultCount  int `json:"resultCount"`
	TotalCount   int `json:"totalCount"`
}

// Game models
type Game struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	DateModified time.Time `json:"dateModified"`
	Assets      GameAssets `json:"assets"`
	Status      int       `json:"status"`
	ApiStatus   int       `json:"apiStatus"`
}

type GameAssets struct {
	IconUrl   string `json:"iconUrl"`
	TileUrl   string `json:"tileUrl"`
	CoverUrl  string `json:"coverUrl"`
}

type GameVersionType struct {
	ID         int    `json:"id"`
	GameID     int    `json:"gameId"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	IsSyncable bool   `json:"isSyncable"`
	Status     int    `json:"status"`
}

type GameVersionInfo struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
}

// Version types and versions
type VersionType struct {
	Type     int      `json:"type"`
	Versions []string `json:"versions"`
}

type VersionTypeV2 struct {
	Type     int                `json:"type"`
	Versions []GameVersionInfo  `json:"versions"`
}

// Mod models
type Mod struct {
	ID                  int       `json:"id"`
	GameID              int       `json:"gameId"`
	Name                string    `json:"name"`
	Slug                string    `json:"slug"`
	Links               ModLinks  `json:"links"`
	Summary             string    `json:"summary"`
	Status              int       `json:"status"`
	DownloadCount       int64     `json:"downloadCount"`
	IsFeatured          bool      `json:"isFeatured"`
	PrimaryCategoryID   int       `json:"primaryCategoryId"`
	Categories          []Category `json:"categories"`
	ClassID             int       `json:"classId"`
	Authors             []ModAuthor `json:"authors"`
	Logo                ModAsset  `json:"logo"`
	Screenshots         []ModAsset `json:"screenshots"`
	MainFileID          int       `json:"mainFileId"`
	LatestFiles         []File    `json:"latestFiles"`
	LatestFilesIndexes  []FileIndex `json:"latestFilesIndexes"`
	LatestEarlyAccessFilesIndexes []FileIndex `json:"latestEarlyAccessFilesIndexes,omitempty"`
	DateCreated         time.Time `json:"dateCreated"`
	DateModified        time.Time `json:"dateModified"`
	DateReleased        time.Time `json:"dateReleased"`
	AllowModDistribution bool     `json:"allowModDistribution"`
	GamePopularityRank  int       `json:"gamePopularityRank"`
	IsAvailable         bool      `json:"isAvailable"`
	ThumbsUpCount       int       `json:"thumbsUpCount"`
	Rating              float64   `json:"rating,omitempty"`
}

type ModLinks struct {
	WebsiteURL string `json:"websiteUrl"`
	WikiURL    string `json:"wikiUrl,omitempty"`
	IssuesURL  string `json:"issuesUrl,omitempty"`
	SourceURL  string `json:"sourceUrl,omitempty"`
}

type Category struct {
	ID               int       `json:"id"`
	GameID           int       `json:"gameId"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	URL              string    `json:"url"`
	IconURL          string    `json:"iconUrl"`
	DateModified     time.Time `json:"dateModified"`
	IsClass          bool      `json:"isClass"`
	ClassID          int       `json:"classId"`
	ParentCategoryID int       `json:"parentCategoryId,omitempty"`
	DisplayIndex     int       `json:"displayIndex"`
}

type ModAuthor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ModAsset struct {
	ID           int    `json:"id"`
	ModID        int    `json:"modId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnailUrl"`
	URL          string `json:"url"`
}

// File models
type FileHash struct {
	Value string `json:"value"`
	Algo  int    `json:"algo"`
}

type FileDependency struct {
	ModID        int `json:"modId"`
	RelationType int `json:"relationType"`
}

type FileModule struct {
	Name        string `json:"name"`
	Fingerprint int64  `json:"fingerprint"`
}

type File struct {
	ID                   int             `json:"id"`
	GameID               int             `json:"gameId"`
	ModID                int             `json:"modId"`
	IsAvailable          bool            `json:"isAvailable"`
	DisplayName          string          `json:"displayName"`
	FileName             string          `json:"fileName"`
	ReleaseType          int             `json:"releaseType"`
	FileStatus           int             `json:"fileStatus"`
	Hashes               []FileHash      `json:"hashes"`
	FileDate             time.Time       `json:"fileDate"`
	FileLength           int64           `json:"fileLength"`
	DownloadCount        int64           `json:"downloadCount"`
	FileSizeOnDisk       int64           `json:"fileSizeOnDisk,omitempty"`
	DownloadURL          string          `json:"downloadUrl,omitempty"`
	GameVersions         []string        `json:"gameVersions"`
	SortableGameVersions []GameVersion   `json:"sortableGameVersions"`
	Dependencies         []FileDependency `json:"dependencies,omitempty"`
	ExposeAsAlternative  bool            `json:"exposeAsAlternative,omitempty"`
	ParentProjectFileID  int             `json:"parentProjectFileId,omitempty"`
	AlternateFileID      int             `json:"alternateFileId,omitempty"`
	IsServerPack         bool            `json:"isServerPack,omitempty"`
	ServerPackFileID     int             `json:"serverPackFileId,omitempty"`
	IsEarlyAccessContent bool            `json:"isEarlyAccessContent,omitempty"`
	EarlyAccessEndDate   time.Time       `json:"earlyAccessEndDate,omitempty"`
	FileFingerprint      int64           `json:"fileFingerprint,omitempty"`
	Modules              []FileModule    `json:"modules,omitempty"`
}

type FileIndex struct {
	GameVersion       string `json:"gameVersion"`
	FileID            int    `json:"fileId"`
	Filename          string `json:"filename"`
	ReleaseType       int    `json:"releaseType"`
	GameVersionTypeID int    `json:"gameVersionTypeId"`
	ModLoader         int    `json:"modLoader"`
}

type GameVersion struct {
	GameVersionName        string    `json:"gameVersionName"`
	GameVersionPadded      string    `json:"gameVersionPadded"`
	GameVersion            string    `json:"gameVersion"`
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"`
	GameVersionTypeID      int       `json:"gameVersionTypeId"`
}

// Fingerprint models
type FingerprintMatch struct {
	ID          int     `json:"id"`
	File        File    `json:"file"`
	LatestFiles []File  `json:"latestFiles"`
	Fingerprints []uint `json:"fingerprints,omitempty"`
}

type FolderFingerprint struct {
	Foldername   string `json:"foldername"`
	Fingerprints []uint `json:"fingerprints"`
}

// Request models
type GetModsByIdsListRequestBody struct {
	ModIDs        []int  `json:"modIds"`
	FilterPcOnly  bool   `json:"filterPcOnly,omitempty"`
}

type GetModFilesRequestBody struct {
	FileIDs []int `json:"fileIds"`
}

type GetFeaturedModsRequestBody struct {
	GameID         int   `json:"gameId"`
	ExcludedModIDs []int `json:"excludedModIds"`
	GameVersionTypeID int `json:"gameVersionTypeId,omitempty"`
}

type GetFingerprintMatchesRequestBody struct {
	Fingerprints []uint `json:"fingerprints"`
}

type GetFuzzyMatchesRequestBody struct {
	GameID       int                 `json:"gameId"`
	Fingerprints []FolderFingerprint `json:"fingerprints"`
}

// Response models
type SearchModsResponse struct {
	Data       []Mod      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type GetModResponse struct {
	Data Mod `json:"data"`
}

type GetModFilesResponse struct {
	Data       []File     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type GetFeaturedModsResponse struct {
	Data struct {
		Featured        []Mod `json:"featured"`
		Popular         []Mod `json:"popular"`
		RecentlyUpdated []Mod `json:"recentlyUpdated"`
	} `json:"data"`
}

type GetModsResponse struct {
	Data []Mod `json:"data"`
}

type GetFilesResponse struct {
	Data []File `json:"data"`
}

type StringResponse struct {
	Data string `json:"data"`
}

type GetGamesResponse struct {
	Data       []Game     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type GetGameResponse struct {
	Data Game `json:"data"`
}

type GetVersionsResponse struct {
	Data []VersionType `json:"data"`
}

type GetVersionsV2Response struct {
	Data []VersionTypeV2 `json:"data"`
}

type GetVersionTypesResponse struct {
	Data []GameVersionType `json:"data"`
}

type GetCategoriesResponse struct {
	Data []Category `json:"data"`
}

type FingerprintMatchesResponse struct {
	Data struct {
		IsCacheBuilt              bool                   `json:"isCacheBuilt"`
		ExactMatches              []FingerprintMatch     `json:"exactMatches"`
		ExactFingerprints         []uint                 `json:"exactFingerprints"`
		PartialMatches            []FingerprintMatch     `json:"partialMatches"`
		PartialMatchFingerprints  map[string][]uint      `json:"partialMatchFingerprints"`
		InstalledFingerprints     []uint                 `json:"installedFingerprints"`
		UnmatchedFingerprints     []uint                 `json:"unmatchedFingerprints"`
	} `json:"data"`
}

type FingerprintFuzzyMatchesResponse struct {
	Data struct {
		FuzzyMatches []FingerprintMatch `json:"fuzzyMatches"`
	} `json:"data"`
}

// Enum constants
const (
	// ModLoaderType
	ModLoaderTypeAny      = 0
	ModLoaderTypeFORGE    = 1
	ModLoaderTypeCauldron = 2
	ModLoaderTypeLiteLoader = 3
	ModLoaderTypeFabric   = 4
	ModLoaderTypeQuilt    = 5
	ModLoaderTypeNeoForge = 6
	
	// ModsSearchSortField
	SortFieldFeatured = 1
	SortFieldPopularity = 2
	SortFieldLastUpdated = 3
	SortFieldName = 4
	SortFieldAuthor = 5
	SortFieldTotalDownloads = 6
	SortFieldCategory = 7
	SortFieldGameVersion = 8
	SortFieldDailyDownloads = 9
	SortFieldRating = 10
	SortFieldCreatedAt = 11
	SortFieldReleaseDate = 12
	
	// SortOrder
	SortOrderAscending = "asc"
	SortOrderDescending = "desc"
	
	// FileHashAlgo
	HashAlgoSha1 = 1
	HashAlgoMd5 = 2
	
	// FileReleaseType
	ReleaseTypeRelease = 1
	ReleaseTypeBeta = 2
	ReleaseTypeAlpha = 3
	
	// FileStatus
	FileStatusProcessing = 1
	FileStatusChangesRequired = 2
	FileStatusUnderReview = 3
	FileStatusApproved = 4
	FileStatusRejected = 5
	FileStatusMalware = 6
	FileStatusDeleted = 7
	FileStatusArchived = 8
	FileStatusTesting = 9
	FileStatusReleased = 10
	FileStatusReadyForReview = 11
	FileStatusDeprecated = 12
	FileStatusBaking = 13
	FileStatusAwaitingPublishing = 14
	FileStatusFailed = 15
	
	// FileDependencyRelationType
	RelationTypeEmbedded = 1
	RelationTypeOptional = 2
	RelationTypeRequired = 3
	RelationTypeToolRequired = 4
	RelationTypeIncompatible = 5
	RelationTypeInclude = 6
)