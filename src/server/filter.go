package server

import (
	"fmt"

	"github.com/GradeyCullins/GoogleVisionFilter/src/vision"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// Images takes a list of image URIs and returns an HTTP payload
// with pass/fail status and any errors for each supplied URI.
func filter(imgURIList []string) (*BatchImgFilterRes, error) {
	cacheRes, err := checkImgCache(conn, imgURIList)
	if err != nil {
		return nil, err
	}

	imgFilterList := make([]ImgFilterRes, len(imgURIList)-len(cacheRes.ImgFilterResList))

	// Filter out URIs returned from the cache.
	for _, imgFilterRes := range cacheRes.ImgFilterResList {
		for i, uri := range imgURIList {
			// If the URI is found in the cache response, remove it from the URI list to avoid
			// redundant Google Vision API request.
			if imgFilterRes.ImgURI == uri {
				imgURIList[i] = imgURIList[len(imgURIList)-1]
				imgURIList = imgURIList[:len(imgURIList)-1]
			}
		}
	}

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
				Pass:   isImgSafe(res),
			}
		}
	}

	// Insert the new entries into the image cache.
	imgFilterList = append(imgFilterList, cacheRes.ImgFilterResList...)

	// Merge the cache response and the response.

	return &BatchImgFilterRes{ImgFilterResList: imgFilterList}, nil
}

func isImgSafe(air *pb.AnnotateImageResponse) bool {
	ssa := air.SafeSearchAnnotation

	// TODO: make the upper and lower bounds of this check configurable.
	if ssa.Adult >= pb.Likelihood_POSSIBLE || ssa.Racy >= pb.Likelihood_POSSIBLE {
		return false
	}

	return true
}
