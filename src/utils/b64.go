package utils

import (
	"bufio"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Base64EncodeF(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	r := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(r)
	return base64.StdEncoding.EncodeToString(content), nil
}

func Base64EncodeR(r io.Reader) (io.Reader, error) {
	content, _ := ioutil.ReadAll(r)
	return strings.NewReader(base64.StdEncoding.EncodeToString(content)), nil
}
