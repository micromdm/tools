package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"net/http"
)

const (
	certificatePEMBlockType = "CERTIFICATE"
)

func pemCert(derBytes []byte) []byte {
	pemBlock := &pem.Block{
		Type:    certificatePEMBlockType,
		Headers: nil,
		Bytes:   derBytes,
	}
	out := pem.EncodeToMemory(pemBlock)
	return out
}

func loadDERCertFromFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	crt, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, err
	}
	return crt.Raw, nil
}

func loadCertfromHTTP(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	crt, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, err
	}
	return crt.Raw, nil
}
