package template

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TemplateData struct {
	Time      time.Time
	Server    string
	Database  string
	Extension string
	Template  string
}

func GetTemplateFileName(data TemplateData) string {

	if data.Template == "" {
		data.Template = "{%srv%}_{%db%}_{%date%}"
	}

	if data.Time.IsZero() {
		data.Time = time.Now()
	}

	replacements := map[string]string{
		"{%srv%}":      data.Server,
		"{%db%}":       data.Database,
		"{%date%}":     data.Time.Format("2006.01.02"),
		"{%time%}":     data.Time.Format("15-04-05"),
		"{%datetime%}": data.Time.Format("2006.01.02_15-04-05"),
		"{%ts%}":       strconv.FormatInt(data.Time.Unix(), 10),
	}

	result := fmt.Sprintf("%s", data.Template)
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return strings.ReplaceAll(result, " ", "_")
}

func GetFullPath(dir string, file string) string {
	fullPath := fmt.Sprintf("%s/%s", dir, file)
	return filepath.Clean(fullPath)
}
