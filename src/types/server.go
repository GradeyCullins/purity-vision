package types

// BatchImgFilterReq is the form of an incoming JSON payload
// for retrieving pass/fail status of each supplied image URI.
type BatchImgFilterReq struct {
	ImgURIList []string `json:"imgURIList"`
}

// ImgFilterRes returns the pass/fail status and any errors for a single image URI.
//
// TODO: consider adding 'URI' as a key. Currently, we are taking advantage of the fact
// that the JSON RFC guarantees array ordering. Comment in the below link suggests that
// there are known cases where order is not guaranteed.
// https://stackoverflow.com/a/7214312
type ImgFilterRes struct {
	Error string `json:"error"`
	Pass  bool   `json:"pass"`
}

// BatchImgFilterRes represents a list of pass/fail statuses and any errors for each
// supplied image URI.
type BatchImgFilterRes struct {
	ImgFilterRes []ImgFilterRes `json:"imgFilterRes"`
}
