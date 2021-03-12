package images

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"purity-vision-filter/src/config"
	"purity-vision-filter/src/db"
	"purity-vision-filter/src/utils"
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
	conn, err = db.Init(config.DefaultDBTestName)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestNewImage(t *testing.T) {
	uri := "https://google.com"
	time := time.Now()
	fakeHash := utils.Hash("some string")
	expected := Image{
		fakeHash,
		uri,
		sql.NullString{},
		true,
		time,
	}
	img, err := NewImage(fakeHash, uri, fmt.Errorf(""), true, time)
	if err != nil {
		t.Fatal(err)
	}

	if (expected.Hash != img.Hash) ||
		(expected.Error != img.Error) ||
		(expected.Pass != img.Pass) ||
		(expected.DateAdded != img.DateAdded) {
		t.Fatalf("Expected %v to equal %v", expected, img)
	}
}

func TestInsertImage(t *testing.T) {
	for _, uri := range imgURIList {
		fakeHash := utils.Hash(uri)
		img, err := NewImage(fakeHash, uri, fmt.Errorf(""), true, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		err = Insert(conn, img)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestFindImagesByURI(t *testing.T) {
	smallURIList := imgURIList[:1]

	imgList, err := FindByURI(conn, smallURIList)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(imgList) != 1 {
		t.Fatalf("Expected 1 image in response but received %d", len(imgList))
		t.FailNow()
	}

	smallURIList = []string{}
	imgList, err = FindByURI(conn, smallURIList)
	if err == nil {
		t.Fatal("Expected FindImagesByURI to return an error because imgURIList cannot be empty")
	}
}

func TestDeleteImagesByURI(t *testing.T) {
	for _, uri := range imgURIList {
		err = DeleteByURI(conn, uri)
		if err != nil {
			t.Fatal(err)
		}
	}
}
