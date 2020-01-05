package src

import (
	"fmt"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func filterImages(imgURIs []string) (*batchImgFilterRes, error) {
	imgFilterList := make([]imgFilterRes, len(imgURIs))

	imgAnnotationsRes, err := getImgAnnotations(imgURIs)
	if err != nil {
		return nil, err
	}

	for i, res := range imgAnnotationsRes.Responses {
		if res.Error != nil {
			imgFilterList[i] = imgFilterRes{
				fmt.Sprintf("Failed to annotate image: %s with error: %s", imgURIs[i], res.Error),
				false,
			}
		} else {
			imgFilterList[i] = imgFilterRes{
				"",
				isImgSafe(res),
			}
		}
	}

	return &batchImgFilterRes{imgFilterList}, nil
}

func isImgSafe(air *pb.AnnotateImageResponse) bool {
	ssa := air.SafeSearchAnnotation

	// TODO: make the upper and lower bounds of this check configurable.
	if ssa.Adult > pb.Likelihood_POSSIBLE || ssa.Racy > pb.Likelihood_POSSIBLE {
		return false
	}

	return true
}
