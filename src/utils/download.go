package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Download fetches a web resource at "uri" and returns a file handle to the downloaded response.
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

	if res.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("Request to download image at: %s returned a 404", uri)
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
