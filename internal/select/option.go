package _select

import (
	"dumper/internal/domain/config/database"
	dbConnect "dumper/internal/domain/config/db-connect"
	"dumper/internal/domain/config/server"
	"sort"
)

type DataOption interface {
	database.Database | server.Server
}

type pair struct {
	Display  string
	Original string
}

func SelectOptionList[T DataOption](
	options map[string]T,
	filter string,
) (map[string]string, []string) {
	pairs := make([]pair, 0, len(options))

	for idx, item := range options {
		var display string

		switch v := any(item).(type) {
		case database.Database:
			if filter != "" && v.Server != filter {
				continue
			}
			display = v.GetTitle()

		case server.Server:
			display = v.GetTitle()
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

func OptionDataBaseList(
	options map[string]dbConnect.DBConnect,
	filter string,
) (map[string]string, []string) {
	pairs := make([]pair, 0, len(options))

	for idx, dbConn := range options {
		var display string

		if filter != "" && dbConn.Database.Server != filter {
			continue
		}

		display = dbConn.Database.GetTitle()

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
