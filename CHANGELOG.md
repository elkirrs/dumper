# Changelog
All notable changes to this project will be documented here.

## [Unreleased]
- Planned features
- Support different databases like Mysql, MongoDB and etc.
- Support connect in server and database

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