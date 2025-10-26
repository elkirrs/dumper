package _select

import (
	"dumper/internal/domain/config/database"
	dbConnect "dumper/internal/domain/config/db-connect"
	"dumper/internal/domain/config/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectOptionList_Database(t *testing.T) {
	options := map[string]database.Database{
		"db1": {Name: "DB One", Server: "srv1"},
		"db2": {Name: "DB Two", Server: "srv2"},
		"db3": {Name: "", Server: "srv1"},
	}

	result, keys := SelectOptionList[database.Database](options, "")
	assert.Equal(t, 3, len(result))
	assert.Contains(t, keys, "DB One")
	assert.Contains(t, keys, "DB Two")
	assert.Contains(t, keys, "db3")

	result, keys = SelectOptionList[database.Database](options, "srv1")
	assert.Equal(t, 2, len(result))
	assert.Contains(t, keys, "DB One")
	assert.Contains(t, keys, "db3")
	assert.NotContains(t, keys, "DB Two")
}

func TestSelectOptionList_Server(t *testing.T) {
	options := map[string]server.Server{
		"srv1": {Name: "Server A"},
		"srv2": {Name: ""},
	}

	result, keys := SelectOptionList[server.Server](options, "")
	assert.Equal(t, 2, len(result))
	assert.Contains(t, keys, "Server A")
	assert.Contains(t, keys, "srv2")
}

func TestOptionDataBaseList(t *testing.T) {
	options := map[string]dbConnect.DBConnect{
		"db1": {Database: database.Database{Name: "DB One", Server: "srv1"}},
		"db2": {Database: database.Database{Name: "DB Two", Server: "srv2"}},
		"db3": {Database: database.Database{Name: "", Server: "srv1"}},
	}

	result, keys := OptionDataBaseList(options, "")
	assert.Equal(t, 3, len(result))
	assert.Contains(t, keys, "DB One")
	assert.Contains(t, keys, "DB Two")
	assert.Contains(t, keys, "db3")

	result, keys = OptionDataBaseList(options, "srv1")
	assert.Equal(t, 2, len(result))
	assert.Contains(t, keys, "DB One")
	assert.Contains(t, keys, "db3")
	assert.NotContains(t, keys, "DB Two")
}

func TestSelectOptionList_Sorting(t *testing.T) {
	options := map[string]server.Server{
		"srv1": {Name: "C Server"},
		"srv2": {Name: "A Server"},
		"srv3": {Name: "B Server"},
	}

	_, keys := SelectOptionList[server.Server](options, "")
	assert.Equal(t, []string{"A Server", "B Server", "C Server"}, keys)
}
