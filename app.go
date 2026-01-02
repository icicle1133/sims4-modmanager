package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func setupApp() fyne.App {
	a := app.New()
	a.Settings().SetTheme(newDarkTheme())
	
	mainWindow := a.NewWindow("Sims 4 Mod Manager")
	mainWindow.Resize(fyne.NewSize(900, 600))
	
	tabs := container.NewAppTabs(
		container.NewTabItem("Mods", setupModsTab()),
		container.NewTabItem("Browse", setupBrowserTab()),
		container.NewTabItem("Settings", setupSettingsTab()),
	)
	
	tabs.SetTabLocation(container.TabLocationTop)
	
	mainWindow.SetContent(tabs)
	mainWindow.ShowAndRun()
	
	return a
}

func welcomeScreen() fyne.CanvasObject {
	return widget.NewLabel("Welcome to Sims 4 Mod Manager")
}