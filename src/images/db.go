package images

import (
	"fmt"
	"os"
	"purity-vision-filter/src/utils"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

// ImgTableName is the SQL table name for images.
const ImgTableName = "images"

// FindByURI returns images that have matching URI's.
func FindByURI(conn *pg.DB, imgURIList []string) ([]Image, error) {
	var imgList []Image

	if len(imgURIList) == 0 {
		return nil, fmt.Errorf("imgURIList cannot be empty")
	}

	// Build slice of img URI hashes.
	uriHashList := make([]string, 0)
	for _, uri := range imgURIList {
		uriHashList = append(uriHashList, utils.Hash(uri))
	}

	conn.Model(&imgList).Where("uri IN (?)", pg.In(imgURIList)).Select()

	for _, img := range imgList {
		logger.Debug().Msgf("Found cached image: %s...", img.URI)
	}

	logger.Debug().Msgf("Cached %d/%d images", len(imgList), len(imgURIList))

	return imgList, nil
}

// Insert inserts the image into the DB.
func Insert(conn *pg.DB, image *Image) error {
	_, err := conn.Model(image).Insert()
	if err != nil {
		return err
	}
	logger.Debug().Msgf("inserted image: %s", image.URI)

	return nil
}

// DeleteByURI deletes the images with matching URI.
func DeleteByURI(conn *pg.DB, uri string) error {
	img := Image{URI: uri}

	if _, err := conn.Model(&img).Where("uri = ?", uri).Delete(); err != nil {
		return err
	}

	return nil
}
