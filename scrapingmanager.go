package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	nocardUrl = "https://nocard.cz"
)

func Start() {
	Scrape()
	fmt.Println("Scraping done, check the data folder.")
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
	codes := Codes{}

	jsonData := strings.Split(rawHtml, "<script>")[2]
	jsonData = strings.Split(jsonData, "</script>")[0]
	jsonData = strings.Join(strings.Split(jsonData, "{")[1:], "{")
	jsonData = strings.Split(jsonData, ";")[0]
	jsonData = "{"+jsonData

	var codeMap map[string]CodesMap

	if err := json.Unmarshal([]byte(jsonData), &codeMap); err != nil {
		panic(err)
	}

	for key, val := range codeMap {
		codes = append(codes, val.extractCodes(key)...)
	}

	return codes
}
