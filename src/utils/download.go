package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Download(uri string) (string, error) {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return "", err
	}

	// Create temp file for download.
	f, err := ioutil.TempFile("", "purity-img")
	if err != nil {
		return "", err
	}

	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
