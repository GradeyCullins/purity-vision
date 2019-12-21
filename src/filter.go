package src

import (
	"fmt"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func filterImages(imgURIs []string) (*filterRes, error) {
	imgPassList := make([]bool, len(imgURIs))

	imgAnnotationsRes, err := getImgAnnotations(imgURIs)
	if err != nil {
		fmt.Println("err")
		return nil, err
	}

	for i, res := range imgAnnotationsRes.Responses {
		if res.Error != nil {
			return nil, fmt.Errorf("Failed to annotate image: %s with error: %s", imgURIs[i], res.Error)
		}
		imgPassList[i] = isImgSafe(res)
	}

	return &filterRes{imgPassList}, nil
}

func isImgSafe(air *pb.AnnotateImageResponse) bool {
	ssa := air.SafeSearchAnnotation
	if ssa.Adult > pb.Likelihood_POSSIBLE || ssa.Racy > pb.Likelihood_POSSIBLE {
		return false
	}

	return true
}
