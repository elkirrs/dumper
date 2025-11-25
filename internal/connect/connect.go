package connect

import (
	"bytes"
	"context"
	connectDomain "dumper/internal/domain/connect"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Connect struct {
	ctx     context.Context
	connect *connectDomain.Connect
	client  *ssh.Client
}

func NewApp(
	ctx context.Context,
	connect *connectDomain.Connect,
) *Connect {
	return &Connect{
		ctx:     ctx,
		connect: connect,
	}
}

func (c *Connect) buildSSHConfig() (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	if c.connect.PrivateKey != "" {
		key, err := os.ReadFile(c.connect.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("error couldn't read SSH key: %v", err)
		}

		if c.connect.IsPassphrase && c.connect.Passphrase == "" {
			fmt.Printf("Enter the passphrase : ")
			passphrase, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Printf("\r")
			if err != nil {
				return nil, fmt.Errorf("input error: %v", err)
			}
			c.connect.Passphrase = strings.TrimSpace(string(passphrase))
		}

		var signer ssh.Signer
		if c.connect.Passphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(c.connect.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}

		if err != nil {
			return nil, fmt.Errorf("error couldn't parse SSH key: %v", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if c.connect.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.connect.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("error the authentication method is not specified")
	}

	return &ssh.ClientConfig{
		User:            c.connect.Username,
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

	fmt.Printf("Trying to connect to %s...\n", c.connect.Server)
	logging.L(c.ctx).Info(
		"Trying to test connection to ",
		logging.StringAttr("server", c.connect.Server),
	)

	client, err := ssh.Dial("tcp", c.connect.Server+":"+c.connect.Port, config)
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

func (c *Connect) Client() *ssh.Client {
	return c.client
}

func (c *Connect) IsConnected() bool {
	if c.client == nil {
		return false
	}
	_, _, err := c.client.SendRequest("keepalive@openssh.com", true, nil)
	return err == nil
}

func (c *Connect) Reconnect() error {
	fmt.Println("[SSH] Attempting reconnect...")

	_ = c.Close()
	time.Sleep(2 * time.Second)
	return c.Connect()
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
