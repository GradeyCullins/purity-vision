package server

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"purity-vision-filter/src/images"
	"purity-vision-filter/src/utils"
	"purity-vision-filter/src/vision"
	"time"

	"github.com/rs/zerolog/log"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// FilterSettings represents user-configurable image filter settings.
type FilterSettings struct {
	Adult      bool          `json:"adult"`
	Medical    bool          `json:"medical"`
	Violence   bool          `json:"violence"`
	Racy       bool          `json:"racy"`
	Likelihood pb.Likelihood `json:"likelihood"` // Not used for now.
}

func (fs *FilterSettings) IsEmpty() bool {
	return !fs.Adult && !fs.Medical && !fs.Violence && !fs.Racy
}

// const MAX_IMAGES_PER_REQUEST = 16

func filterSingle(uri string, fs FilterSettings) (*ImgFilterRes, error) {
	var res *ImgFilterRes

	// check if URI cached in db
	img, err := images.FindByURI(conn, uri)
	if err != nil {
		return nil, err
	}

	// if cached, return cached response
	if img != nil {
		log.Debug().Msgf("Found cached image %s...", uri)
		return &ImgFilterRes{
			ImgURI: uri,
			Error:  img.Error.String,
			Pass:   img.Pass,
			Reason: img.Reason,
		}, nil
	}

	// otherwise download file
	path, hash, err := fetchAndReadFile(uri)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)

	// Make filter request to Google Vision API.
	imgAnnotation, err := vision.GetImgAnnotation(path)
	if err != nil {
		return nil, err
	}

	if imgAnnotation.Error != nil {
		res = &ImgFilterRes{
			ImgURI: uri,
			Error:  fmt.Sprintf("Failed to annotate image: %s with error: %s", uri, imgAnnotation.Error),
			Pass:   false,
			Reason: "",
		}
	} else {
		isPass, reason := isImgSafe(imgAnnotation, fs)
		if reason != "" {
			logger.Debug().Msgf("Image %s failed the filter because it had %s content", uri, reason)
		}
		res = &ImgFilterRes{
			ImgURI: uri,
			Error:  "",
			Pass:   isPass,
			Reason: reason,
		}
	}

	// Cache filter response to db.
	img, err = images.NewImage(
		hash,
		uri,
		errors.New(res.Error),
		res.Pass,
		res.Reason,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	if err = images.Insert(conn, img); err != nil {
		return nil, err
	} else {
		logger.Debug().Msgf("Adding %s to DB cache", res.ImgURI)
	}

	return res, nil
}

func fetchAndReadFile(uri string) (string, string, error) {
	path, err := utils.Download(uri)

	// If the download fails, log the error and skip to the next download.
	if err != nil {
		return "", "", err
	}
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return "", "", err
	}

	return path, utils.Hash(base64.StdEncoding.EncodeToString(content)), nil
}

func isImgSafe(air *pb.AnnotateImageResponse, fs FilterSettings) (bool, string) {
	ssa := air.SafeSearchAnnotation

	// Enable all the filter rules if the filter settings are empty.
	if fs.IsEmpty() {
		fs = FilterSettings{
			Adult:    true,
			Violence: true,
			Racy:     true,
			Medical:  true,
		}
	}

	if fs.Adult && ssa.Adult >= pb.Likelihood_LIKELY {
		return false, "adult"
	}

	if fs.Medical && ssa.Medical >= pb.Likelihood_LIKELY {
		return false, "medical"
	}

	if fs.Violence && ssa.Violence >= pb.Likelihood_LIKELY {
		return false, "violence"
	}

	// Set this to "very likely" because anything less strict
	// was filtering out cartoons that I did not personally deem
	// to be innapropriate or "racy".
	if fs.Racy && ssa.Racy >= pb.Likelihood_VERY_LIKELY {
		return false, "racy"
	}

	return true, ""
}
