package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	rsaPrivateKeyPEMBlockType = "RSA PRIVATE KEY"
	privateKeyPEMBlockType    = "PRIVATE KEY"
)

// protect an rsa key with a password
func encryptedKey(key *rsa.PrivateKey, password []byte) ([]byte, error) {
	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPEMBlock, err := x509.EncryptPEMBlock(rand.Reader, rsaPrivateKeyPEMBlockType, privBytes, password, x509.PEMCipher3DES)
	if err != nil {
		return nil, err
	}

	out := pem.EncodeToMemory(privPEMBlock)
	return out, nil
}

// load an encrypted private key from disk
func loadKeyFromFile(path string, password []byte) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, errors.New("PEM decode failed")
	}

	if string(password) != "" {
		b, err := x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, err
		}
		return x509.ParsePKCS1PrivateKey(b)
	}
	return x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}
