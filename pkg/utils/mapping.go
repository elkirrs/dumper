package utils

func GetDBSource(driverName, format string) string {
	switch driverName {
	case "mariadb":
		return "mariadb-dump"
	case "mongodb":
		return "mongodump"
	case "mssql":
		if format == "bacpac" {
			return "sqlpackage"
		}
		return "sqlcmd"
	case "mysql":
		return "mysqldump"
	case "psql":
		return "pg_dump"
	case "redis":
		return "redis-cli"
	case "sqlite":
		return "sqlite3"
	default:
		return ""
	}
}
