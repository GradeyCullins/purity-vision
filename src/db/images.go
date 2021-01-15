package db

import (
	"database/sql"
	"fmt"
	"google-vision-filter/src/utils"

	"github.com/lib/pq"
)

// ImgTableName is the SQL table name for images.
const ImgTableName = "images"

// FindImagesByURI returns images that have matching URI's.
func FindImagesByURI(db *sql.DB, imgURIList []string) ([]Image, error) {
	if len(imgURIList) == 0 {
		return nil, fmt.Errorf("imgURIList cannot be empty")
	}

	imgList := make([]Image, 0)

	uriHashList := make([]string, 0)
	// Build slice of img URI hashes.
	for _, uri := range imgURIList {
		uriHashList = append(uriHashList, utils.Hash(uri))
	}

	// Build the SQL 'IN' clause.
	// for i, uri := range imgURIList {
	// 	uriHash := utils.Hash(uri)
	// 	imgHashes[uriHash] = uri
	// 	imgHashIn = imgHashIn + "'" + uriHash + "'"
	// 	if i < len(imgURIList)-1 {
	// 		imgHashIn = imgHashIn + ", "
	// 	}
	// }
	// inClause := "?" + strings.Repeat(", ?", len(imgURIList)-1)

	// Query DB for URI hashes using IN (val1, val2,... valn) syntax
	// query := fmt.Sprintf("SELECT * FROM %s i WHERE i.img_uri_hash IN ($1);", ImgTableName)
	query := `SELECT * FROM images WHERE img_uri_hash = ANY($1)`

	rows, err := db.Query(query, pq.Array(uriHashList))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err = rows.Scan(&img.ImgURIHash, &img.Error, &img.Pass); err != nil {
			return nil, err
		}
		fmt.Printf("Found cached image: %v\n", img)
		imgList = append(imgList, img)
	}

	return imgList, nil
}

// InsertImage inserts the image into the DB.
func InsertImage(db *sql.DB, image Image) error {
	statement := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3)", ImgTableName)
	_, err := db.Exec(statement, image.ImgURIHash, image.Error.String, image.Pass)
	if err != nil {
		return err
	}

	return nil
}
