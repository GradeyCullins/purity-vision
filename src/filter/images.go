package filter

import (
	"database/sql"
	"fmt"

	"github.com/GradeyCullins/GoogleVisionFilter/src/cache"
	"github.com/GradeyCullins/GoogleVisionFilter/src/types"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// Images takes a list of image URIs and returns an HTTP payload
// with pass/fail status and any errors for each supplied URI.
func Images(db *sql.DB, imgURIs []string) (*types.BatchImgFilterRes, error) {
	imgFilterList := make([]types.ImgFilterRes, len(imgURIs))

	// TODO: implement cache for avoiding redundant Vision API calls.
	_, err := cache.CheckImgCache(db, imgURIs)
	if err != nil {
		return nil, err
	}

	imgAnnotationsRes, err := getImgAnnotations(imgURIs)
	if err != nil {
		return nil, err
	}

	for i, res := range imgAnnotationsRes.Responses {
		if res.Error != nil {
			imgFilterList[i] = types.ImgFilterRes{
				Error: fmt.Sprintf("Failed to annotate image: %s with error: %s", imgURIs[i], res.Error),
				Pass:  false,
			}
		} else {
			imgFilterList[i] = types.ImgFilterRes{
				Error: "",
				Pass:  isImgSafe(res),
			}
		}
	}

	return &types.BatchImgFilterRes{ImgFilterRes: imgFilterList}, nil
}

func isImgSafe(air *pb.AnnotateImageResponse) bool {
	ssa := air.SafeSearchAnnotation

	// TODO: make the upper and lower bounds of this check configurable.
	if ssa.Adult > pb.Likelihood_POSSIBLE || ssa.Racy > pb.Likelihood_POSSIBLE {
		return false
	}

	return true
}
