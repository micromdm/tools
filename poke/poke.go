package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RobotsAndPencils/buford/certificate"
	"github.com/RobotsAndPencils/buford/payload"
	"github.com/RobotsAndPencils/buford/push"
)

var (
	version = "unreleased"
	gitHash = "unknown"
)

func main() {
	// flags
	var (
		flMagic    = flag.String("magic", "", "pushmagic")
		flToken    = flag.String("token", "", "deviceToken")
		flVersion  = flag.Bool("version", false, "print version information")
		flPushCert = flag.String("push-cert", envString("MDM_PUSH_CERT", ""), "path to push certificate")
		flPushPass = flag.String("push-pass", envString("MDM_PUSH_PASS", ""), "push certificate password")
	)
	flag.Parse()

	if *flVersion {
		fmt.Printf("poke - %v\n", version)
		fmt.Printf("git revision - %v\n", gitHash)
		os.Exit(0)
	}

	// load apns cert
	cert, key, err := certificate.Load(*flPushCert, *flPushPass)
	if err != nil {
		log.Fatal(err)
	}

	// check the validity of the token
	if !push.IsDeviceTokenValid(*flToken) {
		log.Fatal("invalid token")
	}

	// create bufford client
	client, err := push.NewClient(certificate.TLS(cert, key))
	if err != nil {
		log.Fatal(err)
	}

	// push service
	service := push.Service{
		Client: client,
		Host:   push.Production,
	}

	expiration := time.Now().Add(5 * time.Minute) // expire in 5 minutes
	headers := &push.Headers{
		LowPriority: true,
		Expiration:  expiration,
	}

	// notification payload
	p := payload.MDM{Token: *flMagic}
	id, err := service.Push(*flToken, headers, p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("notification sent successfuly: id: %s", id)
}

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return false
}
