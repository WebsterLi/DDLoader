// Package main provides various examples of Fyne API capabilities.
package main

import (
	"hdlayout"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {
	//Initial Setting
	a := app.NewWithID("hentai.dl")
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("HDownloader")
	topWindow = w
	content := container.NewMax()
	title := widget.NewLabel("Hentai Downloader")
	//function pointer
	setTutorial := func(t hdlayout.Section) {
		title.SetText(t.Title)
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	split := container.NewHSplit(makeNav(setTutorial), tutorial)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(640, 480))
	w.ShowAndRun()
}

func makeNav(setTutorial func(tutorial hdlayout.Section)) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return hdlayout.SectionIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := hdlayout.SectionIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := hdlayout.Sections[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := hdlayout.Sections[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}
	currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
	tree.Select(currentPref)

	themes := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
