package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/big"
	"net/http"
	url2 "net/url"
	"os"
	"time"
)

// this is a base64'ed rsa pubkey
const KeycloakPubKeyModulusB64 = "hYB3d4JAuUC1pGlP45Ym2wMr62itCdnRkoNnwIUWRxD4UzNrKJ5vuZZMeFZqz-PYas0u1xV_7aCUbn6soPhzlC4Okp3kePoi79cTcahFqzT5veNiKUWYMKp4yweMtkaGdaB9lOffpYFv6XeSTqcooWMTVqeTl2xamseelzdxjHwwR0qgCEmBVYjaOiETm6v0s0aoWRYASX2Dhj_aVa2bxkfKNJXASsc74BpqzqGkXECk5J3O6ddLellM0zyeVfxdJCtPxOr5Hlf1KJM4GPh99Y-8EVVRX4XOCDLXI-OTY1g5Qy7TdKVsW3RsuJP05JuzFw4Cfn0skJUVJUukodxMHw"
const KeycloakPubKeyExponentB64 = "AQAB"

type authorizer struct {

}

type KeycloakClaims struct {
	Groups []string `json:"groups"`
	Username string `json:"preferred_username"`
	jwt.StandardClaims
}

func (a *authorizer) Get(ctx context.Context, stateArg string, sessionstateArg string, idtokenArg string) (string, error) {

	// need to verify the idtoken that we get from keycloak
	var claimsStruct KeycloakClaims
	token, err := jwt.ParseWithClaims(idtokenArg, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// TODO: fetch public key from jwks endpoint or config instead

		// get the keycloak public key or certificate and return it as a *rsa.PublicKey object
		modBytes, err := base64.RawURLEncoding.DecodeString(KeycloakPubKeyModulusB64)
		if err != nil {
			return nil, err
		}
		expBytes, err := base64.RawURLEncoding.DecodeString(KeycloakPubKeyExponentB64)
		if err != nil {
			return nil, err
		}
		expBig := new(big.Int).SetBytes(expBytes)
		pubkey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(modBytes),
			E: int(expBig.Int64()),
		}

		return pubkey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*KeycloakClaims); ok && token.Valid {
			// now create the choria jwt and send it back in the body

		if len(claims.Groups) == 0 {
			return errors.New("group claims missing from jwt").Error(), nil
		}
		hostname, err := os.Hostname()
		if err != nil {
			return "", err
		}
		// TODO: check that token nonce is valid. that means doing a lookup in state and if found, delete nonce
		//       from state and proceed successfully
		choriaToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"groups": claims.Groups,
			"callerid": fmt.Sprintf("up=%s", claims.Username),
			"sub": fmt.Sprintf("up=%s", claims.Username),
			"nbf": time.Now().Unix(),
			"exp": time.Now().Add(time.Second * 3600).Unix(),
			"aud": "choria",
			"iat": time.Now().Unix(),
			"jti": uuid.New(),
			"iss": hostname,
		})

		// make an *rsa.PrivateKey object
		priv, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			return "", err
		}
		choriaTokenStr, err := choriaToken.SignedString(priv)
		if err != nil {
			return "", err
		}
		return choriaTokenStr, nil
	}
	return fmt.Sprint("token is not valid!"), nil
	// generate choria-friendly JWT and send it back to the user
	//return fmt.Sprintf("state: %s / session-state: %s", stateArg, sessionstateArg), nil
}

type redirecter struct {
	redirectURI string
}

func (r *redirecter) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO: get host and path from config
	url := url2.URL{
		Host: "redacted",
		Scheme: "https",
		Path: "/auth/realms/rockland/protocol/openid-connect/auth",
	}
	nonce := uuid.New()
	state := uuid.New()
	query := url.Query()
	query.Set("response_type", "id_token")
	query.Set("client_id", "choria")
	// TODO: get redirect_uri from config
	query.Set("redirect_uri", "https://localhost:8080/choria-keycloak-oidc/authorize")
	//url.Query().Add("redirect_uri", r.redirectURI)
	query.Set("scope", "openid roles")
	query.Set("state", state.String())
	query.Set("nonce", nonce.String())
	url.RawQuery = query.Encode()

	// TODO: store nonce in state, to be checked when we get the jwt back from keycloak
	http.Redirect(resp, req, url.String(), http.StatusFound)
	fmt.Fprintf(os.Stderr, "target: %s", url.String())
}
