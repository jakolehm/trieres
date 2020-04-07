package hosts

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type RemoteHost interface {
	Connect() error
	Disconnect() error
}

// Config for host
type Config struct {
	Address          string   `yaml:"address" validate:"required,hostname|ip"`
	User             string   `yaml:"user"`
	SSHPort          int      `yaml:"sshPort" validate:"gt=0,lte=65535"`
	SSHKeyPath       string   `yaml:"sshKeyPath" validate:"file"`
	Role             string   `yaml:"role" validate:"oneof=master worker"`
	ExtraArgs        []string `yaml:"extraArgs"`
	PrivateInterface string   `validate:"omitempty,gt=2"`
}

// HostMetadata resolved metadata for host
type HostMetadata struct {
	Hostname        string
	InternalAddress string
}

// Host describes connectable host
type Host struct {
	Config
	sshClient *ssh.Client
	Metadata  *HostMetadata
}

// Hosts array of hosts
type Hosts []*Host

// FullAddress returns address and non-standard ssh port
func (h *Host) FullAddress() string {
	address := h.Address
	if h.Metadata != nil {
		address = h.Metadata.Hostname
	}
	if h.SSHPort != 22 {
		address = fmt.Sprintf("%s:%s", address, strconv.Itoa(h.SSHPort))
	}

	return address
}

// Connect to the host
func (h *Host) Connect() error {
	key, err := ioutil.ReadFile(h.SSHKeyPath)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	config := ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", h.Address, h.SSHPort)
	client, err := ssh.Dial("tcp", address, &config)
	if err != nil {
		return err
	}
	h.sshClient = client

	return nil
}

// Exec a command on the host
func (h *Host) Exec(cmd string) error {
	session, err := h.sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}

	logrus.Debugf("executing command: %s", cmd)
	if err := session.Start(cmd); err != nil {
		return err
	}

	multiReader := io.MultiReader(stdout, stderr)
	outputScanner := bufio.NewScanner(multiReader)

	for outputScanner.Scan() {
		logrus.Debugf("%s:  %s", h.FullAddress(), outputScanner.Text())
	}
	if err := outputScanner.Err(); err != nil {
		logrus.Errorf("%s:  %s", h.FullAddress(), err.Error())
	}

	return nil
}

// ExecWithOutput execs a command on the host and return output
func (h *Host) ExecWithOutput(cmd string) (string, error) {
	session, err := h.sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(output)), nil
}

// CopyFile copies a local file to host
func (h *Host) CopyFile(file os.File, remotePath string, permissions string) error {
	session, err := h.sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	scpClient := scp.Client{
		Session:      session,
		Timeout:      time.Second * 60,
		RemoteBinary: "scp",
	}

	return scpClient.CopyFromFile(file, remotePath, permissions)
}

// Disconnect from the host
func (h *Host) Disconnect() error {
	if h.sshClient == nil {
		return nil
	}

	return h.sshClient.Close()
}

// UnmarshalYAML unmarshals yaml
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config
	homeDir, _ := homedir.Dir()
	raw := rawConfig{
		Address:          "127.0.0.1",
		User:             "root",
		SSHKeyPath:       path.Join(homeDir, ".ssh", "id_rsa"),
		SSHPort:          22,
		Role:             "worker",
		ExtraArgs:        []string{},
		PrivateInterface: "eth0",
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}
	if strings.HasPrefix(raw.SSHKeyPath, "~") {
		raw.SSHKeyPath = path.Join(homeDir, raw.SSHKeyPath[2:])
	}

	*c = Config(raw)
	return nil
}
