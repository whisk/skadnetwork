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
	version, ok := p.params["version"].(string)
	if !ok {
		return false, errors.New("invalid version")
	}
	publicKey, err := publicKey(version)
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
