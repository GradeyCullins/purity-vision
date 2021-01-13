package db

import (
	"database/sql"
	"fmt"
	"google-vision-filter/src/utils"
)

// FindImagesByURI returns images that have matching URI's.
func FindImagesByURI(conn *sql.DB, imgURIList []string) ([]Image, error) {
	if len(imgURIList) == 0 {
		return nil, fmt.Errorf("imgURIList cannot be empty")
	}

	imgList := make([]Image, 0)
	var imgHashIn string = ""
	imgHashes := make(map[string]string)

	// Build the SQL 'IN' clause.
	for i, uri := range imgURIList {
		uriHash := utils.Hash(uri)
		imgHashes[uriHash] = uri
		imgHashIn = imgHashIn + "'" + uriHash + "'"
		if i < len(imgURIList)-1 {
			imgHashIn = imgHashIn + ", "
		}
	}

	// Query DB for hashes using IN (val1, val2,... valn) syntax
	rows, err := conn.Query("SELECT * from images i WHERE i.img_uri_hash IN (" + imgHashIn + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var img Image
		if err = rows.Scan(&img.ImgURIHash, &img.Error, &img.Pass); err != nil {
			return nil, err
		}
		imgList = append(imgList, img)
	}

	return imgList, nil
}
