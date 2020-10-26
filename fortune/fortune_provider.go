package fortune

import (
	"os"
	"os/exec"
)

const (
	defaultBrewFortuneLocation   = "/tmp/usr/local/bin/fortune"
	defaultUbuntuFortuneLocation = "/tmp/usr/games/fortune"
	defaultFedoraFortuneLocation = "/tmp/usr/bin/fortune"
)

// A Provider provides a fortune
type Provider interface {
	// Get returns a new Fortune or an error if one cannot be found
	Get() ([]byte, error)
}

// NewProvider returns a new Fortune provider that uses either the default fortune
// if the OS provided fortune command is not available
func NewProvider(defaultFortune []byte) Provider {
	return &OSFortuneProvider{defaultFortune: defaultFortune}
}

// A OSFortuneProvider uses the fortune command from the OS to return a fortune
type OSFortuneProvider struct {
	defaultFortune []byte
}

// Get uses the OS provided fortune function which it expects to be on the PATH
func (p *OSFortuneProvider) Get() ([]byte, error) {
	var cmd *exec.Cmd

	if _, existErr := os.Stat(defaultUbuntuFortuneLocation); existErr == nil {
		cmd = exec.Command(defaultUbuntuFortuneLocation)
	} else if _, existErr := os.Stat(defaultFedoraFortuneLocation); existErr == nil {
		cmd = exec.Command(defaultFedoraFortuneLocation)
	} else if _, existErr := os.Stat(defaultBrewFortuneLocation); existErr == nil {
		cmd = exec.Command(defaultBrewFortuneLocation)
	} else if _, pathErr := exec.LookPath("fortune"); pathErr == nil {
		cmd = exec.Command("fortune")
	}

	if cmd == nil {
		return p.defaultFortune, nil
	}

	fortune, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return fortune, nil
}
