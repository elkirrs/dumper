package mapping

type DriverInfo struct {
	DefaultCommand string
	DefaultPort    string
	Formats        map[string]struct{}
	Overrides      map[string]string
}

var dbDrivers = map[string]DriverInfo{
	"mariadb": {
		DefaultCommand: "mariadb-dump",
		DefaultPort:    "3306",
		Formats:        map[string]struct{}{"sql": {}},
	},
	"mongo": {
		DefaultCommand: "mongodump",
		DefaultPort:    "27015",
		Formats:        map[string]struct{}{"bson": {}, "archive": {}},
	},
	"mssql": {
		DefaultCommand: "sqlcmd",
		DefaultPort:    "27015",
		Formats:        map[string]struct{}{"bac": {}, "bacpac": {}},
		Overrides:      map[string]string{"bacpac": "sqlpackage"},
	},
	"mysql": {
		DefaultCommand: "mysqldump",
		DefaultPort:    "3306",
		Formats:        map[string]struct{}{"sql": {}},
	},
	"psql": {
		DefaultCommand: "pg_dump",
		DefaultPort:    "5432",
		Formats:        map[string]struct{}{"plain": {}, "dump": {}, "tar": {}},
	},
	"redis": {
		DefaultCommand: "redis-cli",
		DefaultPort:    "6379",
		Formats:        map[string]struct{}{"rdb": {}},
	},
	"sqlite": {
		DefaultCommand: "sqlite3",
		DefaultPort:    "",
		Formats:        map[string]struct{}{"sql": {}},
	},
	"neo4j": {
		DefaultCommand: "neo4j-admin",
		DefaultPort:    "",
		Formats:        map[string]struct{}{"dump": {}},
	},
	"dynamodb": {
		DefaultCommand: "aws",
		DefaultPort:    "",
		Formats:        map[string]struct{}{"json": {}},
	},
	"influxdb": {
		DefaultCommand: "influx",
		DefaultPort:    "8086",
		Formats:        map[string]struct{}{"tar": {}},
	},
	"db2": {
		DefaultCommand: "db2",
		DefaultPort:    "50000",
		Formats:        map[string]struct{}{"0.db2": {}},
	},
}

func IsValidFormatDump(driverName, format string) bool {
	if driver, ok := dbDrivers[driverName]; ok {
		_, exists := driver.Formats[format]
		return exists
	}
	return false
}

func GetDBSource(driverName, format string) string {
	driver, ok := dbDrivers[driverName]
	if !ok {
		return ""
	}

	if driver.Overrides != nil {
		if cmd, exists := driver.Overrides[format]; exists {
			return cmd
		}
	}

	return driver.DefaultCommand
}

func GetDefaultDBPort(driverName string) string {
	driver, ok := dbDrivers[driverName]
	if !ok {
		return ""
	}
	return driver.DefaultPort
}
