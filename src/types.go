package src

type batchImgFilterReq struct {
	ImgURIList []string `json:"imgURIList"`
}

// TODO: consider adding 'URI' as a key. Currently, we are taking advantage of the fact
// that the JSON RFC guarantees array ordering. Comment in the below link suggests that
// there are known cases where order is not guaranteed.
//
// https://stackoverflow.com/a/7214312
type imgFilterRes struct {
	Error string `json:"error"`
	Pass  bool   `json:"pass"`
}

type batchImgFilterRes struct {
	ImgFilterRes []imgFilterRes `json:"imgFilterRes"`
}
