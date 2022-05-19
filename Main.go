/*
File Name:  Main.go
Copyright:  2021 Peernet s.r.o.
Author:     Peter Kleissner
*/

package main

import (
	"fmt"
	"os"

	"github.com/PeernetOfficial/core"
)

const configFile = "Config.yaml"
const appName = "Peernet Web Gateway"

var config struct {
	// HTTP settings for the web gateway
	WebListen          []string `yaml:"WebListen"`          // WebListen is in format IP:Port and declares where the web-interface should listen on. IP can also be ommitted to listen on any.
	WebUseSSL          bool     `yaml:"WebUseSSL"`          // Enables SSL.
	WebCertificateFile string   `yaml:"WebCertificateFile"` // This is the certificate received from the CA. This can also include the intermediate certificate from the CA.
	WebCertificateKey  string   `yaml:"WebCertificateKey"`  // This is the private key.
	WebTimeoutRead     string   `yaml:"WebTimeoutRead"`     // The maximum duration for reading the entire request, including the body.
	WebTimeoutWrite    string   `yaml:"WebTimeoutWrite"`    // The maximum duration before timing out writes of the response. This includes processing time and is therefore the max time any HTTP function may take.
}

func main() {
	userAgent := appName + "/" + core.Version

	backend, status, err := core.Init(userAgent, configFile, nil, &config)

	if status != core.ExitSuccess {
		switch status {
		case core.ExitErrorConfigAccess:
			fmt.Printf("Unknown error accessing config file '%s': %s\n", configFile, err.Error())
		case core.ExitErrorConfigRead:
			fmt.Printf("Error reading config file '%s': %s\n", configFile, err.Error())
		case core.ExitErrorConfigParse:
			fmt.Printf("Error parsing config file '%s' (make sure it is valid YAML format): %s\n", configFile, err.Error())
		case core.ExitErrorLogInit:
			fmt.Printf("Error opening log file '%s': %s\n", backend.Config.LogFile, err.Error())
		default:
			fmt.Printf("Unknown error %d initializing backend: %s\n", status, err.Error())
		}
		os.Exit(status)
	}

	backend.Stdout.Subscribe(os.Stdout)

	//startWebGateway(backend, )

	backend.Connect()
}
