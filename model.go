package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Code struct {
	Store string
	Value string
	Type string
}

type Codes []*Code

const (
	barcodeapi = "https://barcodeapi.org/api"
)
	
var typeMap = map[string]string{
	"code128": "128",
	"ean13": "13",
	"qr": "qr",
}

func NewCode(store string, value string, typeValue string) *Code {
	return &Code{
		Store: store,
		Value: value,
		Type: typeMap[typeValue],
	}
}

func (c *Code) DumpImage() error {
	dataFolder := "data/"+c.Store
	fileName := fmt.Sprintf("%s/%s.jpg", dataFolder, c.Value)
	if _, err := os.Stat(fileName); err == nil {
		return nil
	}

	image, err := c.Image()

	if err != nil {
		return err
	}

	if _, err := os.Stat(dataFolder); err != nil {
		err = os.MkdirAll(dataFolder, 0755)

		if err != nil {
			return err
		}

	}


	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	file.WriteString(image)

	return nil
}

func (c *Code) Image() (string, error) {
	url := fmt.Sprintf("%s/%s/%s?",barcodeapi, c.Type, c.Value)
	request, err := http.Get(url)

	if err != nil {
		return "", err
	}

	rawBody, err := io.ReadAll(request.Body)

	if err != nil {
		return "", err
	}

	return string(rawBody), nil
}
