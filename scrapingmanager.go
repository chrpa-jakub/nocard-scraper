package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	nocardUrl = "https://nocard.cz"
)

func Start() {
	counter := 0
	ticker := time.NewTicker(time.Millisecond*50)
	defer ticker.Stop()

	for range ticker.C {
		counter++
		fmt.Println(counter)
		Scrape()

		if counter == 500 {
			fmt.Println("Scraping done, check the data folder.")
			return
		}
	}

}

func Scrape() {
	html, err := NocardHtml()

	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	codes := FilterCodes(html)

	var wg sync.WaitGroup
	for _, v := range codes {
		wg.Add(1)
		go func() {
			v.DumpImage()
			defer wg.Done()
		}()
	}

	wg.Wait()
}

func NocardHtml() (string, error) {
	request, err := http.Get(nocardUrl)

	if err != nil {
		return "", err
	}

	defer request.Body.Close()

	requestBody, err := io.ReadAll(request.Body)

	if err != nil {
		return "", err
	}


	return string(requestBody), nil;	
}

func FilterCodes(rawHtml string) Codes {
	htmlSplit := strings.Split(rawHtml, "\n")
	codes := Codes{}

	for _, line := range htmlSplit {
		if strings.Contains(line, `"card `) {
			lineSplit := strings.Split(line, `"`)

			if len(lineSplit) < 10 {
				break
			}

			codes = append(codes, NewCode(lineSplit[5], lineSplit[9], lineSplit[7]))
		}
	}
	
	return codes
}
