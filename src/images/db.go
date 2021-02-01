package images

import (
	"fmt"
	"google-vision-filter/src/utils"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
// var logger zerolog.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
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

	conn.Model(&imgList).Where("img_uri_hash IN (?)", pg.In(uriHashList)).Select()

	for _, img := range imgList {
		logger.Info().Msgf("Found cached image: %s...", img.ImgURIHash[:8])
	}

	return imgList, nil
}

// Insert inserts the image into the DB.
func Insert(conn *pg.DB, image *Image) error {
	_, err := conn.Model(image).Insert()
	if err != nil {
		return err
	}
	logger.Info().Msgf("inserted image: %s", image.ImgURIHash)

	return nil
}

// DeleteByURI deletes the images with matching URI.
func DeleteByURI(conn *pg.DB, uri string) error {
	hash := utils.Hash(uri)
	img := Image{ImgURIHash: hash}

	if _, err := conn.Model(&img).Where("img_uri_hash = ?", hash).Delete(); err != nil {
		return err
	}

	return nil
}
