package src

type filterReq struct {
	ImgURIList []string `json:"imgUriList"`
}

type filterRes struct {
	// ImgPassList is an array of booleans that stores a true or false
	// depending on whether the image of the given URI from the request
	// passes the criterion of being safe to view. The ordering mimics
	// that of the initial ImgURIList array of the request body.
	ImgPassList []bool `json:"imgPassList"`
}
