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

// filter takes a list of image URIs and returns a response
// with pass/fail statuses and any errors for each supplied URI.
// func filter(imgURIList []string) (*BatchImgFilterRes, error) {
func filter(filterRequest BatchImgFilterReq) (*BatchImgFilterRes, error) {
	//TODO: cleanup this function

	// TODO: if filterRequest.ImgURIList is longer than 16, page the request.
	// num pages = len(imgURIList) / 16 rounded up to the nearest int
	// pageCount := int(math.Ceil(float64(len(filterRequest.ImgURIList)) / MAX_IMAGES_PER_REQUEST))
	// for pageNum := 0; pageNum < pageCount; pageNum++ {
	// }

	imgURIList := filterRequest.ImgURIList
	uriPathMap := make(map[string]string, len(imgURIList)) // Map URI to image file path on disk to be fed to Vision API.
	uriHashMap := make(map[string]string, len(imgURIList))
	filterSettings := filterRequest.FilterSettings
	cachedImgFilterList := make([]ImgFilterRes, 0)

	for _, uri := range imgURIList {
		path, err := utils.Download(uri)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r := bufio.NewReader(f)
		content, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		uriHashMap[uri] = utils.Hash(base64.StdEncoding.EncodeToString(content))
		uriPathMap[uri] = f.Name()
	}

	// Delete image files.
	defer func() {
		for _, path := range uriPathMap {
			os.Remove(path)
		}
	}()

	// Images that are cached in DB.
	dbImgList, err := images.FindByURI(conn, imgURIList)
	if err != nil {
		return nil, err
	}

	imgFilterList := make(BatchImgFilterRes, len(imgURIList)-len(dbImgList))
	size := len(imgURIList)

	// Populate map of uriHashes -> uris to be used in results later.
	for _, img := range dbImgList {
		for i := 0; i < size; i++ {
			uri := imgURIList[i]
			if img.URI == uri {
				cachedImgFilterList = append(cachedImgFilterList, ImgFilterRes{uri, img.Error.String, img.Pass})

				// If the URI is found in the cache response, remove it from the URI list to avoid
				// redundant Google Vision API request.
				imgURIList, err = utils.StringSliceRemove(imgURIList, i)
				if err != nil {
					return nil, err
				}
				size -= 1
			}
		}
	}

	if len(imgURIList) > 0 {
		imgAnnotationsRes, err := vision.GetImgAnnotations(uriPathMap)
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
		logger.Debug().Msg("No new images were sent to Vision API")
	}

	// Cache the new image filter response entries in the image table.
	for _, filterRes := range imgFilterList {
		img, err := images.NewImage(
			uriHashMap[filterRes.ImgURI],
			filterRes.ImgURI,
			errors.New(filterRes.Error),
			filterRes.Pass, time.Now(),
		)
		if err != nil {
			return nil, err
		}
		if err = images.Insert(conn, img); err != nil {
			return nil, err
		}
		logger.Debug().Msgf("Adding %s to DB cache", filterRes.ImgURI)
	}

	// Merge the cache response and the response.
	imgFilterList = append(imgFilterList, cachedImgFilterList...)

	return &imgFilterList, nil
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
