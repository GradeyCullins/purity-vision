package images

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

// ImgTableName is the SQL table name for images.
const ImgTableName = "images"

// FindByURI returns an image with the matching URI.
func FindByURI(conn *pg.DB, imgURI string) (*Image, error) {
	var img Image

	err := conn.Model(&img).Where("uri = ?", imgURI).Select()
	if err != nil {
		log.Error().Msgf("err: %v", err)
		return nil, nil
	}

	return &img, nil
}

// FindAllByURI returns images that have matching URI's.
func FindAllByURI(conn *pg.DB, imgs []string) ([]Image, error) {
	var imgList []Image

	if len(imgs) == 0 {
		return nil, fmt.Errorf("imgURIList cannot be empty")
	}

	// Build slice of img URI hashes.

	conn.Model(&imgList).Where("uri IN (?)", pg.In(imgs)).Select()

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
