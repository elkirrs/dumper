# 📦 Dumper

**Dumper** — This is a CLI utility for creating backups databases of various types (PostgreSQL, MySQL and etc.) with flexible connection and storage settings.

---

![Demo](assets/dumper-proccess.gif)

## 🚀 Opportunities

- Support **PostgreSQL**, **MySQL** and etc.
- Connect for DB:
    - with the dump performed directly on the server and download (server).
- Work with **SSH-Keys** (include passphrase).
- Custom dump name templates.
- Archiving old dumps.
- Different formats:
    - PostgreSQL: `plain`, `dump`, `tar`
  
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
  test:
    name: "test server"
    host: "10.234.23.45"
    port: "22"
    user: "user"
    password: "password"

databases:
  demo:
    name: "demo"
    user: "demo_user"
    password: "password"
    server: "test"
    port: "5432"
    driver: "psql"

  app:
    user: "app_user"
    password: "pass_user"
    server: "test"

```
---

### 📑 Configuration Description

#### The configuration consists of three sections:

#### 🔧 1. settings — global settings
Apply to all servers and databases, unless redefined locally.

| Parameter           | Description                                                                               | is        |
|---------------------|-------------------------------------------------------------------------------------------|-----------|
| `db_port`           | Default database connection port                                                          | option    |
| `driver`            | The default DB driver: `psql`                                                             | required  |
| `ssh.private_key`   | The path to the private SSH key.                                                          | option    |
| `ssh.passphrase`    | Passphrase for the key (optional).                                                        | option    |
| `ssh.is_passphrase` | whether to use passphrase from the config                                                 | option    |
| `template`          | File Name Template: `{%srv%}`, `{%db%}`, `{%datetime%}`, `{%date%}`, `{%time%}`, `{%ts%}` | option    |
| `archive`           | Archiving old dumps (need `{%srv%}_{%db%}` in template).                                  | option    |
| `location`          | Dump execution method: `server`                                                           | required  |
| `format`            | Dump format: `plain`, `dump`, `tar`.                                                      | required  |
| `dir_dump`          | Directory for saving dumps                                                                | option    |
| `dir_archived`      | Archive Directory                                                                         | option    |
| `logging`           | Create logging                                                                            | option    |
| `retry_connect`     | attempts reconnect to server (default 5)                                                  | option    |


#### Params:

- #### template:
    - `{%srv%}` —  Name server
    - `{%db%}` —  Name db
    - `{%datetime%}` —  Date and time
    - `{%date%}` — Date
    - `{%time%}` — Time
    - `{%ts%}` — Time unix
- #### location:
    - `server` — create dump in server and download
- #### format:
    - PostgreSQL: `plain`, `dump`, `tar`

#### 🖥 2. servers
Defines the connections through which databases can be backed up.


| Parameter   | Description                         | is                                     |
|-------------|-------------------------------------|----------------------------------------|
| `name`      | Human-readable server name          | option                                 |
| `host`      | The IP address or domain name       | required                               |
| `port`      | Connection port                     | required<br/> (if not set global)      |
| `user`      | Username.                           | required                               |
| `password`  | Password (if there is no key)       | required<br/> (if not set key)         |


#### 🗄 3. databases
A list of databases that need to be backed up.

| Parameter   | Description                                            | is                                |
|-------------|--------------------------------------------------------|-----------------------------------|
| `name`      | Database name (by default, the key name)               | option                            |
| `user`      | The database user                                      | required                          |
| `password`  | DB user's password                                     | required                          |
| `server`    | The link to the server from the `servers` section      | required                          |
| `port`      | Connection port (if different from `settings.db_port`) | required<br/> (if not set global) |
| `driver`    | driver: `psql`                                         | required<br/> (if not set global) |

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