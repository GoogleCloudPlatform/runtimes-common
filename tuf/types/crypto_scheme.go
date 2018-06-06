package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	ECDSA256 string = "ECDSA256"
	RSA256   string = "RSA256"
)

var VALID_CRYPTO_SCHEMES = map[string]bool{
	ECDSA256: true,
	RSA256:   false, // Not implemented.
}

var CryptoSchemes []string
var ImplementedSchemes []string

type CryptoScheme struct {
	Scheme string
}

func (scheme *CryptoScheme) String() string {
	return scheme.Scheme
}

func (scheme *CryptoScheme) Set(s string) error {
	value, ok := VALID_CRYPTO_SCHEMES[s]
	if ok && value {
		scheme.Scheme = s
		return nil
	}
	if !ok {
		return fmt.Errorf(`%s is not a valid CryptoScheme.
		Please Provide one of %s`, s, strings.Join(CryptoSchemes, ", "))
	}
	return fmt.Errorf(`%s is not a Not Implemented Yet!
		Please Provide one of %s`, s, strings.Join(ImplementedSchemes, ", "))
}

func (scheme *CryptoScheme) Type() string {
	return "types.CryptoScheme"
}

func NewCryptoScheme(val string, p *CryptoScheme) *CryptoScheme {
	value, ok := VALID_CRYPTO_SCHEMES[val]
	if ok && value {
		*p = CryptoScheme{
			Scheme: val,
		}
		return p
	}
	return nil
}

func (scheme *CryptoScheme) Store(filename string) error {
	schemeJson, err := json.Marshal(scheme)
	if err != nil {
		return fmt.Errorf("Error while marshalling json %s", err.Error())
	}
	return ioutil.WriteFile(filename, schemeJson, 0644)
}

func init() {
	for k, v := range VALID_CRYPTO_SCHEMES {
		CryptoSchemes = append(CryptoSchemes, k)
		if v {
			ImplementedSchemes = append(ImplementedSchemes, k)
		}
	}
}
