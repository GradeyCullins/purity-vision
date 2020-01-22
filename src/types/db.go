package types

// User represents a user in the Purity system.
type User struct {
	UID      int    `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ImageCacheEntry represents a pass/fail status of a previously queried img URI.
type ImageCacheEntry struct {
	ImgURIHash string `json:"imgURIHash"`
	Pass       bool   `json:"pass"`
}
