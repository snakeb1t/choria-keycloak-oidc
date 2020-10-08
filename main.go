package main

import (
	"os"

	"github.com/palantir/pkg/signals"
	"github.palantir.build/ptittle/choria-keycloak-oidc/cmd"
)

func main() {
	signals.RegisterStackTraceWriter(os.Stderr, nil)
	os.Exit(cmd.Execute())
}
