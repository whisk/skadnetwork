package skadnetwork

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Apple's NIST P-256 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/verifying_an_install-validation_postback
const APPLE_PUBLIC_KEY_21 = `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEWdp8GPcGqmhgzEFj9Z2nSpQVddayaPe4FMzqM9wib1+aHaaIzoHoLN9zW4K8y4SPykE3YVK3sVqW6Af0lfx3gg==`

// Apple's P-192 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/skadnetwork_release_notes/skadnetwork_2_release_notes
const APPLE_PUBLIC_KEY_20 = `MEkwEwYHKoZIzj0CAQYIKoZIzj0DAQEDMgAEMyHD625uvsmGq4C43cQ9BnfN2xslVT5V1nOmAMP6qaRRUll3PB1JYmgSm+62sosG`

// Apple's P-192 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/skadnetwork_release_notes/skadnetwork_1_release_notes
const APPLE_PUBLIC_KEY_10 = `MEkwEwYHKoZIzj0CAQYIKoZIzj0DAQEDMgAEMyHD625uvsmGq4C43cQ9BnfN2xslVT5V1nOmAMP6qaRRUll3PB1JYmgSm+62sosG`

var supportedVersions = map[string]bool{"2.2": true, "3.0": true, "4.0": true}

type Postback struct {
	bytes  []byte
	fields map[string]any
}

func NewPostback(bytes []byte) (Postback, error) {
	var p Postback
	var fields map[string]any
	err := json.Unmarshal(bytes, &fields)
	if err != nil {
		return p, err
	}
	p.bytes = bytes
	p.fields = fields
	return p, nil
}

func (p Postback) VerifySignature() (bool, error) {
	return true, nil
}

func (p Postback) ValidateSchema() (bool, []ValidationError, error) {
	version, ok := p.fields["version"]
	if !ok {
		return false, nil, errors.New("no version information found")
	}
	versionStr, _ := version.(string)
	if !supportedVersions[versionStr] {
		return false, nil, fmt.Errorf("validation of version %s is not supported", version)
	}
	schemaHelper, err := NewSchemaHelper(versionStr)
	if err != nil {
		return false, nil, err
	}
	return schemaHelper.Validate(p)
}
