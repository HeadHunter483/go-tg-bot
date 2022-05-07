package forismatic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Aphorism struct {
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	SenderName  string `json:"senderName"`
	SenderLink  string `json:"senderLink"`
	QuoteLink   string `json:"quoteLink"`
}

// getResponse makes http get request by provided url and retrieves
// response body.
func getResponse(url string) ([]byte, error) {
	// get response by URL
	response := make([]byte, 0)

	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	for {
		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)
		response = append(response, bs[:n]...)
		if n == 0 || err != nil {
			break
		}
	}

	return response, nil
}

// GetAphorism gets and unmarshals aphorism from https://forismatic.com/.
func GetAphorism() (*Aphorism, error) {
	url := "http://api.forismatic.com/api/1.0/" +
		"?method=getQuote&format=json&lang=en"
	response, err := getResponse(url)
	aphorism := Aphorism{}

	if err != nil {
		return &aphorism, err
	}

	err = json.Unmarshal(response, &aphorism)
	if err != nil {
		log.Fatal(err)
		return &aphorism, err
	}

	return &aphorism, nil
}

// GetAphorismText gets aphorism as text in format
// "«<aphorism>»\n(c) <author>".
func GetAphorismText() (string, error) {
	result := ""
	aphorism, err := GetAphorism()
	if err != nil {
		log.Fatal(err)
		result = `Something went wrong ಠ~ಠ`
	} else {
		result = fmt.Sprintf(
			"«%s»\n(c) %s", strings.TrimSpace(aphorism.QuoteText),
			aphorism.QuoteAuthor)
	}
	return result, err
}
