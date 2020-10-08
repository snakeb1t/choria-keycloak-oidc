package dictionary_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.palantir.build/ptittle/choria-keycloak-oidc/dictionary"
)

func TestStore_Add(t *testing.T) {
	s := dictionary.New()
	s.Add("hello")
	assert.True(t, s.Lookup("hello"))
}

func TestStore_Remove(t *testing.T) {
	s := dictionary.New()
	s.Add("hello")
	assert.True(t, s.Lookup("hello"))

	s.Remove("HeLlo")
	assert.False(t, s.Lookup("hello"))
}

func TestStore_Lookup(t *testing.T) {
	s := dictionary.New()
	assert.False(t, s.Lookup("hello"))

	s.Add("HelLo")
	assert.True(t, s.Lookup("hello"))
}
