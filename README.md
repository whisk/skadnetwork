# Description

This package provides a validator for [Apple SKAdNetwork](https://developer.apple.com/documentation/storekit/skadnetwork) postbacks, covering both JSON schema validation and [signature verification](https://developer.apple.com/documentation/storekit/skadnetwork/verifying_an_install-validation_postback) for versions from 2.1 up to 4.0.

# Installation

```
go get github.com/whisk/skadnetwork
```

# Usage

```go
validator := skadnetwork.NewPostbackValidator()
ok, err := validator.Validate(jsonBytes)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error validating postback: %s\n", err)
    os.Exit(1)
}
if !ok {
    fmt.Println("Postback is not valid. Errors found:")
    for _, e := range validator.Errors() {
        fmt.Println(e.Error())
    }
    os.Exit(1)
}
fmt.Println("Postback is valid")
```

# Known issues

* Versions 2.0 and older are not supported due to difficulties with P-192 public key verification
* JSON schema validation does not cover all edge cases, for example, the mutual exclusion of `source-domain` and `source-app-id` params
* Any extra params in the postback are not checked or validated
