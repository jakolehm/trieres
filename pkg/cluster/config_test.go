package cluster

import (
	"testing"

	_ "github.com/jakolehm/trieres/pkg/hosts"
)

func TestNonExistingHostsFails(t *testing.T) {
	data := `
hosts:
`
	c := loadYaml(t, data)
	if err := c.Validate(); err == nil {
		t.Error("config with no hosts should fail validation")
	}
}

func TestHostAddressValidationWithIP(t *testing.T) {
	data := `
hosts:
- address: "512.1.2.3"
`
	c := loadYaml(t, data)

	err := c.Validate()
	if err == nil {
		t.Error("Host with invalid address should fail validation")
	}
}

func TestHostAddressValidationWithValidHostname(t *testing.T) {
	data := `
hosts:
- address: "foo.bar.com"
`
	c := loadYaml(t, data)

	err := c.Validate()
	if err != nil {
		t.Error("Host with invalid address should fail validation")
	}
}

func TestHostSshPortValidation(t *testing.T) {
	data := `
hosts:
- address: "1.2.3.4"
  sshPort: 0
`
	c := loadYaml(t, data)

	err := c.Validate()
	if err == nil {
		t.Error("Host with invalid ssh port should fail validation")
	}
}

func TestHostSshKeyValidation(t *testing.T) {
	data := `
hosts:
- address: "1.2.3.4"
  sshPort: 22
  sshKeyPath: /path/to/nonexisting/key
`
	c := loadYaml(t, data)

	err := c.Validate()
	if err == nil {
		t.Error("Host with invalid ssh key should fail validation")
	}
}

func TestHostRoleValidation(t *testing.T) {
	data := `
hosts:
- address: "1.2.3.4"
  sshPort: 22
  role: foobar
`
	c := loadYaml(t, data)
	err := c.Validate()
	if err == nil {
		t.Error("Host with invalid role should fail validation")
	}
}

// Just a small helper to load the config struct from yaml to get defaults etc. in place
func loadYaml(t *testing.T, data string) *Config {
	c := &Config{}
	if err := c.FromYaml([]byte(data)); err != nil {
		t.Error(err)
	}

	return c
}
