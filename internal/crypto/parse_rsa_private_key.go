package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ParseRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	if block, _ := pem.Decode(data); block != nil {
		data = block.Bytes
	}

	key, err := x509.ParsePKCS8PrivateKey(data)
	if err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is %T, not *rsa.PrivateKey", key)
		}

		return rsaKey, nil
	}

	rsaKey, pkcs1Err := x509.ParsePKCS1PrivateKey(data)
	if pkcs1Err == nil {
		return rsaKey, nil
	}

	return nil, fmt.Errorf("failed parsing RSA private key as PKCS#8: %w; failed parsing as PKCS#1: %w", err, pkcs1Err)
}
