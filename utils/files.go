package utils

import "strings"

func GetFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 0 {
		return parts[1]
	}
	return ""
}
