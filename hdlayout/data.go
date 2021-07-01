package hdlayout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
)

type Section struct {
	Title	string
	View	func(w fyne.Window) fyne.CanvasObject
}

var (
	// Sections defines the metadata for each section
	Sections = map[string]Section{
		"welcome": {"Welcome", welcomeScreen},
		"hbook": {"GalleryDL",
			makeGalleryTab,
		},
		"setting" : {"Setting",
			settings.NewSettings().LoadAppearanceScreen,
		},
	}

	//SectionIndex  defines how our tutorials should be laid out in the index tree
	SectionIndex = map[string][]string{
		"":	{"welcome", "hbook", "setting"},
	}
)
