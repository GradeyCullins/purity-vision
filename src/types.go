package src

type batchImgFilterReq struct {
	ImgURIList []string `json:"imgUriList"`
}

type imgFilterRes struct {
	Error string `json:"error"`
	Pass  bool   `json:"pass"`
}

type batchImgFilterRes struct {
	ImgFilterRes []imgFilterRes `json:"imgFilterRes"`
}
