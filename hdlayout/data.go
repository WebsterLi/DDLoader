package hdlayout

import (
	"fyne.io/fyne/v2"
)

type Section struct {
	Title	string
	View	func(w fyne.Window) fyne.CanvasObject
}

var (
	// Sections defines the metadata for each section
	Sections = map[string]Section{
		"hbook": {"GalleryDL",
			makeGalleryTab,
		},
		"theme" : {"Theme",
			NewSettings().LoadAppearanceScreen,
		},
		"license": {"License",
			licenseScreen,
		},
	}

	//SectionIndex  defines how our tutorials should be laid out in the index tree
	SectionIndex = map[string][]string{
		"":	{"hbook", "theme", "license", },
	}
)
