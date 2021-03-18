package images

import (
	"database/sql"
	"errors"
	"net/url"
	"time"
)

// Image represents a pass/fail status of a previously queried img URI.
type Image struct {
	Hash      string         `json:"hash"`   // Sha254 hash of the base64 encoded contents of the image.
	URI       string         `json:"uri"`    // The original URI where the image resides on the web.
	Error     sql.NullString `json:"error"`  // Any error returned when trying to filter the image.
	Pass      bool           `json:"pass"`   // True if the image passed the filter.
	Reason    string         `json:"reason"` // Why the image did not pass the filter if non-empty e.g. "adult", "violence", etc.
	DateAdded time.Time      `json:"dateAdded"`
}

// NewImage returns a new Image instance (wow).
func NewImage(hash string, uri string, err error, pass bool, reason string, dateAdded time.Time) (*Image, error) {
	var sqlErr sql.NullString
	var validErr bool

	if _, _err := url.Parse(uri); _err != nil {
		return nil, err
	}

	if err == nil {
		validErr = false
		sqlErr = sql.NullString{String: "", Valid: validErr}
	} else {
		validErr = true
		sqlErr = sql.NullString{String: err.Error(), Valid: validErr}
	}

	// If an image fails the filter it can only have an empty reason property if and only if the err property is non-empty
	if !pass && err == nil && reason == "" {
		return nil, errors.New("an image cannot fail the filter and have both an empty error property AND reason property")
	}

	return &Image{
		Hash:      hash,
		URI:       uri,
		Error:     sqlErr,
		Pass:      pass,
		Reason:    reason,
		DateAdded: dateAdded,
	}, nil
}
