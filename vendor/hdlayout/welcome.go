package hdlayout

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(228, 167))

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Hentai Downloader", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
	))
}
