package hdlayout

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
	"strings"
	"regexp"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)
func NhName(doc *goquery.Document, win fyne.Window)string{
	var author,title []string
	var dir string
	// Find the authors
	doc.Find("span.before").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		re := regexp.MustCompile(`\[(.*?)\]`)
		author= re.FindStringSubmatch(s.Text())
	})
	//Find Title
	doc.Find("span.pretty").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title = append(title, s.Text())
	})
	if len(author) >= 2 && len(title) >=2 {
		dir = filepath.Join(author[1], title[1])
	}else if len(author) == 0  && len(title) != 0{
		dir = title[len(title)-1]
	}else{
		dialog.ShowInformation("Unknown", fmt.Sprintf("Info: %d, %d,\n %d, %d.",author,title,len(author),len(title)), win)
	}
	return dir
}

func WnName(doc *goquery.Document, win fyne.Window)string{
	var title []string
	var dir string
	//Empty gallery exception
	doc.Find("div.title").Each(func(i int, s *goquery.Selection) {
		re := regexp.MustCompile(`您要訪問的相冊不存在！`)
		title = re.FindStringSubmatch(s.Text())
	})
	//title get wubmatch means gallery doesn't exist.
	if title != nil{
		dialog.ShowInformation("Invalid URL", "Gallery Not Found.", win)
		return dir
	}
	// Find the authors
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the author and title
		title = strings.Split(s.Text()," - 紳士漫畫")
	})
	dir = title[0]
	return dir
}

func GetTitle(site string, win fyne.Window) string{
	var bname string
	// Request the HTML page.
	res, err := http.Get(site)
	if err != nil {
		dialog.ShowInformation("Error", "Get http error.", win)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		dialog.ShowInformation("Invalid URL", "Gallery Not Found.", win)
		return bname
	}

	body, err := ioutil.ReadAll(res.Body)
	reader := bytes.NewReader(body)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		dialog.ShowInformation("Error", "Read query error.", win)
	}

	u,_ := url.Parse(site)
	switch {
	case u.Hostname() == "nhentai.net":
		bname = NhName(doc, win)
	case u.Hostname() == "www.wnacg.com":
		bname = WnName(doc, win)
	default:
		dialog.ShowInformation("Error", "Not supported hostname.", win)
	}
	return bname
}

