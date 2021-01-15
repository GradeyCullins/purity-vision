package db

import (
	"context"
	"fmt"
	"google-vision-filter/src/utils"
	"log"

	"github.com/jackc/pgx/v4"
)

// ImgTableName is the SQL table name for images.
const ImgTableName = "images"

// FindImagesByURI returns images that have matching URI's.
func FindImagesByURI(conn *pgx.Conn, imgURIList []string) ([]Image, error) {
	if len(imgURIList) == 0 {
		return nil, fmt.Errorf("imgURIList cannot be empty")
	}

	imgList := make([]Image, 0)

	uriHashList := make([]string, 0)
	// Build slice of img URI hashes.
	for _, uri := range imgURIList {
		uriHashList = append(uriHashList, utils.Hash(uri))
	}

	// Can't use IN because it is not supported:
	// https://github.com/jackc/pgx/issues/334
	// https://github.com/lib/pq/issues/515
	query := fmt.Sprintf(`SELECT * FROM %s WHERE img_uri_hash = ANY($1)`, ImgTableName)

	rows, err := conn.Query(context.Background(), query, uriHashList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err = rows.Scan(&img.ImgURIHash, &img.Error, &img.Pass, &img.DateAdded); err != nil {
			return nil, err
		}
		log.Printf("Found cached image: %v\n", img)
		imgList = append(imgList, img)
	}

	return imgList, nil
}

// InsertImage inserts the image into the DB.
func InsertImage(conn *pgx.Conn, image Image) error {
	statement := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3)", ImgTableName)
	_, err := conn.Exec(context.Background(), statement, image.ImgURIHash, image.Error.String, image.Pass)
	if err != nil {
		return err
	}

	return nil
}
