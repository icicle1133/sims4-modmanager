package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type darkTheme struct{
	base fyne.Theme
}

func newDarkTheme() fyne.Theme {
	return &darkTheme{base: theme.DefaultTheme()}
}

func (d *darkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 30, G: 30, B: 30, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 190, G: 190, B: 190, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 65, G: 105, B: 225, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 128}
	default:
		return d.base.Color(name, variant)
	}
}

func (d *darkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return d.base.Font(style)
}

func (d *darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return d.base.Icon(name)
}

func (d *darkTheme) Size(name fyne.ThemeSizeName) float32 {
	return d.base.Size(name)
}