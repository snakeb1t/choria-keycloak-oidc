package cmd

import (
	"github.com/palantir/pkg/cobracli"
	"github.com/spf13/cobra"
)

var (
	// Version of the program.
	// TODO: Enable this by fixing the import path in godel/config/dist.yml
	Version = "unspecified"

	rootCmd = &cobra.Command{
		Use:   "choria-keycloak-oidc",
		Short: "An example service that provides some text operations",
	}
)

// Execute runs the application and returns the exit code.
func Execute() int {
	return cobracli.ExecuteWithDefaultParams(rootCmd, cobracli.VersionFlagParam(Version))
}
