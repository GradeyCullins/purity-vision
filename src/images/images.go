package images

import (
	"database/sql"
	"net/url"
	"purity-vision-filter/src/utils"
	"time"
)

// Image represents a pass/fail status of a previously queried img URI.
type Image struct {
	ImgURIHash string         `json:"img_uri_hash"`
	Error      sql.NullString `json:"error"`
	Pass       bool           `json:"pass"`
	DateAdded  time.Time      `json:"dateAdded"`
}

// NewImage returns a new Image instance (wow).
func NewImage(uri string, err error, pass bool, dateAdded time.Time) (*Image, error) {
	var sqlErr sql.NullString
	var validErr bool

	if _, err := url.Parse(uri); err != nil {
		return nil, err
	}

	if err.Error() == "" {
		validErr = false
	} else {
		validErr = true
	}
	sqlErr = sql.NullString{String: err.Error(), Valid: validErr}

	return &Image{
		ImgURIHash: utils.Hash(uri),
		Error:      sqlErr,
		Pass:       pass,
		DateAdded:  dateAdded,
	}, nil
}
