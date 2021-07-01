package hdlayout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
)

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"welcome": {"Welcome", "", welcomeScreen},
		"hbook": {"GalleryDL",
			"Enter ID/URL to download gallery.",
			makeDownloadTab,
		},
		"split": {"Split_Test",
			"A split container divides the container in two pieces that the user can resize.",
			makeSplitTab,
		},
		"card": {"Card_Test",
			"Group content and widgets.",
			makeCardTab,
		},
		"progress": {"Progress_Test",
			"Show duration or the need to wait for a task.",
			makeProgressTab,
		},
		"setting" : {"Setting",
			"Fyne theme appearence setting.",

			settings.NewSettings().LoadAppearanceScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":	{"welcome", "hbook", "split", "card", "progress", "setting"},
	}
)
