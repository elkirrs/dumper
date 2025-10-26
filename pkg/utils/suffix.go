package utils

import "strings"

func RemoveSuffix(str string, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return strings.TrimSuffix(str, suffix)
	}
	return str
}
