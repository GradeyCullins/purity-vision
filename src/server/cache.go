package server

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"google-vision-filter/src/db"
)

func checkImgCache(conn *sql.DB, imgURIList []string) (*BatchImgFilterRes, error) {
	var cacheRes BatchImgFilterRes
	var imgFilterResList []ImgFilterRes
	var imgHashIn string = ""
	imgHashes := make(map[string]string)

	// Build the SQL 'IN' clause.
	for i, uri := range imgURIList {
		h := sha256.New()
		h.Write([]byte(uri))
		uriHash := fmt.Sprintf("%x", h.Sum(nil))
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
		var imgCacheEntry db.ImageCacheEntry
		if err = rows.Scan(&imgCacheEntry.ImgURIHash, &imgCacheEntry.Error, &imgCacheEntry.Pass); err != nil {
			return nil, err
		}
		filterRes := ImgFilterRes{
			ImgURI: imgHashes[imgCacheEntry.ImgURIHash],
			Error:  imgCacheEntry.Error.String,
			Pass:   imgCacheEntry.Pass,
		}
		imgFilterResList = append(imgFilterResList, filterRes)
	}

	cacheRes.ImgFilterResList = imgFilterResList
	return &cacheRes, nil
}
