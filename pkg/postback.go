package skadnetwork

// Apple's NIST P-256 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/verifying_an_install-validation_postback
const APPLE_PUBLIC_KEY_21 = `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEWdp8GPcGqmhgzEFj9Z2nSpQVddayaPe4FMzqM9wib1+aHaaIzoHoLN9zW4K8y4SPykE3YVK3sVqW6Af0lfx3gg==`
// Apple's P-192 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/skadnetwork_release_notes/skadnetwork_2_release_notes
const APPLE_PUBLIC_KEY_20 = `MEkwEwYHKoZIzj0CAQYIKoZIzj0DAQEDMgAEMyHD625uvsmGq4C43cQ9BnfN2xslVT5V1nOmAMP6qaRRUll3PB1JYmgSm+62sosG`
// Apple's P-192 public key
// https://developer.apple.com/documentation/storekit/skadnetwork/skadnetwork_release_notes/skadnetwork_1_release_notes
const APPLE_PUBLIC_KEY_10 = `MEkwEwYHKoZIzj0CAQYIKoZIzj0DAQEDMgAEMyHD625uvsmGq4C43cQ9BnfN2xslVT5V1nOmAMP6qaRRUll3PB1JYmgSm+62sosG`

type Postback map[string]any

func (p *Postback) Verify() (bool, error) {
	return true, nil
}
