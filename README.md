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
- Encrypting:
    - Encrypting dump file
    - Encrypt and Decrypt config file
- Different formats.

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
  template: "{%srv%}_{%db%}_{%datetime%}"
  archive: true
  location: "server"
  format: "plain"
  dir_dump: "./dumps"
  dir_archived: "./archived"
  remove_dump: true

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

| Parameter           | Description                                              | type     |
|---------------------|----------------------------------------------------------|----------|
| `db_port`           | Default database connection port                         | option   |
| `driver`            | [The DB driver list](#Driver)                            | required |
| `ssh.private_key`   | The path to the private SSH key.                         | option   |
| `ssh.passphrase`    | Passphrase for the key (optional).                       | option   |
| `ssh.is_passphrase` | whether to use passphrase from the config                | option   |
| `template`          | [File name template](#Template)                          | option   |
| `archive`           | Archiving old dumps (need `{%srv%}_{%db%}` in template). | option   |
| `location`          | [Dump execution method](#Location)                       | required |
| `format`            | [The dump format](#Format)                               | required |
| `dir_dump`          | Directory for saving dumps                               | option   |
| `dir_archived`      | Archive Directory                                        | option   |
| `logging`           | Create logging                                           | option   |
| `retry_connect`     | attempts reconnect to server (default 5)                 | option   |
| `remove_dump`       | remove dump file after created (default true)            | option   |
| `encrypt.type`      | Type encrypting (only aes)                               | option   |
| `encrypt.password`  | Password for encrypting (only aes)                       | option   |

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

#### 🖥 2. servers

Defines the connections through which databases can be backed up.

| Parameter   | Description                   | type                                       |
|-------------|-------------------------------|--------------------------------------------|
| `title`     | Human-readable server name    | option                                     |
| `name`      | Server name                   | option                                     |
| `host`      | The IP address or domain name | required                                   |
| `port`      | Connection port               | required<br/> (if not set `settings.port`) |
| `user`      | Username                      | required                                   |
| `password`  | Password (if there is no key) | required<br/> (if not set `key`)           |
| `key`       | Key (if there is no password) | required<br/> (if not set `password`)      |
| `conf_path` | Path remote config            | option (if set read only remote config)    | 

The configuration file on the remote `servers` must contain the servers and `databases` section.

#### 🗄 3. databases

A list of databases that need to be backed up.

| Parameter          | Description                                       | Type     | Additional info                                  |
|--------------------|---------------------------------------------------|----------|--------------------------------------------------|
| `title`            | Human-readable database name                      | option   |                                                  |
| `name`             | Database name (by default, the key name)          | option   | if set up driver sqlite need set up <path_to_db> |
| `user`             | The database user                                 | required |                                                  |
| `password`         | DB user's password                                | required |                                                  |
| `server`           | The link to the server from the `servers` section | required |                                                  |
| `port`             | Connection port                                   | required | if not set `settings.db_port`                    |
| `driver`           | [The DB driver list](#Driver)                     | required | if not set `settings.driver`                     |
| `format`           | [The dump format](#Format)                        | required | if not set `settings.format`                     |
| `options.*`        | Additional option for another databases           | option   | [Option list](#Options)                          |
| `remove_dump`      | remove dump file after created (default true)     | option   |                                                  |
| `encrypt.type`     | Type encrypting (only aes)                        | option   |                                                  |
| `encrypt.password` | Password for encrypting (only aes)                | option   |                                                  |

#### Options

| Parameter             | Description            | Type     | Additional info               |
|-----------------------|------------------------|----------|-------------------------------|
| `options.auth_source` | Name database for auth | option   | if set up driver mongo        |
| `options.ssl`         | SSL/TLS                | option   | if set up driver mongo, mssql |
| `options.mode`        | Mode create dump       | option   | if set up driver redis        |
| `options.role`        | Role for create dump   | option   | if set up driver firebird     |
| `options.path`        | Path database SQLite   | option   | if set up driver sqlite       |

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

Or only for a specific purpose

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
./dumper --crypt backup --password 123456 --input ./dump.sql.gz.enc
```
or
```
openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in dump.sql.gz.enc -out dump.sql.gz -k 123456
```


####  🔐 5. encrypt and decrypt file config

How it works: 
1. Password encryption (outputs recovery token)
    ```
   ./dumper --crypt config --mode encrypt --password <password> --input config.yaml
   ```
2. Launching the application (reading without password)
    ```
    ./dumper --crypt config --mode decrypt --password <password> --input config.yaml
    ```
3. Decryption on the same device
    ```
    ./dumper --crypt config --mode decrypt --password <password> --input config.yaml
    ```
4. Recovery on another device
    ```
    ./dumper --crypt config --mode recover --recovery <recovery_token> --input config.yaml
    ```


### ▶ Launch examples

#### Backup with a choice of database from config file

```
./dumper --config ./cfg.yaml
````

#### Flag list

| Flag           | Example           | Description                                      | 
|----------------|-------------------|--------------------------------------------------|
| `--config`     | ./cfg.yaml        | path to config file                              |
| `--db`         | demo,app          | backup databases from list                       |
| `--all`        |                   | backup all databases from config file            |
| `--file-log`   | file.log          | file name log file (if settings.logging == true) |
| `--crypt`      | config            | Crypt file: `dump`, `config`                     |                                 |
| `--input`      | ./dump.sql.gz.enc | path to encrypt file                             |
| `--mode`       | encrypt           | Mode crypt: `encrypt`, `decrypt`, `recover`      |
| `--password`   | 123456            | password for decrypt (only for aes)              |
| `--recovery`   | 4j3k4lc7na09s     | Recovery token for recovery                      |

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