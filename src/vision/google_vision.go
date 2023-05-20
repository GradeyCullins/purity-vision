package vision

import (
	"context"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type BatchAnnotateResponse map[string]*pb.AnnotateImageResponse

func BatchGetImgAnnotation(uris []string) (*pb.BatchAnnotateImagesResponse, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	requests := make([]*pb.AnnotateImageRequest, 0)

	for _, uri := range uris {
		requests = append(requests, &pb.AnnotateImageRequest{
			Image: vision.NewImageFromURI(uri),
			Features: []*pb.Feature{
				{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: 5},
			},
		})
	}

	req := &pb.BatchAnnotateImagesRequest{Requests: requests}

	res, err := client.BatchAnnotateImages(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetImgAnnotation annotates an image.
func GetImgAnnotation(uri string) (*pb.AnnotateImageResponse, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// f, err := os.Open(filePath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer f.Close()

	// image, err := vision.NewImageFromReader(f)
	// if err != nil {
	// 	return nil, err
	// }

	image := vision.NewImageFromURI(uri)

	req := &pb.AnnotateImageRequest{
		Image: image,
		Features: []*pb.Feature{
			{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: 5},
		},
	}

	res, err := client.AnnotateImage(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func passSafeSearch(file string) error {
	// [START init]
	ctx := context.Background()

	// Create the client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	// [END init]

	// [START request]
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}

	res, err := client.AnnotateImage(ctx, &pb.AnnotateImageRequest{
		Image: image,
		Features: []*pb.Feature{
			{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: 5},
		},
	})
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

// findLabels gets labels from the Vision API for an image at the given file path.
func findLabels(file string) ([]string, error) {
	// [START init]
	ctx := context.Background()

	// Create the client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// [END init]

	// [START request]
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}

	// Perform the request.
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}
	// [END request]
	// [START transform]
	var labels []string
	for _, annotation := range annotations {
		labels = append(labels, annotation.Description)
	}
	return labels, nil
	// [END transform]
}
