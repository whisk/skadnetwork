package skadnetwork_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	skadnetwork "github.com/whisk/skadnetwork/pkg"
)

func TestVersionFromFuture(t *testing.T) {
	v := skadnetwork.NewPostbackValidator()
	jsonBytes, err := readAllFromFile("../testdata/err/version-from-future.json")
	if err != nil {
		t.Fatal("failed to read test postback", err)
	}
	ok, err := v.Validate(jsonBytes)
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestBadSignature(t *testing.T) {
	v := skadnetwork.NewPostbackValidator()
	jsonBytes, err := readAllFromFile("../testdata/invalid/bad-signature.json")
	if err != nil {
		t.Fatal("failed to read test postback", err)
	}
	ok, err := v.Validate(jsonBytes)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.NotEmpty(t, v.Errors())
}

func TestNoRequiredParams(t *testing.T) {
	v := skadnetwork.NewPostbackValidator()
	jsonBytes, err := readAllFromFile("../testdata/invalid/no-required-params.json")
	if err != nil {
		t.Fatal("failed to read test postback", err)
	}
	ok, err := v.Validate(jsonBytes)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.NotEmpty(t, v.Errors())
}

func TestReferencePostbacks(t *testing.T) {
	v := skadnetwork.NewPostbackValidator()
	files, err := filepath.Glob("../testdata/*.json")
	if err != nil || len(files) == 0 {
		t.Fatal("failed to read reference postbacks", err)
	}
	for _, file := range files {
		jsonBytes, _ := readAllFromFile(file)
		ok, err := v.Validate(jsonBytes)
		assert.NoError(t, err, fmt.Sprintf("postback %s should be validated without errors", file))
		assert.True(t, ok, fmt.Sprintf("postback %s should be valid", file))
	}
}

func TestPostbackInit(t *testing.T) {
	v := skadnetwork.NewPostbackValidator()
	assert.IsType(t, skadnetwork.PostbackValidator{}, v)
	err := v.Init([]byte(`{"version":"3.0"}`))
	assert.NoError(t, err)
}

func TestPostbackSupportedVersion(t *testing.T) {
	v := skadnetwork.NewPostbackValidator()
	_ = v.Init([]byte(`{"version":"3.0"}`))
	ok, err := v.VersionSupported()
	assert.NoError(t, err)
	assert.True(t, ok)
}

func readAllFromFile(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
