package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var apiClient *ApiClient
var currentPage = 1
var lastSearch = ""

func setupBrowserTab() fyne.CanvasObject {
	if apiClient == nil {
		settings, _ := LoadSettings()
		if settings.ApiKey == "" {
			return setupApiKeyPrompt()
		} else {
			apiClient = NewApiClient(settings.ApiKey)
			// Verify API key works by making a simple request
			_, err := apiClient.GetGames(0, 1)
			if err != nil {
				fmt.Printf("API key validation failed: %v\n", err)
				return setupApiKeyPrompt()
			}
		}
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search mods...")
	
	contentContainer := container.NewVBox()
	
	searchButton := widget.NewButton("Search", func() {
		currentPage = 1
		lastSearch = searchEntry.Text
		refreshModBrowser(lastSearch, currentPage, contentContainer)
	})

	nextButton := widget.NewButton("Next Page", func() {
		currentPage++
		refreshModBrowser(lastSearch, currentPage, contentContainer)
	})
	
	prevButton := widget.NewButton("Previous Page", func() {
		if currentPage > 1 {
			currentPage--
			refreshModBrowser(lastSearch, currentPage, contentContainer)
		}
	})
	
	pageControls := container.NewHBox(
		prevButton,
		widget.NewLabel("Page:"),
		widget.NewLabel(strconv.Itoa(currentPage)),
		nextButton,
	)
	
	searchRow := container.NewBorder(nil, nil, nil, searchButton, searchEntry)
	
	loadFeaturedMods(contentContainer)
	
	return container.NewBorder(
		container.NewVBox(widget.NewLabel("Browse Mods"), searchRow),
		pageControls,
		nil, nil,
		container.NewVScroll(contentContainer),
	)
}

func setupApiKeyPrompt() fyne.CanvasObject {
	keyEntry := widget.NewPasswordEntry()
	
	env := loadEnvFile()
	if apiKey, ok := env["CURSEFORGE_API_KEY"]; ok && apiKey != "" {
		keyEntry.SetText(apiKey)
	}
	
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "CurseForge API Key", Widget: keyEntry},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			if keyEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("API key cannot be empty"), fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}
			
			// Create a temporary client to validate the API key
			tempClient := NewApiClient(keyEntry.Text)
			_, err := tempClient.GetGames(0, 1)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Invalid API key: %v", err), fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}
			
			settings, _ := LoadSettings()
			settings.ApiKey = keyEntry.Text
			SaveSettings(settings)
			
			updateEnvApiKey(keyEntry.Text)
			
			apiClient = tempClient
			
			tabs := fyne.CurrentApp().Driver().AllWindows()[0].Content().(*container.AppTabs)
			browserTab := container.NewTabItem("Browse", setupBrowserTab())
			
			for i, tab := range tabs.Items {
				if tab.Text == "Browse" {
					tabs.Items[i] = browserTab
					tabs.Refresh()
					tabs.Select(browserTab)
					return
				}
			}
			
			tabs.Append(browserTab)
			tabs.Select(browserTab)
		},
	}
	
	return container.NewVBox(
		widget.NewLabel("Enter your CurseForge API Key to browse mods"),
		widget.NewLabel("The key will be saved in the .env file"),
		form,
	)
}

func loadFeaturedMods(container *fyne.Container) {
	container.RemoveAll()
	loadingLabel := widget.NewLabel("Loading featured mods...")
	container.Add(loadingLabel)
	container.Refresh()
	
	go func() {
		featured, err := apiClient.GetFeaturedMods()
		if err != nil {
			container.RemoveAll()
			container.Add(widget.NewLabel("Error loading featured mods: " + err.Error()))
			container.Refresh()
			return
		}
		
		container.RemoveAll()
		
		container.Add(widget.NewLabel("Popular Mods"))
		for _, mod := range featured.Data.Popular {
			container.Add(createModCard(mod))
		}
		
		container.Add(widget.NewSeparator())
		container.Add(widget.NewLabel("Featured Mods"))
		for _, mod := range featured.Data.Featured {
			container.Add(createModCard(mod))
		}
		
		container.Add(widget.NewSeparator())
		container.Add(widget.NewLabel("Recently Updated Mods"))
		for _, mod := range featured.Data.RecentlyUpdated {
			container.Add(createModCard(mod))
		}
		
		container.Refresh()
	}()
}

func refreshModBrowser(search string, page int, container *fyne.Container) {
	container.RemoveAll()
	loadingLabel := widget.NewLabel("Searching for mods...")
	container.Add(loadingLabel)
	container.Refresh()
	
	go func() {
		var searchResults SearchModsResponse
		var err error
		
		if search == "" {
			loadFeaturedMods(container)
			return
		} else {
			searchResults, err = apiClient.SearchMods(search, page)
			if err != nil {
				container.RemoveAll()
				container.Add(widget.NewLabel("Error searching mods: " + err.Error()))
				container.Refresh()
				return
			}
		}
		
		container.RemoveAll()
		
		if len(searchResults.Data) == 0 {
			container.Add(widget.NewLabel("No mods found matching your search."))
			container.Refresh()
			return
		}
		
		resultLabel := widget.NewLabel(fmt.Sprintf("Showing %d of %d results", len(searchResults.Data), searchResults.Pagination.TotalCount))
		container.Add(resultLabel)
		
		for _, mod := range searchResults.Data {
			container.Add(createModCard(mod))
		}
		
		container.Refresh()
	}()
}

func createModCard(mod Mod) fyne.CanvasObject {
	nameLabel := widget.NewLabel(mod.Name)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	authorText := "By: "
	for i, author := range mod.Authors {
		if i > 0 {
			authorText += ", "
		}
		authorText += author.Name
	}
	authorLabel := widget.NewLabel(authorText)
	
	downloadsLabel := widget.NewLabel(fmt.Sprintf("Downloads: %d", mod.DownloadCount))
	
	summaryLabel := widget.NewLabel(mod.Summary)
	summaryLabel.Wrapping = fyne.TextWrapWord
	
	installButton := widget.NewButton("Install", func() {
		showModFiles(mod)
	})
	
	detailsButton := widget.NewButton("Details", func() {
		showModDetails(mod)
	})
	
	buttonRow := container.NewHBox(layout.NewSpacer(), detailsButton, installButton)
	
	card := container.NewVBox(
		nameLabel,
		authorLabel,
		downloadsLabel,
		summaryLabel,
		buttonRow,
	)
	
	return container.NewBorder(nil, nil, nil, nil, card)
}

func showModDetails(mod Mod) {
	detailsWindow := fyne.CurrentApp().NewWindow("Mod Details")
	detailsWindow.Resize(fyne.NewSize(800, 600))
	
	nameLabel := widget.NewLabel(mod.Name)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	authorText := "By: "
	for i, author := range mod.Authors {
		if i > 0 {
			authorText += ", "
		}
		authorText += author.Name
	}
	authorLabel := widget.NewLabel(authorText)
	
	downloadsLabel := widget.NewLabel(fmt.Sprintf("Downloads: %d", mod.DownloadCount))
	dateLabel := widget.NewLabel(fmt.Sprintf("Updated: %s", mod.DateModified.Format("2006-01-02")))
	
	installButton := widget.NewButton("Install", func() {
		detailsWindow.Close()
		showModFiles(mod)
	})
	
	headerBox := container.NewVBox(
		nameLabel,
		authorLabel,
		downloadsLabel,
		dateLabel,
	)
	
	descriptionLabel := widget.NewLabel("Loading description...")
	descriptionLabel.Wrapping = fyne.TextWrapWord
	
	categoriesText := "Categories: "
	for i, cat := range mod.Categories {
		if i > 0 {
			categoriesText += ", "
		}
		categoriesText += cat.Name
	}
	categoriesLabel := widget.NewLabel(categoriesText)
	
	websiteButton := widget.NewButton("Visit Website", func() {
		// like i said in other files, this is not implemented. ignore.
	})
	
	headerBox.Add(categoriesLabel)
	
	content := container.NewBorder(
		headerBox,
		container.NewHBox(layout.NewSpacer(), websiteButton, installButton),
		nil, nil,
		container.NewScroll(descriptionLabel),
	)
	
	detailsWindow.SetContent(content)
	
	go func() {
		descResp, err := apiClient.GetModDescription(mod.ID)
		if err != nil {
			descriptionLabel.SetText("Failed to load description: " + err.Error())
			return
		}
		
		plainText := StripHTML(descResp.Data)
		descriptionLabel.SetText(plainText)
		detailsWindow.Content().Refresh()
	}()
	
	detailsWindow.Show()
	detailsWindow.CenterOnScreen()
}

func showModFiles(mod Mod) {
	filesWindow := fyne.CurrentApp().NewWindow("Mod Files")
	filesWindow.Resize(fyne.NewSize(600, 400))
	
	loadingLabel := widget.NewLabel("Loading files...")
	filesWindow.SetContent(loadingLabel)
	filesWindow.Show()
	
	go func() {
		filesResp, err := apiClient.GetModFiles(mod.ID)
		if err != nil {
			filesWindow.SetContent(widget.NewLabel("Error loading files: " + err.Error()))
			return
		}
		
		if len(filesResp.Data) == 0 {
			filesWindow.SetContent(widget.NewLabel("No files available for this mod"))
			return
		}
		
		filesList := widget.NewList(
			func() int { return len(filesResp.Data) },
			func() fyne.CanvasObject {
				return container.NewHBox(
					widget.NewLabel("Filename"),
					widget.NewLabel("Version"),
					widget.NewButton("Download", func() {}),
				)
			},
			func(id widget.ListItemID, item fyne.CanvasObject) {
				file := filesResp.Data[id]
				
				container := item.(*fyne.Container)
				
				nameLabel := container.Objects[0].(*widget.Label)
				nameLabel.SetText(file.FileName)
				
				versionLabel := container.Objects[1].(*widget.Label)
				if len(file.GameVersions) > 0 {
					versionLabel.SetText("Game: " + file.GameVersions[0])
				} else {
					versionLabel.SetText("Unknown version")
				}
				
				downloadButton := container.Objects[2].(*widget.Button)
				downloadButton.OnTapped = func() {
					filesWindow.Close() // i may or may not have forgot to add this when i first did this
					downloadFile(mod, file)
				}
			},
		)
		
		filesWindow.SetContent(container.NewBorder(
			widget.NewLabel(fmt.Sprintf("Files for %s", mod.Name)),
			nil, nil, nil,
			container.NewScroll(filesList),
		))
	}()
}

func downloadFile(mod Mod, file File) {
	settings, err := LoadSettings()
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	filename := file.FileName
	progressMessage := "Preparing download for " + filename
	progress := dialog.NewProgress("Downloading", progressMessage, fyne.CurrentApp().Driver().AllWindows()[0])
	progress.Show()
	
	go func() {
		downloadURL := file.DownloadURL
		
		if downloadURL == "" {
			urlResp, err := apiClient.GetModFileDownloadURL(mod.ID, file.ID)
			if err == nil && urlResp.Data != "" {
				downloadURL = urlResp.Data
				fmt.Printf("Got a fucking download URL: %s\n", downloadURL)
			}
		}
		
		if downloadURL == "" {
			fingerprints := []uint{}
			
			if file.FileFingerprint > 0 {
				fingerprints = append(fingerprints, uint(file.FileFingerprint))
			} else {
				for i := range file.Modules {
					if file.Modules[i].Fingerprint > 0 {
						fingerprints = append(fingerprints, uint(file.Modules[i].Fingerprint))
					}
				}
			}
			
			if len(fingerprints) == 0 {
				fingerprints = append(fingerprints, uint(file.ID))
			}
			
			fmt.Printf("Using fingerprints: %v\n", fingerprints)
			
			fingerprintResp, err := apiClient.MatchFingerprints(fingerprints)
			if err == nil {
				fmt.Printf("Found matches: Exact=%d, Partial=%d\n", 
					len(fingerprintResp.Data.ExactMatches), len(fingerprintResp.Data.PartialMatches))
				
				for i := range fingerprintResp.Data.ExactMatches {
					if fingerprintResp.Data.ExactMatches[i].File.DownloadURL != "" {
						downloadURL = fingerprintResp.Data.ExactMatches[i].File.DownloadURL
						break
					}
				}
				
				if downloadURL == "" && len(fingerprintResp.Data.PartialMatches) > 0 {
					for i := range fingerprintResp.Data.PartialMatches {
						if fingerprintResp.Data.PartialMatches[i].File.DownloadURL != "" {
							downloadURL = fingerprintResp.Data.PartialMatches[i].File.DownloadURL
							break
						}
					}
				}
			}
		}
		
		if downloadURL == "" {
			fmt.Printf("Fuck, no download URL. Let's try to make one...\n")
			fileID := file.ID
			thousands := fileID / 1000
			remainder := fileID % 1000
			
			downloadURL = fmt.Sprintf("https://mediafilez.forgecdn.net/files/%d/%d/%s", 
				thousands, remainder, file.FileName)
			fmt.Printf("Made a URL: %s\n", downloadURL)
		}
		
		if downloadURL == "" {
			progress.Hide()
			
			websiteURL := ""
			if mod.Links.WebsiteURL != "" {
				websiteURL = mod.Links.WebsiteURL
			} else {
				websiteURL = fmt.Sprintf("https://www.curseforge.com/sims4/mods/%s/files/%d", mod.Slug, file.ID)
			}
			
			confirmDialog := dialog.NewConfirm(
				"Shit! Can't Download",
				fmt.Sprintf("I can't download '%s' directly. Wanna go to the website and do it yourself?", file.FileName),
				func(ok bool) {
					if ok {
						fmt.Printf("Would open browser to: %s\n", websiteURL)
						
						dialog.ShowInformation(
							"Manual Download",
							fmt.Sprintf("Download it and put it here:\n%s", settings.ModsDirectory),
							fyne.CurrentApp().Driver().AllWindows()[0],
						)
					}
				},
				fyne.CurrentApp().Driver().AllWindows()[0],
			)
			confirmDialog.Show()
			return
		}
		
		fmt.Printf("Download URL: %s\n", downloadURL)
		
		err = ensureDirectoryExists(settings.ModsDirectory)
		if err != nil {
			progress.Hide()
			dialog.ShowError(fmt.Errorf("failed to create mods directory: %w", err), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		
		client := &http.Client{
			Timeout: 60 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects, what the fuck")
				}
				return nil
			},
		}

		req, err := http.NewRequest("GET", downloadURL, nil) // let's try to download this shit
		if err != nil {
			progress.Hide()
			dialog.ShowError(fmt.Errorf("download request failed: %w", err), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}

		req.Header.Add("User-Agent", "Sims4ModManager/1.0") // fake being a real browser
		req.Header.Add("Accept", "*/*") // take any content type, we're desperate

		// Execute the request
		resp, err := client.Do(req)
		if err != nil {
			progress.Hide()
			dialog.ShowError(fmt.Errorf("download failed: %w", err), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		if resp == nil || resp.StatusCode != http.StatusOK {
			progress.Hide()
			dialog.ShowError(fmt.Errorf("download failed: server returned status %d", resp.StatusCode), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		defer resp.Body.Close()
		
		targetPath := filepath.Join(settings.ModsDirectory, file.FileName)
		
		if _, err := os.Stat(targetPath); err == nil {
			progress.Hide()
			confirmDialog := dialog.NewConfirm(
				"File Already Exists",
				"A file with the name "+file.FileName+" already exists. Do you want to overwrite it?",
				func(confirmed bool) {
					if confirmed {
						downloadToFile(resp, targetPath, file.FileLength, progress)
					}
				},
				fyne.CurrentApp().Driver().AllWindows()[0],
			)
			confirmDialog.Show()
		} else {
			downloadToFile(resp, targetPath, file.FileLength, progress)
		}
	}()
	
	progress.Show()
}

func downloadToFile(resp *http.Response, targetPath string, fileSize int64, progress *dialog.ProgressDialog) {
	out, err := os.Create(targetPath)
	if err != nil {
		progress.Hide()
		dialog.ShowError(fmt.Errorf("can't create the damn file: %w", err), fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	defer out.Close()
	
	buffer := make([]byte, 4096)
	var downloaded int64
	
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := out.Write(buffer[:n])
			if writeErr != nil {
				progress.Hide()
				dialog.ShowError(fmt.Errorf("fuck, can't write to file: %w", writeErr), fyne.CurrentApp().Driver().AllWindows()[0])
				return
			}
			
			downloaded += int64(n)
			if fileSize > 0 {
				progress.SetValue(float64(downloaded) / float64(fileSize))
			} else {
				progress.SetValue(0.5) // who the hell knows how big this file is
			}
		}
		
		if err != nil {
			if err == io.EOF {
				break
			}
			progress.Hide()
			dialog.ShowError(fmt.Errorf("shit broke during download: %w", err), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
	}
	
	progress.Hide()
	
	successDialog := dialog.NewInformation(
		"Got it!",
		"Mod installed. Enjoy your game!",
		fyne.CurrentApp().Driver().AllWindows()[0],
	)
	successDialog.Show()
}