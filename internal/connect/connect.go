package connect

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Connect struct {
	Server           string
	Username         string
	Port             string
	SSHLocalKeyPath  string
	SSHServerKeyPath string
	Passphrase       string
	IsPassphrase     bool
	Password         string
	client           *ssh.Client
}

func New(server, username, port, sshLocalKeyPath, sshServerKeyPath, passphrase, password string, isPassphrase bool) *Connect {
	return &Connect{
		Server:           server,
		Username:         username,
		Port:             port,
		SSHLocalKeyPath:  sshLocalKeyPath,
		SSHServerKeyPath: sshServerKeyPath,
		Passphrase:       passphrase,
		IsPassphrase:     isPassphrase,
		Password:         password,
	}
}

func (c *Connect) buildSSHConfig() (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	if c.SSHLocalKeyPath != "" {
		key, err := os.ReadFile(c.SSHLocalKeyPath)
		if err != nil {
			return nil, fmt.Errorf("error couldn't read SSH key: %v", err)
		}

		if c.IsPassphrase && c.Passphrase == "" {
			fmt.Print("Enter the passphrase : \n")
			passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return nil, fmt.Errorf("input error: %v", err)
			}
			c.Passphrase = string(passphrase)
		}

		var signer ssh.Signer
		if c.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(c.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}

		if err != nil {
			return nil, fmt.Errorf("error couldn't parse SSH key: %v", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if c.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("error the authentication method is not specified")
	}

	return &ssh.ClientConfig{
		User:            c.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}, nil
}

func (c *Connect) Connect() error {
	config, err := c.buildSSHConfig()
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", c.Server+":"+c.Port, config)
	if err != nil {
		return fmt.Errorf("error couldn't connect via SSH: %v", err)
	}

	c.client = client
	return nil
}

func (c *Connect) NewSession() (*ssh.Session, error) {
	if c.client == nil {
		return nil, fmt.Errorf("SSH client is not connected")
	}
	return c.client.NewSession()
}

func (c *Connect) RunCommand(cmd string) (string, error) {
	session, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

func (c *Connect) TestConnection() error {
	_, err := c.RunCommand("true")
	return err
}

func (c *Connect) Close() error {
	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		return err
	}
	return nil
}
