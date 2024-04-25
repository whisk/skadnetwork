package skadnetwork

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var supportedVersions = map[string]bool{"2.1": true, "2.2": true, "3.0": true, "4.0": true}

type Postback struct {
	bytes  []byte
	params map[string]any
}

// NewPostback returns a new postback from given JSON bytes presentation
func NewPostback(bytes []byte) (Postback, error) {
	var p Postback
	var params map[string]any
	err := json.Unmarshal(bytes, &params)
	if err != nil {
		return p, err
	}
	p.bytes = bytes
	p.params = params
	return p, nil
}

// NewPostbackFromString returns a new postback from JSON string presentation
func NewPostbackFromString(s string) (Postback, error) {
	return NewPostback([]byte(s))
}

// VersionSupported checks if postback version is supported.
func (p Postback) VersionSupported() (bool, error) {
	version, ok := p.params["version"].(string)
	if !ok {
		return false, errors.New("no version information found")
	}
	_, ok = supportedVersions[version]
	return ok, nil
}

// versionString returns string representation of postback version. Returns an empty string if
// no version data found.
func (p Postback) versionString() string {
	return p.params["version"].(string)
}

// VerifySignature verifies postback cryptographic signature. Returns an error if the version is not
// supported or the signature has an invalid format.
func (p Postback) VerifySignature() (bool, error) {
	signableString := p.signableString()
	publicKey, err := publicKey(p.versionString())
	if err != nil {
		return false, err
	}

	attrSign, _ := p.params["attribution-signature"].(string)
	attrSignBytes, err := base64.StdEncoding.DecodeString(attrSign)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256([]byte(signableString))
	return ecdsa.VerifyASN1(publicKey, hash[:], attrSignBytes), nil
}

// ValidateSchema checks postback structure using JSON schema. Returns a slice of validation errors,
// which is empty if the postback is valid.
// Returns an error if the validation itself has failed.
func (p Postback) ValidateSchema() (bool, []ValidationError, error) {
	schemaHelper, err := newSchemaHelper(p.versionString())
	if err != nil {
		return false, nil, err
	}
	return schemaHelper.validate(p)
}

func (p Postback) signableString() string {
	var partNames []string
	switch p.params["version"] {
	case "4.0":
		partNames = []string{
			"version",
			"ad-network-id",
			"source-identifier",
			"app-id",
			"transaction-id",
			"redownload",
			"source-app-id", "source-domain", // mutually exclusive
			"fidelity-type",
			"did-win",
			"postback-sequence-index",
		}
	case "3.0":
		partNames = []string{
			"version",
			"ad-network-id",
			"campaign-id",
			"app-id",
			"transaction-id",
			"redownload",
			"source-app-id",
			"fidelity-type",
			"did-win",
		}
	case "2.2":
		// see https://developer.apple.com/documentation/storekit/skadnetwork/verifying_an_install-validation_postback/combining_parameters_for_previous_skadnetwork_postback_versions#3761660
		partNames = []string{
			"version",
			"ad-network-id",
			"campaign-id",
			"app-id",
			"transaction-id",
			"redownload",
			"source-app-id",
			"fidelity-type",
		}
	case "2.1":
		fallthrough
	case "2.0":
		// see https://developer.apple.com/documentation/storekit/skadnetwork/verifying_an_install-validation_postback/combining_parameters_for_previous_skadnetwork_postback_versions#3626226
		partNames = []string{
			"version",
			"ad-network-id",
			"campaign-id",
			"app-id",
			"transaction-id",
			"redownload",
			"source-app-id",
		}
	case "1.0":
	}

	parts := []string{}
	for _, name := range partNames {
		paramVal, ok := p.params[name]
		if !ok {
			// skip non-existing fields
			continue
		}

		var strVal string
		switch v := (paramVal).(type) {
		case string:
			strVal = v
		case int, uint:
			strVal = fmt.Sprintf("%d", v)
		case float64:
			strVal = fmt.Sprintf("%0.0f", v)
		case bool:
			if v {
				strVal = "true"
			} else {
				strVal = "false"
			}
		default:
			strVal = fmt.Sprintf("%s_has_unsupported_type_%T", name, v)
		}
		parts = append(parts, strVal)
	}

	return strings.Join(parts, "\u2063")
}
