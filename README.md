<p align="center"><a href="https://elkirrs.github.io/dumper" target="_blank"><img src="assets/dumper.svg" width="400" alt="Laravel Logo"></a></p>

## About Dumper

This is a CLI utility for creating backups databases of various types with flexible connection and storage settings.

Opportunities:

- Multiple database systems can be managed.
- Support **PostgreSQL**, **MySQL**, **MongoDB** and etc.
- Work with **SSH-Keys** (include passphrase).
- Custom dump name templates.
- Archiving old dumps.
- Encrypting and Decrypting backup and config file
- Support different storages **SFTP**, **Azure** and etc.
- Backup from docker
- Shell script after and before backup

## Required Tools

Dumper requires the following database tools to be installed and available on the host where the database dump is being performed:

| Database   | Required Tool    | Installation Example                          |
|------------|------------------|-----------------------------------------------|
| PostgreSQL | `pg_dump`        | `apt-get install postgresql-client`           |
| MySQL      | `mysqldump`      | `apt-get install mysql-client`                |
| MariaDB    | `mariadb-dump`   | `apt-get install mariadb-client`              |
| MongoDB    | `mongodump`      | `apt-get install mongodb-database-tools`      |
| SQLite     | `sqlite3`        | `apt-get install sqlite3`                     |
| Redis      | `redis-cli`      | `apt-get install redis-tools`                 |
| MSSQL      | `sqlcmd` or `sqlpackage` | Download from Microsoft                   |
| Neo4j      | `neo4j-admin`    | Included with Neo4j installation              |
| DynamoDB   | `aws`            | `pip install awscli` or `apt-get install awscli` |

**Note:** These tools must be installed on the server/host where the database resides, not necessarily where the dumper CLI is running (unless dumping locally).

### Installation:
```shell
    curl -sSL https://elkirrs.github.io/dumper/install.sh | sh 
```

### Documentation:
Click to docs page: [Documentation](https://elkirrs.github.io/dumper)


### Checking the receipt by the user:

- Mac/Linux:

```
    shasum -a 256 dumper_linux_amd64.tar.gz
    cat checksums.txt | grep dumper_linux_amd64.tar.gz
```

- Windows (PowerShell):

```
    Get-FileHash .\dumper_windows_amd64.zip -Algorithm SHA256
```

## Support the Project

#### If you like this project, you can support development:

- #### PayPal: [https://www.paypal.com/donate/?hosted_button_id=86QWWZSYNY4JN](https://www.paypal.com/donate/?hosted_button_id=86QWWZSYNY4JN)
- #### BTC [bc1qqrrtkymdck9q4h764hejjyenyfnyrpt4pgxd6h](bc1qqrrtkymdck9q4h764hejjyenyfnyrpt4pgxd6h)
- #### ETH [0xfe25171F3763E789d50279c2d4e16d2bAf14F701](0xfe25171F3763E789d50279c2d4e16d2bAf14F701)

### Thank you for your support!