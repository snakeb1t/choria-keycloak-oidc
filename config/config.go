package config

import (
	"strings"

	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-server/config"
	"github.com/snakeb1t/choria-keycloak-oidc/echo"
)

// InstallConfig contains the install configuration of the server
type InstallConfig struct {
	config.Install `yaml:",inline"`
}

// RuntimeConfig contains the runtime configuration of the server
type RuntimeConfig struct {
	config.Runtime `yaml:",inline"`

	EchoCount int                 `yaml:"echo-count"`
	StringOp  StringOperationType `yaml:"string-op"`
}

// StringOperationType is a strongly typed string that represents an operation type
type StringOperationType string

const (
	// Noop is a noop operation
	Noop StringOperationType = "noop"
	// Reverse reverses the string
	Reverse StringOperationType = "reverse"
	// Randomize randomizes the string
	Randomize StringOperationType = "randomize"
)

// StringOp returns the StringOperation of the StringOperationType
func (s StringOperationType) StringOp() echo.StringOperation {
	return stringOps[s]
}

var stringOps = map[StringOperationType]echo.StringOperation{
	Noop:      echo.NoOp(),
	Reverse:   echo.Reverse(),
	Randomize: echo.Randomize(),
}

// UnmarshalYAML initializes the StringOperationType from a provided unmarshal function.
func (s *StringOperationType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return werror.Wrap(err, "failed to unmarshal yaml")
	}
	str = strings.ToLower(str)
	switch StringOperationType(str) {
	default:
		return werror.Error("unrecognized StringOperationType", werror.SafeParam("type", str))
	case Noop, Reverse, Randomize:
		*s = StringOperationType(str)
	}
	return nil
}
