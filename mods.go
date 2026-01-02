package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func setupModsTab() fyne.CanvasObject {
	modsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return container.NewBorder(
				nil, nil, nil, widget.NewButton("Remove", func() {}),
				container.NewVBox(
					widget.NewLabel("Mod Name"),
					container.NewHBox(widget.NewIcon(theme.InfoIcon()), widget.NewLabel("Install Date")),
					widget.NewLabel("Size"),
				),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	refreshButton := widget.NewButton("Refresh Mods", func() {
		refreshModsList(modsList)
	})
	
	installButton := widget.NewButton("Browse Mods", func() {
		showModBrowser(modsList)
	})
	
	return container.NewBorder(
		widget.NewLabel("Installed Mods"),
		container.NewHBox(refreshButton, installButton),
		nil, nil, container.NewVScroll(modsList),
	)
}

func refreshModsList(list *widget.List) {
	settings, err := LoadSettings()
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	mods, err := scanMods(settings.ModsDirectory)
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}

	list.Length = func() int { return len(mods) }
	list.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		mod := mods[id]
		container := item.(*fyne.Container)
		
		innerContainer := container.Objects[0].(*fyne.Container)
		
		nameLabel := innerContainer.Objects[0].(*widget.Label)
		nameLabel.SetText(mod.Name)
		
		dateContainer := innerContainer.Objects[1].(*fyne.Container)
		dateLabel := dateContainer.Objects[1].(*widget.Label)
		dateLabel.SetText(mod.InstallDate.Format("2006-01-02 15:04:05"))
		
		sizeLabel := innerContainer.Objects[2].(*widget.Label)
		sizeLabel.SetText(formatFileSize(mod.FileSize))
		
		removeButton := container.Objects[1].(*widget.Button)
		removeButton.OnTapped = func() {
			removeMod(mod, list)
		}
	}
	
	list.Refresh()
}

func scanMods(directory string) ([]ModInfo, error) {
	var mods []ModInfo
	
	err := filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && (filepath.Ext(path) == ".package" || filepath.Ext(path) == ".ts4script") {
			mod := ModInfo{
				Name:        filepath.Base(path),
				InstallDate: info.ModTime(),
				FilePath:    path,
				FileSize:    info.Size(),
			}
			mods = append(mods, mod)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].InstallDate.After(mods[j].InstallDate)
	})
	
	saveRecentMods(mods)
	
	return mods, nil
}

func saveRecentMods(mods []ModInfo) error {
	data, err := json.Marshal(mods)
	if err != nil {
		return err
	}
	
	return compressAndSave(data, "recent_mods.dat")
}

func removeMod(mod ModInfo, list *widget.List) {
	confirmDialog := dialog.NewConfirm(
		"Confirm Removal",
		"Are you sure you want to remove " + mod.Name + "?",
		func(confirmed bool) {
			if confirmed {
				err := os.Remove(mod.FilePath)
				if err != nil {
					dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
					return
				}
				refreshModsList(list)
			}
		},
		fyne.CurrentApp().Driver().AllWindows()[0],
	)
	
	confirmDialog.Show()
}

func showModBrowser(list *widget.List) {
	tabs := fyne.CurrentApp().Driver().AllWindows()[0].Content().(*container.AppTabs)
	
	for _, tab := range tabs.Items {
		if tab.Text == "Browse" {
			tabs.Select(tab)
			break
		}
	}
}

func installMod(reader fyne.URIReadCloser, list *widget.List) {
	settings, err := LoadSettings()
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	filename := filepath.Base(reader.URI().Path())
	targetPath := filepath.Join(settings.ModsDirectory, filename)
	
	if _, err := os.Stat(targetPath); err == nil {
		confirmDialog := dialog.NewConfirm(
			"File Already Exists",
			"A mod with the name " + filename + " already exists. Do you want to overwrite it?",
			func(confirmed bool) {
				if confirmed {
					copyModFile(reader, targetPath, list)
				}
			},
			fyne.CurrentApp().Driver().AllWindows()[0],
		)
		confirmDialog.Show()
	} else {
		copyModFile(reader, targetPath, list)
	}
}

func copyModFile(reader fyne.URIReadCloser, targetPath string, list *widget.List) {
	targetFile, err := os.Create(targetPath)
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	defer targetFile.Close()
	
	_, err = io.Copy(targetFile, reader)
	if err != nil {
		dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		return
	}
	
	dialog.ShowInformation("Mod Installed", "The mod has been successfully installed.", fyne.CurrentApp().Driver().AllWindows()[0])
	
	refreshModsList(list)
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}