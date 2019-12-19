// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command label uses the Vision API's label detection capabilities to find a label based on an image's content.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	// [START imports]
	"context"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	// [END imports]
)

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
	fmt.Println(res)
	return nil
}

func passSafeSearchList(files []string) error {
	// [START init]
	ctx := context.Background()

	// Create the client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	var annotateReqs []*pb.AnnotateImageRequest

	// Loop over the files and create annotate requests.
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		image, err := vision.NewImageFromReader(f)
		if err != nil {
			return err
		}
		annotateReqs = append(annotateReqs, &pb.AnnotateImageRequest{
			Image: image,
			Features: []*pb.Feature{
				{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: 5},
			},
		})

	}
	// [END init]
	req := &pb.AsyncBatchAnnotateImagesRequest{
		Requests: annotateReqs
		// TODO: Fill request struct fields.
	}
	op, err := c.AsyncBatchAnnotateImages(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	
	resp, err := op.Wait(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp


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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	err := passSafeSearch(args[0])
	if err != nil {
		fmt.Println(err)
	}
	//labels, err := findLabels(args[0])
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	os.Exit(1)
	//}
	//if len(labels) == 0 {
	//	fmt.Println("No labels found.")
	//} else {
	//	fmt.Println("Found labels:")
	//	for _, label := range labels {
	//		fmt.Println(label)
	//	}
	//}
}
