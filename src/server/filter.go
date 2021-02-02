package server

import (
	"errors"
	"fmt"
	"google-vision-filter/src/images"
	"google-vision-filter/src/utils"
	"os"
	"time"

	"github.com/GradeyCullins/GoogleVisionFilter/src/vision"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

// filter takes a list of image URIs and returns a response
// with pass/fail statuses and any errors for each supplied URI.
func filter(imgURIList []string) (*BatchImgFilterRes, error) {
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
					Pass: isImgSafe(res, &FilterSettings{
						Adult:    filterCategory{Enabled: true, Likelihood: pb.Likelihood_LIKELY},
						Medical:  filterCategory{Enabled: true, Likelihood: pb.Likelihood_LIKELY},
						Violence: filterCategory{Enabled: true, Likelihood: pb.Likelihood_LIKELY},
						Racy:     filterCategory{Enabled: true, Likelihood: pb.Likelihood_LIKELY},
					}),
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

type filterCategory struct {
	Enabled    bool          `json:"enabled"`
	Likelihood pb.Likelihood `json:"likelihood"`
}

// FilterSettings represents user-configurable image filter settings.
type FilterSettings struct {
	Adult    filterCategory `json:"adult"`
	Medical  filterCategory `json:"medical"`
	Violence filterCategory `json:"violence"`
	Racy     filterCategory `json:"racy"`
}

func isImgSafe(air *pb.AnnotateImageResponse, fs *FilterSettings) bool {
	ssa := air.SafeSearchAnnotation

	if fs.Adult.Enabled && ssa.Adult >= fs.Adult.Likelihood {
		return false
	}

	if fs.Medical.Enabled && ssa.Medical >= fs.Medical.Likelihood {
		return false
	}

	if fs.Violence.Enabled && ssa.Violence >= fs.Violence.Likelihood {
		return false
	}

	if fs.Racy.Enabled && ssa.Racy >= fs.Racy.Likelihood {
		return false
	}

	// // TODO: make the upper and lower bounds of this check configurable.
	// if ssa.Adult >= pb.Likelihood_POSSIBLE || ssa.Racy >= pb.Likelihood_POSSIBLE {
	// 	return false
	// }

	return true
}
