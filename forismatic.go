package main

import (
	"encoding/json"
	"net/http"
)

type Aphorism struct {
	QuoteText string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	SenderName string `json:"senderName"`
	SenderLink string `json:"senderLink"`
	QuoteLink string `json:"quoteLink"`
}

func getResponse(url string) ([]byte, error) {
	// get response by URL
	response := make([]byte, 0)

	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	for true {
		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)
		response = append(response, bs[:n]...)
		if n == 0 || err != nil {
			break
		}
	}

	return response, nil
}

func getAphorism() (*Aphorism, error) {
	// get aphorism from https://forismatic.com/ in json format
	url := "http://api.forismatic.com/api/1.0/" + 
		   "?method=getQuote&format=json&lang=en"
	response, err := getResponse(url)
	aphorism := Aphorism{}

	if err != nil {
		return &aphorism, err
	}

	err = json.Unmarshal(response, &aphorism)

	return &aphorism, nil
}
