# 📦 Dumper

**Dumper** — This is a CLI utility for creating backups databases of various types (PostgreSQL, MySQL and etc.) with
flexible connection and storage settings.

---

![Demo](assets/dumper.gif)

## 🚀 Opportunities

- Multiple database systems can be managed.
- Support **PostgreSQL**, **MySQL**, **MongoDB** and etc.
- Connect for DB:
    - with the dump performed directly on the server and download (server).
- Work with **SSH-Keys** (include passphrase).
- Custom dump name templates.
- Archiving old dumps.
- Encrypting and Decrypting backup and config file
- Different formats.
- Different storages.

---

## 📄 Configuration

The configuration is set in a YAML file. (example, `config.yaml`).

---

### Configuration example:

```yaml
settings:
  db_port: "5432"
  driver: "psql"
  ssh:
    private_key: "path_your_key"
    passphrase: "your_passphrase"
    is_passphrase: true
  dir_remote: "/var/www/dump/"
  template: "{%srv%}_{%db%}_{%datetime%}"
  archive: true
  location: "server"
  format: "plain"
  dir_dump: "./dumps"
  dir_archived: "./archived"
  remove_dump: true
  parallel_download: 1
  storages:
    - local

storages:
  local:
    dir: "./dumps"

  ftp:
    dir: "./uploads/dumps"
    host: "116.34.17.94"
    port: 21
    username: "ftpuser"
    password: "123456"

  sftp:
    dir: "./dumps"
    host: "56.7.127.64"
    port: 22
    username: "sftpuser"
    private_key: "/Users/sftpuser/.ssh/id_rsa"
    passphrase: "123456"

servers:
  first-server:
    name: "test server"
    host: "10.234.23.45"
    port: "22"
    user: "user"
    password: "password"
  second-server:
    name: "mongo"
    host: "172.0.18.54"
    user: "root"
  remote-config-server:
    name: "mongo"
    host: "43.4.58.64"
    user: "user"
    port: 22
    conf_path: "/var/www/conf.yaml"

databases:
  demo:
    name: "db_name_demo"
    user: "demo_user"
    password: "password"
    server: "first-server"
    port: "5432"
    driver: "psql"
    format: "dump"

  mysql_db:
    name: "mysql_db_dumper"
    user: "user"
    password: "password"
    port: 3306
    driver: "mysql"
    format: "sql"
    server: "first-server"
    remove_dump: false

  mongo:
    name: "mongo_db_name"
    user: "root"
    password: "mongo_password"
    port: 27017
    driver: "mongo"
    format: "bson"
    server: "second-server"
    options:
      auth_source: "admin"
      ssl: true
    remove_dump: false

  maria_db:
    name: "maria_db_dumper"
    user: "user"
    password: "password"
    port: 3306
    driver: "mysql"
    format: "sql"
    server: "first-server"

  redis_db:
    name: "redis_db_dumper"
    user: "user"
    password: "password"
    port: 6379
    driver: "redis"
    format: "rdb"
    server: "first-server"
    options:
      mode: "sync"

```

The file example you can find in repository

---

### 📑 Configuration Description

#### The configuration consists of three sections:

#### 🔧 1. settings — global settings

Apply to all servers and databases, unless redefined locally.

| Parameter           | Description                                           | Type   | Rule     |
|---------------------|-------------------------------------------------------|--------|----------|
| `db_port`           | Default database connection port                      | int    | option   |
| `driver`            | [The DB driver list](#Driver)                         | string | required |
| `ssh.private_key`   | The path to the private SSH key.                      | string | option   |
| `ssh.passphrase`    | Passphrase for the key (optional).                    | string | option   |
| `ssh.is_passphrase` | whether to use passphrase from the config             | bool   | option   |
| `template`          | [File name template](#Template)                       | string | option   |
| `dir_remote`        | Dir remote for dumps                                  | string | option   |
| `archive`           | Archive of backup file                                | bool   | option   |
| `location`          | [Dump execution method](#Location)                    | string | required |
| `format`            | [The dump format](#Format)                            | string | required |
| `dir_dump`          | Directory for saving dumps                            | string | option   |
| `dir_archived`      | Archive Directory (need `{%srv%}_{%db%}` in template) | string | option   |
| `logging`           | Create logging                                        | bool   | option   |
| `retry_connect`     | attempts reconnect to server (default 5)              | int    | option   |
| `remove_dump`       | remove dump file after created (default true)         | bool   | option   |
| `encrypt.type`      | Type encrypting (only aes)                            | string | option   |
| `encrypt.password`  | Password for encrypting (only aes)                    | string | option   |
| `storages`          | Storage list when the dump need to upload             | list   | required |
| `parallel_download` | parallel upload dump file for several storages        | int    | option   |

#### Params:

- #### Driver:
    - `psql` — PostgreSQL
    - `mysql` — MySQL
    - `mongo` — MongoDB
    - `mariadb` — MariaDB
    - `mssql` — Microsoft SQL Server
    - `sqlite` — SQLite
    - `redis` — Redis

- #### Format:
    - PostgreSQL: `plain`, `dump`, `tar`
    - MySQL: `sql`
    - MongoDB: `bson`
    - MariaDB: `sql`
    - MSSQL: `bac`, `bacpac`
    - SQLite: `sql`
    - Redis: `rdb`

- #### Template:
    - `{%srv%}` — Name server
    - `{%db%}` — Name db
    - `{%datetime%}` — Date and time
    - `{%date%}` — Date
    - `{%time%}` — Time
    - `{%ts%}` — Time unix

- #### Location:
    - `server` — create dump in server and download

- #### Encrypt:
    - `aes` — type encrypting (openssl)

- #### Storages:
    - `local` - download local when the app started
    - `ftp` - upload to ftp server
    - `sftp` - upload to sftp server

#### 🖥 2. servers

Defines the connections through which databases can be backed up.

| Parameter   | Description                   | Type   | Rule                                       |
|-------------|-------------------------------|--------|--------------------------------------------|
| `title`     | Human-readable server name    | string | option                                     |
| `name`      | Server name                   | string | option                                     |
| `host`      | The IP address or domain name | string | required                                   |
| `port`      | Connection port               | int    | required<br/> (if not set `settings.port`) |
| `user`      | Username                      | string | required                                   |
| `password`  | Password (if there is no key) | string | required<br/> (if not set `key`)           |
| `key`       | Key (if there is no password) | string | required<br/> (if not set `password`)      |
| `conf_path` | Path remote config            | string | option (if set read only remote config)    | 

The configuration file on the remote `servers` must contain the servers and `databases` section.

#### 🗄 3. databases

A list of databases that need to be backed up.

| Parameter          | Description                                       | Type   | Rule     | Additional info                                  |
|--------------------|---------------------------------------------------|--------|----------|--------------------------------------------------|
| `title`            | Human-readable database name                      | string | option   |                                                  |
| `name`             | Database name (by default, the key name)          | string | option   | if set up driver sqlite need set up <path_to_db> |
| `user`             | The database user                                 | string | required |                                                  |
| `password`         | DB user's password                                | string | required |                                                  |
| `server`           | The link to the server from the `servers` section | string | required |                                                  |
| `port`             | Connection port                                   | int    | required | if not set `settings.db_port`                    |
| `driver`           | [The DB driver list](#Driver)                     | string | required | if not set `settings.driver`                     |
| `format`           | [The dump format](#Format)                        | string | required | if not set `settings.format`                     |
| `archive`          | Archiving a backup file                           | bool   | option   |                                                  |
| `options.*`        | Additional option for another databases           | object | option   | [Option list](#Options)                          |
| `remove_dump`      | remove dump file after created (default true)     | bool   | option   |                                                  |
| `encrypt.type`     | Type encrypting (only aes)                        | string | option   |                                                  |
| `encrypt.password` | Password for encrypting (only aes)                | string | option   |                                                  |
| `storages`         | Storage list when the dump need to upload         | list   | required |

#### Options

| Parameter             | Description            | Type   | Rule   | Additional info               |
|-----------------------|------------------------|--------|--------|-------------------------------|
| `options.auth_source` | Name database for auth | string | option | if set up driver mongo        |
| `options.ssl`         | SSL/TLS                | bool   | option | if set up driver mongo, mssql |
| `options.mode`        | Mode create dump       | string | option | if set up driver redis        |
| `options.role`        | Role for create dump   | string | option | if set up driver firebird     |
| `options.path`        | Path database SQLite   | string | option | if set up driver sqlite       |

---

#### 🔐 4. decrypt database file

If the file need to encrypt your database backup,
you can use encryption (the encryption utility must be
installed in the environment where the database backup
is performed)

In global  (encrypt for all databases)

```yaml
settings:
  encrypt:
    type: "aes"
    password: "123456"

servers:
# several servers
databases:
# several databases
```

Or only for a specific database

```yaml
settings:
# several settings
servers:
# several servers
databases:
  db-psql:
    name: 'mydb'
    user: 'myuser'
    password: 'mypassword'
    port: 5432
    driver: 'psql'
    server: 'srv-psql'
    format: 'plan'
    encrypt:
      type: "aes"
      password: "123456"
```

The file can be decrypted either via dumper or an encryption utility.

- Decrypt command

```
./dumper --crypt backup --mode decrypt --input ./dump.sql.gz.enc
```

or

```
openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in dump.sql.gz.enc -out dump.sql.gz -k 123456
```

- Encrypt command

```
./dumper --crypt backup --mode encrypt --input ./dump.sql.gz
```

#### 🔐 5. encrypt and decrypt file config

`
Encrypting a configuration file with one version will not be decrypted by another version of the application.
You can use the recovery token to decrypt it.
`

How it works:

1. Password encryption (outputs recovery token)
    ```
   ./dumper --crypt config --mode encrypt --input config.yaml
   ```
2. Launching the application (reading without password)
    ```
    ./dumper --config config.yaml
    ```
3. Decryption on the same device
    ```
    ./dumper --crypt config --mode decrypt --input config.yaml
    ```
4. Recovery on another device
    ```
    ./dumper --crypt config --mode recovery --token <recovery_token> --input config.yaml
    ```

#### 🗄️ 6. Storage

Configuration:

```yaml
storages:
  local:
    dir: "./dumps"

  ftp:
    dir: "./uploads/dumps"
    host: "172.168.139.109"
    port: 21
    username: "ftpuser"
    password: "123456"

  sftp:
    dir: "./dumps"
    host: "172.168.139.108"
    port: 22
    username: "sftpuser"
    private_key: "/Users/sftpuser/.ssh/id_rsa"
    passphrase: "123456" #option set up if key has passphrase
```

In global  (storage for all databases)

```yaml
settings:
  storages:
    - local

storages:
  local:
    dir: "./dumps"

  ftp:
    dir: "./uploads/dumps"
    host: "172.168.139.108"
    port: 21
    username: "ftpuser"
    password: "123456"

servers:
# several servers
databases:
# several databases
```

Or only for a specific database

```yaml
settings:
  storages:
    - local

storages:
  local:
    dir: "./dumps"

  ftp:
    dir: "./uploads/dumps"
    host: "172.168.139.108"
    port: 21
    username: "ftpuser"
    password: "123456"

servers:
# several servers

databases:
  db-psql:
    title: "My DB PSQL #1"
    name: "mydb"
    user: "myuser"
    password: "mypassword"
    port: 5432
    driver: "psql"
    server: "srv-psql"
    format: "plain"
    storages:
      - local
      - sftp 
```

### ▶ Launch examples

#### Backup with a choice of database from config file

```
./dumper --config ./cfg.yaml
````

#### Flag list

| Flag         | Example           | Description                                      | 
|--------------|-------------------|--------------------------------------------------|
| `--config`   | ./cfg.yaml        | path to config file                              |
| `--db`       | demo,app          | backup databases from list                       |
| `--all`      |                   | backup all databases from config file            |
| `--file-log` | file.log          | file name log file (if settings.logging == true) |
| `--crypt`    | config            | Crypt file: `backup`, `config`                   |                                 |
| `--input`    | ./dump.sql.gz.enc | path to encrypt file                             |
| `--mode`     | encrypt           | Mode crypt: `encrypt`, `decrypt`, `recovery`     |
| `--password` | 123456            | password for crypt (optional)                    |
| `--recovery` | 4j3k4lc7na09s     | Recovery token for recovery                      |

---

### 📂 Project structure

```
├── dumps/       # Directory for new dumps
├── archived/    # Archive of old dumps
├── config.yaml  # Configuration file
├── dumper       # The executable file of the utility
└── dumper.log   # Log file

```

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

## 💖 Support the Project

### If you like this project, you can support development:

- #### PayPal: [https://www.paypal.com/donate/?hosted_button_id=86QWWZSYNY4JN](https://www.paypal.com/donate/?hosted_button_id=86QWWZSYNY4JN)
- #### BTC [bc1qqrrtkymdck9q4h764hejjyenyfnyrpt4pgxd6h](bc1qqrrtkymdck9q4h764hejjyenyfnyrpt4pgxd6h)
- #### ETH [0xfe25171F3763E789d50279c2d4e16d2bAf14F701](0xfe25171F3763E789d50279c2d4e16d2bAf14F701)

## 🙏 Thank you for your support!