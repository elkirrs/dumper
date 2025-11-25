package validation

import (
	"reflect"
	"strings"
)

func fieldNameExtractor(f reflect.StructField) string {
	if tag := f.Tag.Get("yaml"); tag != "" {
		parts := strings.Split(tag, ",")
		return parts[0]
	}

	if tag := f.Tag.Get("yml"); tag != "" {
		parts := strings.Split(tag, ",")
		return parts[0]
	}

	return f.Name
}
