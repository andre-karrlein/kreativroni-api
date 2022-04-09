package util

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Etsy_request(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	api_key := os.Getenv("API_KEY")
	req.Header.Add("x-api-key", api_key)

	lowerCaseHeader := make(http.Header)

	for key, value := range req.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}

	req.Header = lowerCaseHeader
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
