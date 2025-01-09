package util

import "strings"

func IsOfficialImage(imageName string) bool {
	return !strings.Contains(imageName, "/")
}
