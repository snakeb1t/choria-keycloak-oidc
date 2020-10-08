package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.palantir.build/ptittle/choria-keycloak-oidc/config"
	"github.palantir.build/ptittle/choria-keycloak-oidc/server"
)

const (
	installConfigFlagName = "install-config"
	runtimeConfigFlagName = "runtime-config"
	ecvKeyFlagName        = "ecv-key"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs echo server",
	RunE: func(_ *cobra.Command, _ []string) error {
		return server.New().
			WithSelfSignedCertificate().
			WithInstallConfigType(config.InstallConfig{}).
			WithInstallConfigFromFile(viper.GetString(installConfigFlagName)).
			WithRuntimeConfigType(config.RuntimeConfig{}).
			WithRuntimeConfigFromFile(viper.GetString(runtimeConfigFlagName)).
			WithECVKeyFromFile(viper.GetString(ecvKeyFlagName)).
			Start()
	},
}

func init() {
	serverCmd.Flags().StringP(installConfigFlagName, "s", "", "install configuration file for echo service")
	if err := viper.BindPFlag(installConfigFlagName, serverCmd.Flags().Lookup(installConfigFlagName)); err != nil {
		panic(err)
	}
	viper.SetDefault(installConfigFlagName, "var/conf/install.yml")
	if err := viper.BindEnv(installConfigFlagName, "CHORIA_KEYCLOAK_OIDC_INSTALL_CONFIG"); err != nil {
		panic(err)
	}

	serverCmd.Flags().StringP(runtimeConfigFlagName, "r", "", "runtime configuration file for echo service")
	if err := viper.BindPFlag(runtimeConfigFlagName, serverCmd.Flags().Lookup(runtimeConfigFlagName)); err != nil {
		panic(err)
	}
	viper.SetDefault(runtimeConfigFlagName, "var/conf/runtime.yml")
	if err := viper.BindEnv(runtimeConfigFlagName, "CHORIA_KEYCLOAK_OIDC_RUNTIME_CONFIG"); err != nil {
		panic(err)
	}

	serverCmd.Flags().String(ecvKeyFlagName, "", "path to the file used to decrypt encrypted values in the runtime configuration")
	if err := viper.BindPFlag(ecvKeyFlagName, serverCmd.Flags().Lookup(ecvKeyFlagName)); err != nil {
		panic(err)
	}
	viper.SetDefault(ecvKeyFlagName, "var/conf/encrypted-config-value.key")

	rootCmd.AddCommand(serverCmd)
}
