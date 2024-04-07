package skadnetwork

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	params map[string]any
}

func NewPostback(bytes []byte) (Postback, error) {
	var p Postback
	var fields map[string]any
	err := json.Unmarshal(bytes, &fields)
	if err != nil {
		return p, err
	}
	p.bytes = bytes
	p.params = fields
	return p, nil
}

func (p Postback) VerifySignature() (bool, error) {
	signableString := p.signableString()
	publicKey, err := p.publicKey()
	if err != nil {
		return false, err
	}

	attrSign, _ := p.params["attribution-signature"].(string)
	attrSignBytes, err := base64.StdEncoding.DecodeString(attrSign)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256([]byte(signableString))
	if ecdsa.VerifyASN1(publicKey, hash[:], attrSignBytes) {
		return true, nil
	}
	return false, errors.New("invalid signature")
}

func (p Postback) ValidateSchema() (bool, []ValidationError, error) {
	version, ok := p.params["version"]
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

func (p Postback) publicKey() (*ecdsa.PublicKey, error) {
	version, ok := p.params["version"]
	if !ok {
		return nil, errors.New("version not defined")
	}

	var publicKeyBytes []byte
	var err error
	switch version {
	case "4.0":
		fallthrough
	case "3.0":
		fallthrough
	case "2.1":
		publicKeyBytes, err = base64.StdEncoding.DecodeString(APPLE_PUBLIC_KEY_21)
	case "2.0":
		publicKeyBytes, err = base64.StdEncoding.DecodeString(APPLE_PUBLIC_KEY_20)
	case "1.0":
		publicKeyBytes, err = base64.StdEncoding.DecodeString(APPLE_PUBLIC_KEY_10)
	default:
		err = fmt.Errorf("version %s is not supported", p.params["version"])
	}
	if err != nil {
		return nil, err
	}
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Apple public key: %w. This is probably a bug", err)
	}
	ecdsaKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not an ECDSA public key. This is probably a bug")
	}

	return ecdsaKey, nil
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
	case "2.1":
	case "2.0":
	case "1.0":
	}

	parts := []string{}
	for _, name := range partNames {
		fieldVal, ok := p.params[name]
		if !ok {
			// skip non-existing fields
			continue
		}

		var strVal string
		switch v := (fieldVal).(type) {
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
