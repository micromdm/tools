// Integration with Jesse's mdmcert.download

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	mdmcertRequestURL = "https://mdmcert.download/api/v1/signrequest"
	// see
	// https://github.com/jessepeterson/commandment/blob/1352b51ba6697260d1111eccc3a5a0b5b9af60d0/commandment/mdmcert.py#L23-L28
	mdmcertServerKey = "f847aea2ba06b41264d587b229e2712c89b1490a1208b7ff1aafab5bb40d47bc"
)

// format of a signing request to mdmcert.download
type signRequest struct {
	CSR     string `json:"csr"` // base64 encoded PEM CSR
	Email   string `json:"email"`
	Key     string `json:"key"`     // server key from above
	Encrypt string `json:"encrypt"` // server cert
}

func newSignRequest(email string, pemCSR []byte, serverCertificate []byte) *signRequest {
	encodedCSR := base64.StdEncoding.EncodeToString(pemCSR)
	encodedServerCert := base64.StdEncoding.EncodeToString(serverCertificate)
	return &signRequest{
		CSR:     encodedCSR,
		Email:   email,
		Key:     mdmcertServerKey,
		Encrypt: encodedServerCert,
	}
}

func (sign *signRequest) HTTPRequest() (*http.Request, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(sign); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", mdmcertRequestURL, ioutil.NopCloser(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "micromdm/certhelper")
	return req, nil
}

func sendRequest(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received bad status from mdmcert.download. status=%q", resp.Status)
	}
	var jsn = struct {
		result string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&jsn); err != nil {
		return err
	}
	if jsn.result != "success" {
		return fmt.Errorf("got unexpected result body: %q\n", jsn.result)
	}
	return nil
}
