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