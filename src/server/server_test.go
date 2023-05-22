package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"purity-vision-filter/src/config"
	"purity-vision-filter/src/db"
	"purity-vision-filter/src/images"
	"testing"

	"github.com/go-pg/pg/v10"
)

type TestServe struct {
}

func (s *TestServe) Init(_conn *pg.DB) {
	conn = _conn
}

func TestHealthHandler(t *testing.T) {
	testHealthNoBody(t)
	testHealthJunkBody(t)
}

func testHealthNoBody(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Error("Failed to create test HTTP request")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(health)

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Health endpoint expected response 200 but got %d", rr.Code)
	}
}

type junkData struct {
	Name  string
	Color int
}

// The health endpoint given junk POST data should still simply return a 200 code.
func testHealthJunkBody(t *testing.T) {
	someData := junkData{
		Name:  "pil",
		Color: 221,
	}
	b, err := json.Marshal(someData)
	if err != nil {
		t.Error("Failed to marshal request body struct")
	}
	r := bytes.NewReader(b)
	req, err := http.NewRequest("POST", "/health", r)
	if err != nil {
		t.Error("Failed to create test HTTP request")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(health)

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Health endpoint expected response 200 but got %d", rr.Code)
	}
}

func TestBatchImgFilterHandler(t *testing.T) {
	conn, err := db.Init(config.DefaultDBTestName)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	s := TestServe{}
	s.Init(conn)
	uri := "https://i.ytimg.com/vi/19VZZpzbh6s/maxresdefault.jpg"

	fr := BatchImgFilterReq{
		ImgURIList: []string{},
	}

	var errRes ErrorRes
	res, err := testBatchImgFilterHandler(fr)
	if err != nil {
		t.Error("Shouldn't have thrown an error")
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&errRes); err != nil {
		t.Error("JSON body missing or malformed")
	}

	if res.Code != 400 || errRes.Message != "ImgUriList cannot be empty" {
		t.Error("Web server should have returned a 400 because the ImgURIList was empty")
	}

	fr = BatchImgFilterReq{
		ImgURIList: []string{uri},
	}

	res, err = testBatchImgFilterHandler(fr)
	if err != nil {
		t.Error(err)
	}
	if res.Code != 200 {
		t.Error("Web server should have returned a 200")
	}
	var fRes BatchImgFilterRes
	json.Unmarshal(res.Body.Bytes(), &fRes)
	if len(fRes) != 1 || fRes[0].Pass != true {
		t.Error("Handler didn't return the right results")
	}

	// Delete the img from the DB.
	if err = images.DeleteByURI(conn, uri); err != nil {
		t.Log(err)
	}
}

// TODO rename to something more descriptive
func testBatchImgFilterHandler(fr BatchImgFilterReq) (*httptest.ResponseRecorder, error) {
	b, err := json.Marshal(fr)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body struct")
	}
	r := bytes.NewReader(b)

	req, err := http.NewRequest("POST", "/filter", r)
	if err != nil {
		return nil, errors.New("Failed to create test HTTP request")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleBatchFilter(logger))

	handler.ServeHTTP(rr, req)

	return rr, nil
}
