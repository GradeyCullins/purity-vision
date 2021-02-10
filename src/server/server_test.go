package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google-vision-filter/src/config"
	"google-vision-filter/src/db"
	"google-vision-filter/src/images"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg/v10"
)

type TestServe struct {
}

func (s *TestServe) Init(_conn *pg.DB) {
	conn = _conn
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
	fRes = BatchImgFilterRes{
		ImgFilterResList: []ImgFilterRes{},
	}
	json.Unmarshal(res.Body.Bytes(), &fRes)
	if len(fRes.ImgFilterResList) != 1 || fRes.ImgFilterResList[0].Pass != true {
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
		return nil, fmt.Errorf("Failed to create test HTTP request")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(batchImgFilterHandler)

	handler.ServeHTTP(rr, req)

	return rr, nil
}
