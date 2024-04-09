package skadnetwork

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
)

// Apple's NIST P-256 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/verifying_an_install-validation_postback
var applePublicKey21 = decodePublicKey(`MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEWdp8GPcGqmhgzEFj9Z2nSpQVddayaPe4FMzqM9wib1+aHaaIzoHoLN9zW4K8y4SPykE3YVK3sVqW6Af0lfx3gg==`)

// ECDSA P-192 keys are not supported in x509 due to their poor security
// See https://github.com/golang/go/issues/41035 for more details
// While it's possible to support them, that would require too much effort for almost unused old versions,
// so I leave it as it is

// Apple's P-192 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/skadnetwork_release_notes/skadnetwork_2_release_notes
var _ = decodeInsecurePublicKey(`MEkwEwYHKoZIzj0CAQYIKoZIzj0DAQEDMgAEMyHD625uvsmGq4C43cQ9BnfN2xslVT5V1nOmAMP6qaRRUll3PB1JYmgSm+62sosG`)

// Apple's P-192 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/skadnetwork_release_notes/skadnetwork_1_release_notes
var _ = decodeInsecurePublicKey(`MEkwEwYHKoZIzj0CAQYIKoZIzj0DAQEDMgAEMyHD625uvsmGq4C43cQ9BnfN2xslVT5V1nOmAMP6qaRRUll3PB1JYmgSm+62sosG`)

func publicKey(version string) (*ecdsa.PublicKey, error) {
	switch version {
	case "4.0":
		fallthrough
	case "3.0":
		fallthrough
	case "2.2":
		fallthrough
	case "2.1":
		return applePublicKey21, nil
	case "2.0":
		fallthrough
	case "1.0":
		return nil, fmt.Errorf("apple's public key for version %s is not supported", version)
	case "":
		return nil, errors.New("undefined version")
	default:
		return nil, fmt.Errorf("version %s is not supported", version)
	}
}
func decodePublicKey(key string) *ecdsa.PublicKey {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(fmt.Sprintf("failed to decode Apple public key %s from base64: %s. This is probably a bug", key, err.Error()))
	}
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		panic(fmt.Sprintf("failed to parse Apple public key %s: %s. This is probably a bug", key, err.Error()))
	}
	ecdsaKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		panic(fmt.Sprintf("Apple key %s is not an ECDSA public key. This is probably a bug", key))
	}
	return ecdsaKey
}

func decodeInsecurePublicKey(_ string) *ecdsa.PublicKey {
	return nil
}
