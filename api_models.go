package main

import (
	"time"
)

type Pagination struct {
	Index        int `json:"index"`
	PageSize     int `json:"pageSize"`
	ResultCount  int `json:"resultCount"`
	TotalCount   int `json:"totalCount"`
}

type Mod struct {
	ID                int       `json:"id"`
	GameID            int       `json:"gameId"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	Links             ModLinks  `json:"links"`
	Summary           string    `json:"summary"`
	Status            int       `json:"status"`
	DownloadCount     int64     `json:"downloadCount"`
	IsFeatured        bool      `json:"isFeatured"`
	PrimaryCategoryID int       `json:"primaryCategoryId"`
	Categories        []Category `json:"categories"`
	Authors           []ModAuthor `json:"authors"`
	Logo              ModAsset  `json:"logo"`
	Screenshots       []ModAsset `json:"screenshots"`
	MainFileID        int       `json:"mainFileId"`
	LatestFiles       []File    `json:"latestFiles"`
	DateCreated       time.Time `json:"dateCreated"`
	DateModified      time.Time `json:"dateModified"`
	DateReleased      time.Time `json:"dateReleased"`
	AllowModDistribution bool    `json:"allowModDistribution"`
	GamePopularityRank    int    `json:"gamePopularityRank"`
	IsAvailable           bool   `json:"isAvailable"`
	Thumbs                int    `json:"thumbs"`
}

type ModLinks struct {
	WebsiteURL string `json:"websiteUrl"`
	WikiURL    string `json:"wikiUrl"`
	IssuesURL  string `json:"issuesUrl"`
	SourceURL  string `json:"sourceUrl"`
}

type Category struct {
	ID       int    `json:"id"`
	GameID   int    `json:"gameId"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	URL      string `json:"url"`
	IconURL  string `json:"iconUrl"`
	DateModified time.Time `json:"dateModified"`
	IsClass       bool     `json:"isClass"`
	ClassID       int      `json:"classId"`
	ParentCategoryID int   `json:"parentCategoryId"`
}

type ModAuthor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ModAsset struct {
	ID       int    `json:"id"`
	ModID    int    `json:"modId"`
	Title    string `json:"title"`
	Description string `json:"description"`
	ThumbnailURL string `json:"thumbnailUrl"`
	URL      string `json:"url"`
}

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
	ID                   int           `json:"id"`
	GameID               int           `json:"gameId"`
	ModID                int           `json:"modId"`
	IsAvailable          bool          `json:"isAvailable"`
	DisplayName          string        `json:"displayName"`
	FileName             string        `json:"fileName"`
	ReleaseType          int           `json:"releaseType"`
	FileStatus           int           `json:"fileStatus"`
	Hashes               []FileHash    `json:"hashes"`
	FileDate             time.Time     `json:"fileDate"`
	FileLength           int64         `json:"fileLength"`
	DownloadCount        int64         `json:"downloadCount"`
	FileSizeOnDisk       int64         `json:"fileSizeOnDisk,omitempty"`
	DownloadURL          string        `json:"downloadUrl,omitempty"`
	GameVersions         []string      `json:"gameVersions"`
	SortableGameVersions []GameVersion `json:"sortableGameVersions"`
	Dependencies         []FileDependency `json:"dependencies,omitempty"`
	ExposeAsAlternative  bool          `json:"exposeAsAlternative,omitempty"`
	ParentProjectFileID  int           `json:"parentProjectFileId,omitempty"`
	AlternateFileID      int           `json:"alternateFileId,omitempty"`
	IsServerPack         bool          `json:"isServerPack,omitempty"`
	ServerPackFileID     int           `json:"serverPackFileId,omitempty"`
	IsEarlyAccessContent bool          `json:"isEarlyAccessContent,omitempty"`
	EarlyAccessEndDate   time.Time     `json:"earlyAccessEndDate,omitempty"`
	FileFingerprint      int64         `json:"fileFingerprint,omitempty"`
	Modules              []FileModule  `json:"modules,omitempty"`
}

type GameVersion struct {
	GameVersionName      string    `json:"gameVersionName"`
	GameVersionPadded    string    `json:"gameVersionPadded"`
	GameVersion          string    `json:"gameVersion"`
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"`
	GameVersionTypeID    int       `json:"gameVersionTypeId"`
}

type SearchModsResponse struct {
	Data      []Mod      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type GetModResponse struct {
	Data Mod `json:"data"`
}

type GetModFilesResponse struct {
	Data      []File     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type GetFeaturedModsResponse struct {
	Data struct {
		Featured []Mod `json:"featured"`
		Popular  []Mod `json:"popular"`
		RecentlyUpdated []Mod `json:"recentlyUpdated"`
	} `json:"data"`
}

type StringResponse struct {
	Data string `json:"data"`
}