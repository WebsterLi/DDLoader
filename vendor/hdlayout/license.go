package hdlayout

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func parseURL(urlStr string) *url.URL {
        link, err := url.Parse(urlStr)
        if err != nil {
                fyne.LogError("Could not parse URL", err)
        }

        return link
}

func licenseScreen(_ fyne.Window) fyne.CanvasObject {

	return container.NewVBox(
		widget.NewLabel("DDownloader"),
		container.NewVBox(
			widget.NewLabel("GUI based on fyne by"),
                        widget.NewHyperlink("fyne.io", parseURL("https://fyne.io/")),
                        widget.NewLabel("-"),
			widget.NewLabel("Icon designed by"),
			widget.NewHyperlink("freepik from Flaticon", parseURL("https://www.freepik.com")),
                        //widget.NewLabel("-"),
                        //widget.NewHyperlink("sponsor", parseURL("https://github.com/sponsors/fyne-io")),
                ),
	)
}
