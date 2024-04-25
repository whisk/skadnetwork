package skadnetwork_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	skadnetwork "github.com/whisk/skadnetwork/pkg"
)

func TestInvalidPostback(t *testing.T) {
	p, err := skadnetwork.NewPostback([]byte(`{'invalid json here':`))
	assert.Error(t, err)
	assert.Zero(t, p)
}

func TestPostbackStub(t *testing.T) {
	p, err := skadnetwork.NewPostback([]byte(`{"version":"3.0"}`))
	assert.NoError(t, err)
	assert.IsType(t, skadnetwork.Postback{}, p)
}

func TestPostbackVersion(t *testing.T) {
	p, err := skadnetwork.NewPostback([]byte(`{"version":"3.0"}`))
	assert.NoError(t, err)
	assert.IsType(t, skadnetwork.Postback{}, p)
	ok, err := p.VersionSupported()
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestPostbackNotSupportedVersion(t *testing.T) {
	p, err := skadnetwork.NewPostback([]byte(`{"version":"5.0"}`))
	assert.NoError(t, err)
	assert.IsType(t, skadnetwork.Postback{}, p)
	ok, err := p.VersionSupported()
	assert.NoError(t, err)
	assert.False(t, ok)
}
