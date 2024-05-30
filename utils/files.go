package utils

import "strings"

func GetFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 0 {
		return parts[1]
	}
	return ""
}

func GetFileName(fileName string) string {
	return strings.Split(fileName, ".")[0]
}
