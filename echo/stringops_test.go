package echo_test

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.palantir.build/ptittle/choria-keycloak-oidc/echo"
)

func TestMain(m *testing.M) {
	// manually set rand seed so that randomized output is deterministic
	rand.Seed(0)
	os.Exit(m.Run())
}

func TestOps(t *testing.T) {
	for i, tc := range []struct {
		op   echo.StringOperation
		in   string
		want string
	}{
		{echo.NoOp(), "hello", "hello"},
		{echo.Reverse(), "hello", "olleh"},
		{echo.Randomize(), "hello", "ellho"},
	} {
		assert.Equal(t, tc.want, tc.op(tc.in), "Case %d", i)
	}
}
