package server_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/conjure-go-runtime/conjure-go-client/httpclient"
	"github.com/palantir/pkg/httpserver"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wserverconfig "github.com/palantir/witchcraft-go-server/config"
	"github.com/palantir/witchcraft-go-server/status"
	"github.com/palantir/witchcraft-go-server/witchcraft"
	"github.com/stretchr/testify/require"
	"github.palantir.build/ptittle/choria-keycloak-oidc/config"
	"github.palantir.build/ptittle/choria-keycloak-oidc/conjure/choriakeycloakoidc/api"
	"github.palantir.build/ptittle/choria-keycloak-oidc/server"
)

func TestServer(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	restorer, err := dirs.SetwdWithRestorer(tmpDir)
	require.NoError(t, err)
	defer restorer()

	appPort, err := httpserver.AvailablePort()
	require.NoError(t, err)
	mgmtPort, err := httpserver.AvailablePort()
	require.NoError(t, err)

	server := server.New().
		WithInstallConfigType(config.InstallConfig{}).
		WithInstallConfig(config.InstallConfig{
			Install: wserverconfig.Install{
				Server: wserverconfig.Server{
					Address:        "localhost",
					ContextPath:    "/test",
					Port:           appPort,
					ManagementPort: mgmtPort,
				},
				UseConsoleLog: true,
			},
		}).
		WithRuntimeConfigType(config.RuntimeConfig{}).
		WithRuntimeConfig(config.RuntimeConfig{
			Runtime: wserverconfig.Runtime{
				LoggerConfig: &wserverconfig.LoggerConfig{
					Level: wlog.DebugLevel,
				},
			},
			EchoCount: 3,
			StringOp:  config.Reverse,
		}).
		WithECVKeyProvider(witchcraft.ECVKeyNoOp()).
		WithSelfSignedCertificate()

	go func() {
		if err := server.Start(); err != nil && err.Error() != "http: Server closed" {
			fmt.Println("server failed:", err)
		}
	}()
	defer func() {
		if err := server.Close(); err != nil {
			fmt.Println("failed to close server:", err)
		}
	}()
	require.True(t, <-httpserver.Ready(func() (*http.Response, error) {
		resp, err := testServerClient().Get(fmt.Sprintf("https://localhost:%d/%s/%s", mgmtPort, "test", status.LivenessEndpoint))
		return resp, err
	}))

	httpClient, err := httpclient.NewClient(
		httpclient.WithBaseURLs([]string{fmt.Sprintf("https://localhost:%d/test", appPort)}),
		httpclient.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	)
	require.NoError(t, err)
	dictClient := api.NewChoriaKeycloakOidcDictServiceClient(httpClient)

	getResp1, err := dictClient.Get(context.Background(), "foo")
	require.NoError(t, err)
	require.Equal(t, false, getResp1)

	err = dictClient.Put(context.Background(), "foo")
	require.NoError(t, err)

	getResp2, err := dictClient.Get(context.Background(), "foo")
	require.NoError(t, err)
	require.Equal(t, true, getResp2)
}

func testServerClient() *http.Client {
	return &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
}
