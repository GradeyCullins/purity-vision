package server

import (
	"errors"
	"fmt"
	"purity-vision-filter/src/images"
	"purity-vision-filter/src/utils"
	"purity-vision-filter/src/vision"
	"time"

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

// filter takes a list of image URIs and returns a response
// with pass/fail statuses and any errors for each supplied URI.
// func filter(imgURIList []string) (*BatchImgFilterRes, error) {
func filter(filterRequest BatchImgFilterReq) (*BatchImgFilterRes, error) {
	imgURIList := filterRequest.ImgURIList
	filterSettings := filterRequest.FilterSettings
	cachedImgFilterList := make([]ImgFilterRes, 0)

	// Images that are cached in DB.
	dbImgList, err := images.FindByURI(conn, imgURIList)
	if err != nil {
		return nil, err
	}

	imgFilterList := make([]ImgFilterRes, len(imgURIList)-len(dbImgList))

	// Populate map of uriHashes -> uris to be used in results later.
	for _, img := range dbImgList {
		for i, uri := range imgURIList {
			if img.ImgURIHash == utils.Hash(uri) {
				cachedImgFilterList = append(cachedImgFilterList, ImgFilterRes{uri, img.Error.String, img.Pass})

				// If the URI is found in the cache response, remove it from the URI list to avoid
				// redundant Google Vision API request.
				imgURIList, err = utils.StringSliceRemove(imgURIList, i)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if len(imgURIList) > 0 {
		imgAnnotationsRes, err := vision.GetImgAnnotations(imgURIList)
		if err != nil {
			return nil, err
		}

		for i, res := range imgAnnotationsRes.Responses {
			if res.Error != nil {
				imgFilterList[i] = ImgFilterRes{
					ImgURI: imgURIList[i],
					Error:  fmt.Sprintf("Failed to annotate image: %s with error: %s", imgURIList[i], res.Error),
					Pass:   false,
				}
			} else {
				imgFilterList[i] = ImgFilterRes{
					ImgURI: imgURIList[i],
					Error:  "",
					Pass:   isImgSafe(res, &filterSettings),
				}
			}
		}
	} else {
		logger.Info().Msg("No new images were sent to Vision API")
	}

	// Cache the new image filter response entries in the image table.
	for _, filterRes := range imgFilterList {
		img, err := images.NewImage(filterRes.ImgURI, errors.New(filterRes.Error), filterRes.Pass, time.Now())
		if err != nil {
			return nil, err
		}
		if err = images.Insert(conn, img); err != nil {
			return nil, err
		}
		logger.Info().Msgf("Adding %s to DB cache", filterRes.ImgURI)
	}

	// Merge the cache response and the response.
	imgFilterList = append(imgFilterList, cachedImgFilterList...)

	return &BatchImgFilterRes{ImgFilterResList: imgFilterList}, nil
}

func isImgSafe(air *pb.AnnotateImageResponse, fs *FilterSettings) bool {
	ssa := air.SafeSearchAnnotation

	// Enable all the filter rules if the filter settings are empty.
	if fs.IsEmpty() {
		fs = &FilterSettings{
			Adult:    true,
			Violence: true,
			Racy:     true,
			Medical:  true,
		}
	}

	if fs.Adult && ssa.Adult >= pb.Likelihood_LIKELY {
		return false
	}

	if fs.Medical && ssa.Medical >= pb.Likelihood_LIKELY {
		return false
	}

	if fs.Violence && ssa.Violence >= pb.Likelihood_LIKELY {
		return false
	}

	if fs.Racy && ssa.Racy >= pb.Likelihood_LIKELY {
		return false
	}

	return true
}
