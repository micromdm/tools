// Integration with Jesse's mdmcert.download

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fullsailor/pkcs7"
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
		Result string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&jsn); err != nil {
		return err
	}
	if jsn.Result != "success" {
		return fmt.Errorf("got unexpected result body: %q\n", jsn.Result)
	}
	return nil
}

// The user will receive a hex encoded p7 file as an email attachment.
// the file contents is a pkcs7 file, encrypted using the server certificate as the
// intended recipient.
// We use the server private key to decode the pkcs7 envelope and extract a
// base64 encoded plist (same format as the once created by `mapkePushRequestPlist`
// Once the pkcs7 file is decrypted, we save the file to disk for the user to upload
// to identity.apple.com for a push certificate.
func decodeSignedRequest(p7Path, certPath, privPath, privPass string) error {
	hexBytes, err := ioutil.ReadFile(p7Path)
	if err != nil {
		return err
	}
	key, err := loadKeyFromFile(privPath, []byte(privPass))
	if err != nil {
		return err
	}
	certPemBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}
	pemBlock, _ := pem.Decode(certPemBytes)
	if pemBlock == nil {
		return errors.New("PEM decode failed")
	}
	if pemBlock.Type != "CERTIFICATE" {
		return errors.New("certificate: unmatched type or headers")
	}
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}
	pkcsBytes, err := hex.DecodeString(string(hexBytes))
	if err != nil {
		return err
	}
	p7, err := pkcs7.Parse(pkcsBytes)
	if err != nil {
		return err
	}
	content, err := p7.Decrypt(cert, key)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf("mdmcert.download_%s", pushRequestFilename), content, 0666)
}
