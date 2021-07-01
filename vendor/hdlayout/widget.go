package hdlayout

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
)

const (
	loremIpsum = `This is a book.`
)

var (
	progress    *widget.ProgressBar
	infProgress *widget.ProgressBarInfinite
	endProgress chan interface{}
)

func makeCardTab(_ fyne.Window) fyne.CanvasObject {
	card1 := widget.NewCard("Book a table", "Which time suits?",
		widget.NewRadioGroup([]string{"6:30pm", "7:00pm", "7:45pm"}, func(string) {}))
	card2 := widget.NewCard("With media", "No content, with image", nil)
	card2.Image = canvas.NewImageFromResource(theme.FyneLogo())
	card3 := widget.NewCard("Title 3", "Another card", widget.NewLabel("Content"))
	return container.NewGridWithColumns(2, container.NewVBox(card1, card3),
		container.NewVBox(card2))
}

func makeDownloadTab(win fyne.Window) fyne.CanvasObject {
	galleryRadio := widget.NewRadioGroup([]string{"NHentai", "Wnacg"}, func(s string) {})
	galleryRadio.Horizontal = true
	IDEntry := widget.NewEntry()
	IDEntry.SetPlaceHolder("ID")
	galleryBox := container.NewGridWithRows(1, galleryRadio, container.NewHBox(widget.NewLabel("ID"), IDEntry))
	URLEntry := widget.NewSelectEntry([]string{"https://nhentai.net/g/<id>", "https://www.wnacg.com/photos-index-aid-<id>.html"})
	URLEntry.PlaceHolder = "Type or select"
	setpathButton := widget.NewButton("Select folder",func(){
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}

			children, err := list.List()
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			dialog.ShowInformation("Folder Open", out, win)
		}, win)
	})
	setpathRadio := widget.NewRadioGroup([]string{"Default", "Set path"}, func(s string) {
		switch s{
		case "Default":
			setpathButton.Disable()
		case "Set path":
			setpathButton.Enable()
		default:
			setpathButton.Disable()
		}
	})
	setpathBox := container.NewGridWithRows(1,setpathRadio, setpathButton)
	folderCheck := widget.NewCheck("Create folder with gallery name", func(on bool) { })
	downloadButton := &widget.Button{Text: "Download", Importance: widget.HighImportance, OnTapped: func(){}}

	return container.NewVBox(
		widget.NewSelect([]string{"ID", "URL"}, func(s string) {
			switch s{
			case "ID":
				galleryBox.Show()
				URLEntry.Hide()
			case "URL":
				galleryBox.Hide()
				URLEntry.Show()
			default:
				galleryBox.Hide()
				URLEntry.Hide()
			}
		}),
		galleryBox,
		URLEntry,
		layout.NewSpacer(),
		setpathBox,
		folderCheck,
		downloadButton,
	)
}

func makeProgressTab(_ fyne.Window) fyne.CanvasObject {
	stopProgress()

	progress = widget.NewProgressBar()

	infProgress = widget.NewProgressBarInfinite()
	endProgress = make(chan interface{}, 1)
	startProgress()

	return container.NewVBox(
		widget.NewLabel("Percent"), progress,
		widget.NewLabel("Infinite"), infProgress)
}

func startProgress() {
	progress.SetValue(0)
	select { // ignore stale end message
	case <-endProgress:
	default:
	}

	go func() {
		end := endProgress
		num := 0.0
		for num < 1.0 {
			time.Sleep(16 * time.Millisecond)
			select {
			case <-end:
				return
			default:
			}

			progress.SetValue(num)
			num += 0.002
		}

		progress.SetValue(1)

		// TODO make sure this resets when we hide etc...
		stopProgress()
	}()
	infProgress.Start()
}

func stopProgress() {
	if !infProgress.Running() {
		return
	}

	infProgress.Stop()
	endProgress <- struct{}{}
}

// widgetScreen shows a panel containing widget demos
func widgetScreen(_ fyne.Window) fyne.CanvasObject {
	content := container.NewVBox(
		widget.NewLabel("Labels"),
		widget.NewButtonWithIcon("Icons", theme.HomeIcon(), func() {}),
		widget.NewSlider(0, 1))
	return container.NewCenter(content)
}

