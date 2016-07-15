package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/groob/plist"
)

var (
	version              = "unreleased"
	gitHash              = "unknown"
	wwdrIntermediaryURL  = "https://developer.apple.com/certificationauthority/AppleWWDRCA.cer"
	appleRootCAURL       = "http://www.apple.com/appleca/AppleIncRootCertificate.cer"
	providerCSRFilename  = "ProviderUnsignedPushCertificateRequest.csr"
	providerPKeyFilename = "ProviderPrivateKey.key"
	vendorPKeyFilename   = "VendorPrivateKey.key"
	vendorCSRFilename    = "VendorCertificateRequest.csr"
	pushRequestFilename  = "PushCertificateRequest"
)

func main() {
	// vendor cmd flags
	vendorCMD := flag.NewFlagSet("vendor", flag.ExitOnError)
	vendorCSRFlag := vendorCMD.Bool("csr", false, "create a CSR for MDM vendor certificate")
	vendorCSREmail := vendorCMD.String("email", "", "email address to use in CSR request Subject")
	vendorCSRCountry := vendorCMD.String("country", "US", "two letter country flag for CSR Subject(example: US)")
	vendorCSRCName := vendorCMD.String("cn", "", "common name for certificate request")
	vendorPKeyPass := vendorCMD.String("password", "", "rsa private key password")
	vendorSignFlag := vendorCMD.Bool("sign", false, "sign a provider push csr with the vendor certificate")
	vendorCertPath := vendorCMD.String("cert", "mdm.cer", "path to mdm vendor cert provided by apple")
	vendorProviderCSRPath := vendorCMD.String("provider-csr", providerCSRFilename, "path to csr which needs to be signed")
	vendorPKeyPath := vendorCMD.String("private-key", vendorPKeyFilename, "path to provider csr which needs to be signed")

	// provider cmd flags
	providerCMD := flag.NewFlagSet("provider", flag.ExitOnError)
	providerCSRFlag := providerCMD.Bool("csr", false, "create a CSR for a push certificate request")
	providerCSREmail := providerCMD.String("email", "", "email address to use in CSR request Subject")
	providerCSRCountry := providerCMD.String("country", "US", "two letter country flag for CSR Subject(example: US)")
	providerCSRCName := providerCMD.String("cn", "", "common name for certificate request")
	providerPKeyPass := providerCMD.String("password", "", "rsa private key password")
	// general flags
	flVersion := flag.Bool("version", false, "prints the version")
	// set usage
	flag.Usage = func() {
		fmt.Println("usage: certhelper <command> [<args>]")
		fmt.Println(" vendor <args> manage mdm vendor certs")
		fmt.Println(" provider <args> manage certs as a provider(mdm server administrator)")
		fmt.Println("type <command> --help to see usage for each subcommand")
	}

	flag.Parse()

	if *flVersion {
		fmt.Printf("certhelper - %v\n", version)
		fmt.Printf("git revision - %v\n", gitHash)
		os.Exit(0)
	}

	if len(os.Args) <= 2 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "vendor":
		vendorCMD.Parse(os.Args[2:])
	case "provider":
		providerCMD.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if vendorCMD.Parsed() {
		password := []byte(*vendorPKeyPass)
		if *vendorCSRFlag {
			// validate CSR arguments
			if err := checkCSRFlags(*vendorCSRCName, *vendorCSRCountry, *vendorCSREmail, password); err != nil {
				fmt.Println("private key password, cn, email, and country code must be provided for CSR")
				fmt.Printf("err: %s\n", err)
				os.Exit(1)
			}
			// create CSR
			req := &csrRequest{
				cname:        *vendorCSRCName,
				email:        *vendorCSREmail,
				country:      *vendorCSRCountry,
				password:     password,
				pkeyFilename: vendorPKeyFilename,
				csrFilename:  vendorCSRFilename,
			}
			if err := makeCSR(req); err != nil {
				fmt.Printf("err: %s\n", err)
				os.Exit(1)
			}
		}
		// sign a csr request
		if *vendorSignFlag {
			pushRequest, err := makePushRequestPlist(
				*vendorCertPath,
				*vendorProviderCSRPath,
				*vendorPKeyPath,
				password,
			)
			if err != nil {
				fmt.Printf("err: %s\n", err)
				os.Exit(1)
			}
			if err := writePushCertRequest(pushRequest); err != nil {
				if err != nil {
					fmt.Printf("err: %s\n", err)
					os.Exit(1)
				}
			}
		}
	}

	if providerCMD.Parsed() {
		password := []byte(*providerPKeyPass)
		if *providerCSRFlag {
			// validate CSR arguments
			if err := checkCSRFlags(*providerCSRCName, *providerCSRCountry, *providerCSREmail, password); err != nil {
				fmt.Println("private key password, cn, email, and country code must be provided for CSR")
				fmt.Printf("err: %s\n", err)
				os.Exit(1)
			}
			// create CSR
			req := &csrRequest{
				cname:        *providerCSRCName,
				email:        *providerCSREmail,
				country:      *providerCSRCountry,
				password:     password,
				pkeyFilename: providerPKeyFilename,
				csrFilename:  providerCSRFilename,
			}
			if err := makeCSR(req); err != nil {
				fmt.Printf("err: %s\n", err)
				os.Exit(1)
			}
		}
	}
}

func checkCSRFlags(cname, country, email string, password []byte) error {
	if cname == "" {
		return errors.New("cn flag not specified")
	}
	if email == "" {
		return errors.New("email flag not specified")
	}
	if country == "" {
		return errors.New("country flag not specified")
	}
	if len(password) == 0 {
		return errors.New("private key password empty")
	}
	if len(country) != 2 {
		return errors.New("must be a two letter country code")
	}
	return nil
}

// plist for push certificate request
type pushCertRequest struct {
	PushCertRequestCSR       string
	PushCertCertificateChain string
	PushCertSignature        string
}

// args for a csr request
type csrRequest struct {
	cname, country, email     string
	password                  []byte
	pkeyFilename, csrFilename string
}

// create a private key and CSR and save both to disk
func makeCSR(req *csrRequest) error {
	key, err := newRSAKey(2048)
	if err != nil {
		return err
	}
	pemKey, err := encryptedKey(key, req.password)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(req.pkeyFilename, pemKey, 0600); err != nil {
		return err
	}

	derBytes, err := newCSR(key, strings.ToLower(req.email), strings.ToUpper(req.country), req.cname)
	if err != nil {
		return err
	}
	pemCSR := pemCSR(derBytes)
	return ioutil.WriteFile(req.csrFilename, pemCSR, 0600)
}

// create a push request plist
func makePushRequestPlist(mdmCertPath, providerCSRPath, pKeyPath string, pKeyPass []byte) (*pushCertRequest, error) {
	// private key of the mdm vendor cert
	key, err := loadKeyFromFile(pKeyPath, pKeyPass)
	if err != nil {
		return nil, err
	}

	// provider csr
	csr, err := loadCSRfromFile(providerCSRPath)
	if err != nil {
		return nil, err
	}

	// csr signature
	signature, err := signProviderCSR(csr.Raw, key)
	if err != nil {
		return nil, err
	}

	// vendor cert
	mdmCertBytes, err := loadDERCertFromFile(mdmCertPath)
	if err != nil {
		return nil, err
	}
	mdmPEM := pemCert(mdmCertBytes)

	// wwdr cert
	wwdrCertBytes, err := loadCertfromHTTP(wwdrIntermediaryURL)
	if err != nil {
		return nil, err
	}
	wwdrPEM := pemCert(wwdrCertBytes)

	// apple root certificate
	rootCertBytes, err := loadCertfromHTTP(appleRootCAURL)
	if err != nil {
		return nil, err
	}
	rootPEM := pemCert(rootCertBytes)

	csrB64 := base64.StdEncoding.EncodeToString(csr.Raw)
	sig64 := base64.StdEncoding.EncodeToString(signature)
	pushReq := &pushCertRequest{
		PushCertRequestCSR:       csrB64,
		PushCertCertificateChain: makeCertChain(mdmPEM, wwdrPEM, rootPEM),
		PushCertSignature:        sig64,
	}
	return pushReq, nil
}

// save plist as base64 encoded string
func writePushCertRequest(req *pushCertRequest) error {
	data, err := plist.MarshalIndent(req, "  ")
	if err != nil {
		return err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	if err := ioutil.WriteFile(pushRequestFilename, encoded, 0600); err != nil {
		return err
	}
	return nil
}

func makeCertChain(mdmPEM, wwdrPEM, rootPEM []byte) string {
	return string(mdmPEM) + string(wwdrPEM) + string(rootPEM)
}

func signProviderCSR(csrData []byte, key *rsa.PrivateKey) ([]byte, error) {
	h := sha1.New()
	h.Write(csrData)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, h.Sum(nil))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return nil, err
	}
	return signature, nil
}
