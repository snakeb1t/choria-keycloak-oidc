package server

import (
	"context"

	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-server/witchcraft"
	"github.palantir.build/ptittle/choria-keycloak-oidc/conjure/choriakeycloakoidc/api"
)

func New() *witchcraft.Server {
	return witchcraft.NewServer().
		WithInitFunc(func(ctx context.Context, info witchcraft.InitInfo) (cleanup func(), rErr error) {
			// TODO: create sync.Map, pass it in as pointer to authorizer and redirecter
			if err := api.RegisterRoutesAuthorize(info.Router, &authorizer{}); err != nil {
				return nil, err
			}
			somethingInConfig := "foo"

			err := info.Router.Get("/login", &redirecter{
				redirectURI: somethingInConfig,
			})
			if err != nil {
							  return nil, err
							  }

			return nil, nil
			}).
		WithOrigin(svc1log.CallerPkg(0, 1))
}
