package template

import (
	"strconv"
	"testing"
	"time"
)

func TestGetTemplateFileName(t *testing.T) {
	fixedTime := time.Date(2025, 10, 31, 12, 34, 5, 0, time.UTC)

	tests := []struct {
		name string
		data TemplateData
		want string
	}{
		{
			name: "default template used when empty",
			data: TemplateData{
				Time:     fixedTime,
				Server:   "web",
				Database: "site",
			},
			want: "web_site_2025.10.31",
		},
		{
			name: "custom template with date and time",
			data: TemplateData{
				Time:     fixedTime,
				Server:   "web",
				Database: "site",
				Template: "{%srv%}_{%db%}_{%date%}_{%time%}",
			},
			want: "web_site_2025.10.31_12-34-05",
		},
		{
			name: "custom template with timestamp",
			data: TemplateData{
				Time:     fixedTime,
				Server:   "web",
				Database: "site",
				Template: "{%srv%}_{%db%}_{%ts%}",
			},
			want: "web_site_1761921245",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTemplateFileName(tt.data)

			if tt.name == "custom template with timestamp" {
				tt.want = "web_site_" + strconv.FormatInt(fixedTime.Unix(), 10)
			}

			if got != tt.want {
				t.Errorf("GetTemplateFileName() = %q; want %q", got, tt.want)
			}
		})
	}
}
