package hdlayout

import (
	"fmt"
	"time"
	"strconv"

	"bytes"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"fyne.io/fyne/v2/data/binding"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
)

var (
	progress    *widget.ProgressBar
	infProgress *widget.ProgressBarInfinite
	endProgress chan interface{}
)
type galdl struct{
	name string
	dl binding.Bool
}

func makeGalleryTab(win fyne.Window) fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItem("ID/URL", makeGIUTab(win)),
		container.NewTabItem("Artist/Search", makeGASTab(win)),
		container.NewTabItem("Tag", widget.NewLabel("Coming soon...")),
	)
}

func makeCheckList(num int)[]fyne.CanvasObject{
	var items []fyne.CanvasObject
        for i := 1; i <= num; i++ {
		check := widget.NewCheck(fmt.Sprintf("Check %d", i), func(on bool) {})
                items = append(items, check)
		check.Hide()
        }
        return items
}

func nReadSearch(site string, win fyne.Window)(map[string]galdl, int){
	m := make(map[string]galdl)
	var page int
	// Request the HTML page.
	res, err := http.Get(site)
	if err != nil {
		dialog.ShowInformation("Error", "Get http error.", win)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		dialog.ShowInformation("Invalid URL", "Gallery Not Found.", win)
		return m, page
	}

	body, err := ioutil.ReadAll(res.Body)
	reader := bytes.NewReader(body)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		dialog.ShowInformation("Error", "Read query error.", win)
	}
	// Find the galleries
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			classes, ok := s.Attr("class")
			if !ok {
				for p := s.Parent(); p.Size() > 0 && !ok; p = p.Parent() {
					classes, ok = p.Attr("class")
				}
			}
			if classes == "cover"{
				key := s.Text()
				key = key[strings.IndexByte(key, '>')+1:]//get all char after '>' symbol from s.text() 

				m[key] = galdl{href, binding.NewBool()}
			}else if classes == "last"{
				last := href[len(href)-3:]
				last = last[strings.IndexByte(last, '=')+1:]//get all char after '=' symbol from href 
				page,_ = strconv.Atoi(last)
			}
		}
	})
	return m, page
}

func makeGASTab(win fyne.Window) fyne.CanvasObject {
	var candidate map[string]galdl
	var selected, end int
	var service, site, path string
	galleryRadio := widget.NewRadioGroup([]string{"NHentai", "Wnacg"}, func(s string) {
		switch s{
		case "NHentai":
			service = "N"
		case "Wnacg":
			service = "W"
		default:
			service = ""
		}
	})
	galleryRadio.Horizontal = true
	ASEntry := widget.NewEntry()
	ASEntry.PlaceHolder = "Please enter name or keyword."
	setpathButton := widget.NewButton("Select folder",func(){
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}
			path = list.Path()
			dialog.ShowInformation("Folder Selected", path, win)
		}, win)
	})
	setpathRadio := widget.NewRadioGroup([]string{"Default", "Set path"}, func(s string) {
		switch s{
		case "Default":
			setpathButton.Disable()
			path = "."
		case "Set path":
			setpathButton.Enable()
		default:
			setpathButton.Disable()
		}
	})
	setpathRadio.Horizontal = true
	setpathBox := container.NewGridWithRows(1,setpathRadio, setpathButton)
	serviceSelect := widget.NewSelect([]string{"Artist", "Search"}, func(s string) {})
	serviceSelect.SetSelected("Artist")
	vlist := makeCheckList(25)
	vscrollBox := container.NewVScroll(container.NewVBox(vlist...))

	applyButton := &widget.Button{Text: "Apply", Importance: widget.MediumImportance, OnTapped: func(){
		switch serviceSelect.SelectedIndex(){
		case 0:
			switch service{
			case "W":
				dialog.ShowInformation("Invalid command", "Wnacg does not support Artist label.", win)
			case "N":
				site = fmt.Sprintf("https://nhentai.net/artist/%s",ASEntry.Text)
				candidate, end = nReadSearch(site,win)
				count := 0
				for key, value := range candidate{
					if count < len(vlist){
						vlist[count] = widget.NewCheckWithData(key, value.dl)
					}else{
						//should not work
						dialog.ShowInformation("Error", "Page list out of range.", win)
					}
					count++
				}
				fmt.Println("Total pages: ", end)
			default:
			}
		case 1:
			switch service{
			case "W":
				site = fmt.Sprintf("https://www.wnacg.com/search/?q=%s", ASEntry.Text)
				dialog.ShowInformation("Test", site, win)
			case "N":
				site = fmt.Sprintf("https://nhentai.net/search/?q=%s", ASEntry.Text)
				candidate, end = nReadSearch(site,win)
				count := 0
				for key, value := range candidate{
					if count < len(vlist){
						vlist[count] = widget.NewCheckWithData(key, value.dl)
					}else{
						//should not work
						dialog.ShowInformation("Error", "Page list out of range.", win)
					}
					count++
				}
				fmt.Println("Total pages: ", end)
			default:
			}
		default:
			dialog.ShowInformation("Empty Selection", "Please select mode.", win)
		}
	}}
	downloadButton := &widget.Button{Text: "Download", Importance: widget.HighImportance, OnTapped: func(){
		d := dialog.NewProgress("Info", "Downloading...", win)
		d.Show()
		selected = 0//reset selected gallery num
		for _, value := range candidate{
				want,_ := value.dl.Get()
				if want{
					selected ++
			}
		}
		switch service{
		case "W":
		case "N":
			//TODO change to go func with progress bar
			prog := 0.0
			for _, value := range candidate{
				want,_ := value.dl.Get()
				if want{
					prog += float64(1)/float64(selected)
					d.SetValue(prog)
					dlsite := fmt.Sprintf("https://nhentai.net%s",value.name)
					hbookURL(dlsite,service,path,win)
				}
			}
		default:
		}
		d.Hide()
		dialog.ShowInformation("Info", "All download finished.", win)
	}}
	optionBox := container.NewGridWithRows(2,applyButton,downloadButton,)

	return container.NewBorder(
		container.NewVBox(
			serviceSelect,
			galleryRadio,
			ASEntry,),
		container.NewVBox(
			setpathBox,
			optionBox,),
		nil,
		nil,
		vscrollBox,
	)
}


func makeGIUTab(win fyne.Window) fyne.CanvasObject {
	var nameget bool
	var id int
	var service, site, path string
	galleryRadio := widget.NewRadioGroup([]string{"NHentai", "Wnacg"}, func(s string) {
		switch s{
		case "NHentai":
			service = "N"
		case "Wnacg":
			service = "W"
		default:
			service = ""
		}
	})
	galleryRadio.Horizontal = true
	URLEntry := widget.NewSelectEntry([]string{
		"<id>",
		"https://nhentai.net/g/<id>",
		"https://www.wnacg.com/photos-index-aid-<id>.html",
	})
	URLEntry.PlaceHolder = "Please enter ID or URL."

	setpathButton := widget.NewButton("Select folder",func(){
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}
			path = list.Path()
			dialog.ShowInformation("Folder Selected", path, win)
		}, win)
	})
	setpathRadio := widget.NewRadioGroup([]string{"Default", "Set path"}, func(s string) {
		switch s{
		case "Default":
			setpathButton.Disable()
			path = "."
		case "Set path":
			setpathButton.Enable()
		default:
			setpathButton.Disable()
		}
	})
	setpathRadio.Horizontal = true
	setpathBox := container.NewGridWithRows(1,setpathRadio, setpathButton)
	folderCheck := widget.NewCheck("Create folder with gallery name", func(on bool) {nameget = on})
	folderCheck.Hide()//TODO
	downloadButton := &widget.Button{Text: "Download", Importance: widget.HighImportance, OnTapped: func(){
		site = URLEntry.Text
		var err error
		id, err = strconv.Atoi(URLEntry.Text)

		if err == nil {
			if (id != 0 && service != "") {
				hbookID(id,service,path,win)
			} else {
				dialog.ShowInformation("Empty Entry", "Please enter gallery id.", win)
			}
		} else {
			if site != ""{
				u,_ := url.Parse(site)
				switch {
				case u.Hostname() == "nhentai.net":
					service = "N"
				case u.Hostname() == "www.wnacg.com":
					service = "W"
				default:
					dialog.ShowInformation("Error", "Not supported hostname.", win)
				}
			}
			if service != ""{
				hbookURL(site,service,path,win)
			} else {
				dialog.ShowInformation("Empty Entry", "Please enter gallery url.", win)
			}
		}
	}}

	return container.NewVBox(
		galleryRadio,
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

