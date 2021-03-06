package hdlayout

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"path/filepath"

	"fyne.io/fyne/v2/widget"
)

var saveTo string

func panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func get(url string) string {
	resp, err := http.Get(url)
	panic(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	panic(err)
	return string(body)
}

func getFileName(url string) string {
	slices := strings.Split(url, "/")
	slices = strings.Split(slices[len(slices)-1], "?")
	return slices[0]
}

func getFileExt(url string) string {
	slices := strings.Split(url, ".")
	slices = strings.Split(slices[len(slices)-1], "?")
	return slices[0]
}

func isDirExists(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

func removeDuplicateValues(stringSlice []string) []string {
    keys := make(map[string]bool)
    list := []string{}

    // If the key(values of the slice) is not equal
    // to the already present value in new slice (list)
    // then we append it. else we jump on another element.
    for _, entry := range stringSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

func downImg(url, fn string) {
	resp, err := http.Get(url)
	panic(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if !isDirExists(saveTo) {
		err = os.MkdirAll(saveTo, 0755)
		panic(err)
	}
	var filename string
	if fn != "" {
		ext := getFileExt(url)
		filename = filepath.Join(saveTo, (fn + "." + ext))
	} else {
		filename = filepath.Join(saveTo, getFileName(url))
	}
	err = ioutil.WriteFile(filename, body, 0755)
	panic(err)
}

func NhPage(url string, statusText *widget.TextGrid, progBar *widget.ProgressBar, start, ratio float64) {
	body := get(url)
	//fmt.Println(body)
	var wg sync.WaitGroup
	//Exclude "cover" "thumb" keyword, select string until "jpg" "png" keyword.
	re := regexp.MustCompile("galleries/[^(humb)|^(cover)]+[.jpg|.png]")
	photo_urls := re.FindAllString(body, -1)
	//remove duplicate strings
	photo_urls = removeDuplicateValues(photo_urls)
	//fmt.Println(photo_urls)
	tokens := make(chan int, 10)
	ratio /= (float64)(len(photo_urls))
	for _, photo_url := range photo_urls {
		wg.Add(1)
		tokens <- 1
		photo_url = "http://i.nhentai.net/" + photo_url
		//remove "t" char in img name
		replacer := strings.NewReplacer("t.jpg", ".jpg", "t.png", ".png")
		photo_url = replacer.Replace(photo_url)
		//fmt.Println(photo_url)
		go func(img string) {
			fn := strings.TrimSuffix(filepath.Base(img), filepath.Ext(img))
			downImg(img, fn)
			statusText.SetText(fmt.Sprintf("Downloading: %s", fn))
			start += ratio
			progBar.SetValue(start)
			<-tokens
			defer wg.Done()
		}(photo_url)
	}
	wg.Wait()
	return
}

func WnPage(url string, statusText *widget.TextGrid, progBar *widget.ProgressBar, start, ratio float64) {
	body := get(url)
	var wg sync.WaitGroup
	re := regexp.MustCompile("photos-view-id-[0-9]+.html")
	photo_urls := re.FindAllString(body, -1)
	//fmt.Println(photo_urls)
	tokens := make(chan int, 5)
	tmpratio := ratio/(float64)(len(photo_urls))
	tmpstart := start
	for _, photo_url := range photo_urls {
		wg.Add(1)
		tokens <- 1
		photo_url = "http://www.wnacg.com/" + photo_url
		go func(url string) {
			body := get(url)
			re := regexp.MustCompile("img[0-9].wnacg.com/data/[^\"]+")
			img := "http://" + re.FindString(body)
			re = regexp.MustCompile(`alt="([^"]+)"`)
			search := re.FindAllStringSubmatch(body, -1)
			//fmt.Println(len(search))
			var fn string
			if search == nil || len(search) < 2 {
				fn = ""
			} else {
				fn = search[len(search) - 1][1]
			}
			downImg(img, fn)
			statusText.SetText(fmt.Sprintf("Downloading: %s", fn))
			tmpstart += tmpratio
			progBar.SetValue(start)
			<-tokens
			defer wg.Done()
		}(photo_url)
	}
	wg.Wait()
	re = regexp.MustCompile(`<span class="next"><a href="([^"]+)">`)
	urls := re.FindStringSubmatch(body)
	if len(urls) == 0 {
		return
	}
	next_page := urls[1]
	// fmt.Println(next_page)
	next_page = "http://www.wnacg.com/" + next_page
	WnPage(next_page, statusText, progBar, start, ratio)
}

func hbookURL(request []string, statusText *widget.TextGrid, progBar *widget.ProgressBar, start, ratio float64) {
	url := request[0]
	service := request[1]
	path := request[2]
	title := request[3]
	//turn to url
	switch service{
	case "n","N":
		saveTo = filepath.Join(path, "N_gallery", title)
		NhPage(url, statusText, progBar, start, ratio)
	case "w","W":
		saveTo = filepath.Join(path, "W_gallery", title)
		WnPage(url, statusText, progBar, start, ratio)
	default:
		fmt.Printf("Undefined Symbol.\n")
	}
}
