package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func setupSettingsTab() fyne.CanvasObject {
	settings, err := LoadSettings()
	if err != nil {
		settings = AppSettings{ModsDirectory: DefaultModsPath}
	}

	pathEntry := widget.NewEntry()
	pathEntry.SetText(settings.ModsDirectory)

	browseButton := widget.NewButton("Browse", func() {
		dlg := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			pathEntry.SetText(uri.Path())
		}, fyne.CurrentApp().Driver().AllWindows()[0])
		dlg.Show()
	})

	pathRow := container.NewBorder(nil, nil, nil, browseButton, pathEntry)
	
	saveButton := widget.NewButton("Save Settings", func() {
		settings.ModsDirectory = pathEntry.Text
		err := SaveSettings(settings)
		if err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		dialog.ShowInformation("Success", "Settings saved successfully", fyne.CurrentApp().Driver().AllWindows()[0])
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Mods Directory", Widget: pathRow},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			saveButton.OnTapped()
		},
	}

	return container.NewVBox(
		widget.NewLabel("Settings"),
		form,
	)
}