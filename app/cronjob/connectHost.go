package cronjob

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
}

func ExecCommandOnHost(command string) (string, error) {

	hostIP := os.Getenv("hostip")
	hostUser := os.Getenv("hostusername")
	hostPassword := os.Getenv("hostpassword")

	conn, err := Connect(hostIP, hostUser, hostPassword)
	if err != nil {
		return "Error Connecting", err
	}

	output, err := conn.SendCommands(command)
	if err != nil {
		return "Error Executing Command", err
	}

	return string(output), nil
}

func Connect(addr, user, password string) (*Connection, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil

}

func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return []byte{}, err
	}
	cmd := strings.Join(cmds, "; ")
	output, err := session.Output(cmd)
	if err != nil {
		return output, fmt.Errorf("failed to execute command '%s' on server: %v", cmd, err)
	}

	return output, err
}
