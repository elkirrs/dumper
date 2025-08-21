# 📦 Dumper

**Dumper** — This is a CLI utility for creating backups databases of various types (PostgreSQL, MySQL and etc.) with
flexible connection and storage settings. 

---

![Demo](assets/dumper-proccess.gif)

## 🚀 Opportunities

- Multiple database systems can be managed.
- Support **PostgreSQL**, **MySQL**, **MongoDB** and etc.
- Connect for DB:
    - with the dump performed directly on the server and download (server).
- Work with **SSH-Keys** (include passphrase).
- Custom dump name templates.
- Archiving old dumps.
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

servers:
  test-server:
    name: "test server"
    host: "10.234.23.45"
    port: "22"
    user: "user"
    password: "password"
  mongo-server:
    name: "mongo"
    host: "172.0.18.54"
    user: "root"

databases:
  demo:
    name: "db_name_demo"
    user: "demo_user"
    password: "password"
    server: "test-server"
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
    server: "test-server"
    
  mongo:
    name: "mongo_db_name"
    user: "root"
    password: "mongo_password"
    port: 27017
    driver: "mongo"
    format: "bson"
    server: "mongo-server"
    options:
      auth_source: "admin"
      ssl: true

```

---

### 📑 Configuration Description

#### The configuration consists of three sections:

#### 🔧 1. settings — global settings

Apply to all servers and databases, unless redefined locally.

| Parameter           | Description                                              | is       |
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

#### Params:

- #### Driver:
    - PostgreSQL — `psql`
    - MySQL — `mysql`
    - MongoDB — `mongo`
- #### Format:
    - PostgreSQL: `plain`, `dump`, `tar`
    - MySQL: `sql`
    - MongoDB: `bson`
- #### Template:
    - `{%srv%}` — Name server
    - `{%db%}` — Name db
    - `{%datetime%}` — Date and time
    - `{%date%}` — Date
    - `{%time%}` — Time
    - `{%ts%}` — Time unix
- #### Location:
    - `server` — create dump in server and download

#### 🖥 2. servers

Defines the connections through which databases can be backed up.

| Parameter  | Description                   | is                                         |
|------------|-------------------------------|--------------------------------------------|
| `name`     | Human-readable server name    | option                                     |
| `host`     | The IP address or domain name | required                                   |
| `port`     | Connection port               | required<br/> (if not set `settings.port`) |
| `user`     | Username                      | required                                   |
| `password` | Password (if there is no key) | required<br/> (if not set `key`)           |

#### 🗄 3. databases

A list of databases that need to be backed up.

| Parameter             | Description                                       | is                                            |
|-----------------------|---------------------------------------------------|-----------------------------------------------|
| `name`                | Database name (by default, the key name)          | option                                        |
| `user`                | The database user                                 | required                                      |
| `password`            | DB user's password                                | required                                      |
| `server`              | The link to the server from the `servers` section | required                                      |
| `port`                | Connection port                                   | required<br/> (if not set `settings.db_port`) |
| `driver`              | [The DB driver list](#Driver)                     | required<br/> (if not set `settings.driver`)  |
| `format`              | [The dump format](#Format)                        | required<br/> (if not set `settings.format`)  |
| `options.auth_source` | Name database for auth                            | option (if set up driver mongo)               |
| `options.ssl`         | SSL/TLS                                           | option (if set up driver mongo)               |

---

### ▶ Launch examples

#### Backup with a choice of database from config file

```
./dumper --config ./cfg.yaml
````

- Flags:
    - `--config ./cfg.yaml` — path to config file
    - `--db demo,app` — backup databases from list
    - `--all` — backup all databases from config file
    - `--file-log file.log` — file name log file (if settings.logging == true)

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