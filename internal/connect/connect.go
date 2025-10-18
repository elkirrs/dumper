package connect

import (
	"bytes"
	"dumper/pkg/utils"
	"fmt"
	"os"
	"strings"
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

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	checkBashCmd := "command -v bash >/dev/null 2>&1 && echo OK"
	checkSession, err := c.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to check bash availability: %w", err)
	}
	var checkOut bytes.Buffer
	checkSession.Stdout = &checkOut
	if err := checkSession.Run(checkBashCmd); err != nil {
		_ = checkSession.Close()
		return "", fmt.Errorf("failed to run bash check: %w", err)
	}
	_ = checkSession.Close()

	var fullCmd string
	if strings.Contains(checkOut.String(), "OK") {
		fullCmd = fmt.Sprintf(`bash -c 'set -o pipefail; %s'`, cmd)
	} else {
		fullCmd = fmt.Sprintf(`sh -c '%s; exit ${PIPESTATUS[0]:-0}'`, cmd)
	}

	err = session.Run(fullCmd)
	output := stdout.String()
	errorOutput := stderr.String()

	if err != nil {
		return output + errorOutput, fmt.Errorf(
			"command failed: %v\nstderr: %s",
			err, utils.Mask(errorOutput),
		)
	}

	return output, nil
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
