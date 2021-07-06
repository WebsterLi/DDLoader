package hdlayout

import (
	"fmt"
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
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
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
	setpathRadio.SetSelected("Default")
	//setpathRadio.Hide()//TODO
	//setpathButton.Hide()//TODO
	setpathBox := container.NewGridWithRows(1,setpathRadio, setpathButton)
	vlist := makeCheckList(25)
	scrollBox := container.NewScroll(container.NewVBox(vlist...))

	serviceSelect := widget.NewSelect([]string{"Artist", "Search"}, func(s string) {})
	serviceSelect.SetSelected("Search")

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
	entryBox := container.NewGridWithRows(1, galleryRadio, applyButton,)

	titleText := widget.NewTextGrid()
	statusText := widget.NewTextGrid()
	infoBox := container.NewHScroll(container.NewVBox(titleText, statusText,))
	progBar := widget.NewProgressBar()
	downprogBox := container.NewGridWithRows(2,infoBox,progBar,)

	downloadButton := &widget.Button{Text: "Download", Importance: widget.HighImportance, OnTapped: func(){

		d := dialog.NewCustom("Info", "Hide", downprogBox, win)
		d.Resize(fyne.NewSize(450, 150))
		selected = 0//reset selected gallery num
		for _, value := range candidate{
				want,_ := value.dl.Get()
				if want{
					selected ++
			}
		}
		ratio := float64(1)/float64(selected)
		switch service{
		case "W":
		case "N":
			if len(candidate)==0{
				dialog.ShowInformation("Error","Please select gallery first.",win)
				return
			}
			d.Show()
			prog := 0.0
			for _, value := range candidate{
				want,_ := value.dl.Get()
				if want{
					dlsite := fmt.Sprintf("https://nhentai.net%s",value.name)
					title := strings.TrimSpace(GetTitle(dlsite, win))
					titleText.SetText(fmt.Sprintf("Title: %s", title))
					statusText.SetText("Status: html info parsing...")
					request := []string{dlsite, service, path, title}
					hbookURL(request,statusText,progBar,prog,ratio)
					prog += ratio
					progBar.SetValue(prog)
				}
			}
		default:
			dialog.ShowInformation("Error","Please select service.",win)
			return
		}
		dialog.ShowInformation("Info", "All task finished.", win)
	}}

	return container.NewBorder(
		container.NewVBox(
			serviceSelect,
			ASEntry,
			entryBox,),
		container.NewVBox(
			setpathBox,
			downloadButton,),
		nil,
		nil,
		scrollBox,
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
	setpathRadio.SetSelected("Default")
	//setpathRadio.Hide()//TODO
	//setpathButton.Hide()//TODO
	setpathRadio.Horizontal = true
	setpathBox := container.NewGridWithRows(1,setpathRadio, setpathButton)
	folderCheck := widget.NewCheck("Create folder with gallery name", func(on bool) {nameget = on})
	folderCheck.Hide()//TODO

	titleText := widget.NewTextGrid()
	statusText := widget.NewTextGrid()
	infoBox := container.NewHScroll(container.NewVBox(titleText, statusText,))
	progBar := widget.NewProgressBar()
	downprogBox := container.NewGridWithRows(2,infoBox,progBar,)

	downloadButton := &widget.Button{Text: "Download", Importance: widget.HighImportance, OnTapped: func(){
		d := dialog.NewCustom("Info", "Hide", downprogBox, win)
		d.Resize(fyne.NewSize(450, 150))
		d.Show()
		site = URLEntry.Text
		var err error
		id, err = strconv.Atoi(URLEntry.Text)

		//if URLEntry is pure string => id
		if err == nil {
			if (id != 0 && service != "") {
				switch service{
				case "N":
					site = fmt.Sprintf("https://nhentai.net/g/%d", id)
				case "W":
					site = fmt.Sprintf("http://www.wnacg.com/photos-index-aid-%d.html", id)
				default :
					dialog.ShowInformation("Error", "Please select service.", win)
				}
				title := strings.TrimSpace(GetTitle(site, win))
				titleText.SetText(fmt.Sprintf("Title: %s", title))
				statusText.SetText("Status: html info parsing...")
				request := []string{site, service, path, title}
				hbookURL(request,statusText,progBar,0.0,1.0)
			} else {
				dialog.ShowInformation("Error", "Please enter gallery id or url.", win)
			}
		//URL case
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
				title := strings.TrimSpace(GetTitle(site, win))
				titleText.SetText(fmt.Sprintf("Title: %s", title))
				statusText.SetText("Status: html info parsing...")
				request := []string{site, service, path}
				hbookURL(request,statusText,progBar,0.0,1.0)
			}
		}
	}}

	return container.NewVBox(
		URLEntry,
		galleryRadio,
		layout.NewSpacer(),
		setpathBox,
		folderCheck,
		downloadButton,
	)
}

