package db

import (
	"google-vision-filter/src/config"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
)

var conn *pg.DB
var tx *pg.Tx
var err error
var imgURIList = []string{
	"https://hatrabbits.com/wp-content/uploads/2017/01/random.jpg",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT1ZgCJADylizZLNnOnyuhtwR2qVk5yOi0UoQ&usqp=CAU",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRKsJoGKlOJnxl-GNgfUtluGobgx_M8JBdsng&usqp=CAU",
}

func TestMain(m *testing.M) {
	conn, err = InitDB(config.DefaultDBTestName)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestInsertImage(t *testing.T) {
	for _, uri := range imgURIList {
		err = InsertImage(conn, NewImage(uri, "", true, time.Now()))
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestFindImagesByURI(t *testing.T) {
	smallURIList := imgURIList[:1]

	imgList, err := FindImagesByURI(conn, smallURIList)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(imgList) != 1 {
		t.Fatalf("Expected 1 image in response but received %d", len(imgList))
		t.FailNow()
	}

	smallURIList = []string{}
	imgList, err = FindImagesByURI(conn, smallURIList)
	if err == nil {
		t.Fatal("Expected FindImagesByURI to return an error because imgURIList cannot be empty")
	}
}

func TestDeleteImagesByURI(t *testing.T) {
	for _, uri := range imgURIList {
		err = DeleteImageByURI(conn, uri)
		if err != nil {
			t.Fatal(err)
		}
	}
}
