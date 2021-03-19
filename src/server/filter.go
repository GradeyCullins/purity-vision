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
	imgResList := make(BatchImgFilterRes, 0)
	cachedImgFilterList := make([]ImgFilterRes, 0)

	// Images that are cached in DB.
	dbImgList, err := images.FindByURI(conn, imgURIList)
	if err != nil {
		return nil, err
	}

	size := len(imgURIList)

	// Populate map of uriHashes -> uris to be used in results later.
	for _, img := range dbImgList {
		for i := 0; i < size; i++ {
			uri := imgURIList[i]
			if img.URI == uri {
				cachedImgFilterList = append(cachedImgFilterList, ImgFilterRes{uri, img.Error.String, img.Pass, img.Reason})
				delete(uriPathMap, uri)

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

	for _, uri := range imgURIList {
		path, err := utils.Download(uri)

		// If the download fails, log the error and skip to the next download.
		if err != nil {
			logger.Error().Msg(err.Error())
			continue
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

	if len(imgURIList) > 0 {
		imgAnnotationsRes, err := vision.GetImgAnnotations(uriPathMap)
		if err != nil {
			return nil, err
		}

		for uri, res := range imgAnnotationsRes {
			if res.Error != nil {
				imgResList = append(imgResList, ImgFilterRes{
					ImgURI: uri,
					Error:  fmt.Sprintf("Failed to annotate image: %s with error: %s", uri, res.Error),
					Pass:   false,
					Reason: "",
				})
			} else {
				isPass, reason := isImgSafe(res, filterSettings)
				if reason != "" {
					logger.Debug().Msgf("Image %s failed the filter because it had %s content", uri, reason)
				}
				imgResList = append(imgResList, ImgFilterRes{
					ImgURI: uri,
					Error:  "",
					Pass:   isPass,
					Reason: reason,
				})
			}
		}
	} else {
		logger.Debug().Msg("No new images were sent to Vision API")
	}

	// Cache the new image filter response entries in the image table.
	for _, filterRes := range imgResList {
		img, err := images.NewImage(
			uriHashMap[filterRes.ImgURI],
			filterRes.ImgURI,
			errors.New(filterRes.Error),
			filterRes.Pass,
			filterRes.Reason,
			time.Now(),
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
	imgResList = append(imgResList, cachedImgFilterList...)

	return &imgResList, nil
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
