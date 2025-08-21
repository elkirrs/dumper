package _select

import (
	"dumper/internal/config"
	"sort"
)

type DataOption interface {
	config.Database | config.Server
}

type pair struct {
	Display  string
	Original string
}

func SelectOptionList[T DataOption](options map[string]T, filter string) (map[string]string, []string) {
	pairs := make([]pair, 0, len(options))

	for idx, item := range options {
		var display string

		switch v := any(item).(type) {
		case config.Database:
			if filter != "" && v.Server != filter {
				continue
			}
			display = v.Name

		case config.Server:
			display = v.Name
		default:
			continue
		}

		if display == "" {
			display = idx
		}

		pairs = append(pairs, pair{Display: display, Original: idx})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Display < pairs[j].Display
	})

	result := make(map[string]string, len(pairs))
	keys := make([]string, 0, len(pairs))
	for _, p := range pairs {
		result[p.Display] = p.Original
		keys = append(keys, p.Display)
	}

	return result, keys
}
