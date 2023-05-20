package utils

import "github.com/google/uuid"

func GenerateLicenseKey() string {
	key := uuid.New()
	return key.String()
}
