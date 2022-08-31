package gutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type KeyFormat int

const (
	KeyFormatPKCS1 = 1
	KeyFormatPKCS8 = 2
	KeyFormatPKIX  = 3
)

// ParsePrivateKeyFromPem parses private key from the pem value.
func ParseRSAPrivateKey(pemBytes []byte, keyFormat KeyFormat) (*rsa.PrivateKey, error) {
	key, err := ParsePrivateKey(pemBytes, keyFormat)

	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key type is not RSA")
	}
	return rsaKey, nil
}

func ParsePrivateKey(pemBytes []byte, keyFormat KeyFormat) (key interface{}, err error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	switch keyFormat {
	case KeyFormatPKCS1:
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case KeyFormatPKCS8:
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	return
}

func ParseRSAPublicKey(pemBytes []byte, format KeyFormat) (*rsa.PublicKey, error) {
	pub, err := ParsePublicKey(pemBytes, format)
	if err != nil {
		return nil, err
	}
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key type is not RSA")
	}

	return publicKey, nil
}

func ParsePublicKey(pemBytes []byte, keyFormat KeyFormat) (key interface{}, err error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	switch keyFormat {
	case KeyFormatPKCS1:
		key, err = x509.ParsePKCS1PublicKey(block.Bytes)
	case KeyFormatPKCS8:
		key, err = x509.ParsePKIXPublicKey(block.Bytes)
	}

	return
}
