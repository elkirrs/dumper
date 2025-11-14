# Changelog
All notable changes to this project will be documented here.

## [1.8.3] - 2025-11-15
- Fix bug with create dump mongodb
- Add parameter archive in database section
- Update file name template

## [1.8.2] - 2025-11-09
- Fix static path main db dump to global
- Feature added remote dir for dump

## [1.8.1] - 2025-11-07
- Fix issue in crypt recovery mode
- Feature set password in command line

## [1.8.0] - 2025-11-02
- Add support storage FTP, SFTP
- Multiple upload dump file by several servers

## [1.7.0] - 2025-10-31
- Add encrypt config file

## [1.6.0] - 2025-10-28
- Add support encrypt aes (openssl) for dump files 
- Add Server and Database titles config

## [1.5.0] - 2025-10-26
- Add support MSSQL, SQLite, Redis
- Fix bags

## [1.4.0] - 2025-10-18
- Add support remote config
- Update prepare Database and Server
- Add view error with create dump
- Add mask user and password for log and error

## [1.3.0] - 2025-08-28
- Add driver MariaDB
- Add config remove file dump after from server after created
- View file size before downloading

## [1.2.0] - 2025-08-21
- Add driver MongoDB
- Add driver MySQL

## [1.1.1] - 2025-08-20
- Add retry connect to server
- Add flag --all. Using this flag, you can dump all databases from the configuration.

## [1.1.0] - 2025-08-19
- Add flag --db. It can be used to create backups of multiple databases in one run
- Fix the exit from the app

## [1.0.0] - 2025-08-17
### Added
- Initial public release
- Creating a backup of the postgresql database with connection to the server and downloading
- Connecting to the server. Creating a database dump. Archiving the database and downloading from the server. Deleting a created one.