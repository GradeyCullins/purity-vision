package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn
var err error
var imgURIList = []string{
	"https://hatrabbits.com/wp-content/uploads/2017/01/random.jpg",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT1ZgCJADylizZLNnOnyuhtwR2qVk5yOi0UoQ&usqp=CAU",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRKsJoGKlOJnxl-GNgfUtluGobgx_M8JBdsng&usqp=CAU",
}

func TestMain(m *testing.M) {
	conn, err = InitDB("purity_test")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	// Don't persist any changes to the test database after tests are complete.
	tx.Rollback(context.Background())

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
