package provider

import (
	"fmt"
	"os/exec"
)

// A Fortune provides a fortune
type Fortune interface {
	// Get returns a new Fortune or an error if one cannot be found
	Get() ([]byte, error)
}

// NewFortune returns a new Fortune provider that uses either the default fortune
// if the OS provided fortune command is not available
func NewFortune(defaultFortune []byte) Fortune {
	_, err := exec.LookPath("fortune")
	if err != nil {
		return &StaticFortuneProvider{staticFortune: defaultFortune}
	}

	return &OSFortuneProvider{}
}

// A OSFortuneProvider uses the fortune command from the OS to return a fortune
type OSFortuneProvider struct{}

// Get uses the OS provided fortune function which it expects to be on the PATH
func (p *OSFortuneProvider) Get() ([]byte, error) {
	cmd := exec.Command("fortune")

	fortune, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return fortune, nil
}

// A StaticFortuneProvider returns the static fortune it was created with
type StaticFortuneProvider struct {
	staticFortune []byte
}

//Get returns the static fortune
func (p *StaticFortuneProvider) Get() ([]byte, error) {
	if len(p.staticFortune) == 0 {
		return nil, fmt.Errorf("static fortune provider contains no fortune :-(")
	}

	return p.staticFortune, nil
}
