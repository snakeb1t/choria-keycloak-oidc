package echo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.palantir.build/ptittle/choria-keycloak-oidc/echo"
)

func TestRun(t *testing.T) {
	for i, tc := range []struct {
		input string
		num   int
		op    echo.StringOperation
		want  string
	}{
		{"test-input", 2, echo.NoOp(), "test-input test-input"},
		{"test-input", 3, echo.Reverse(), "tupni-tset tupni-tset tupni-tset"},
	} {
		assert.Equal(t, tc.want, echo.Run(tc.input, tc.num, tc.op), "Case %d", i)
	}
}
