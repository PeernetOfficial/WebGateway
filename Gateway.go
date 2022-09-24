/*
File Name:  Gateway.go
Copyright:  2022 Peernet s.r.o.
Author:     Peter Kleissner
*/

package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/PeernetOfficial/core"
	"github.com/gorilla/mux"
)

func startWebGateway(backend *core.Backend) {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.HandlerFunc(webGatewayHandler)).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(config.WebFiles)))).Methods("GET")

	for _, listen := range config.WebListen {
		go startWebServer(backend, listen, config.WebUseSSL, config.WebCertificateFile, config.WebCertificateKey, router, "Web Listen", parseDuration(config.WebTimeoutRead), parseDuration(config.WebTimeoutWrite))
	}

	if config.Redirect80 != "" {
		go webRedirect80(config.Redirect80)
	}

	// wait forever
	select {}
}

// startWebServer starts the web-server and may block forever
func startWebServer(backend *core.Backend, WebListen string, UseSSL bool, CertificateFile, CertificateKey string, Handler http.Handler, Info string, ReadTimeout, WriteTimeout time.Duration) {
	//func startWebServer(backend *core.Backend, webListen string, useSSL bool, certificateFile, certificateKey string, server *http.Server) {
	backend.LogError("startWebServer", "Web Gateway to listen on '%s'", WebListen)

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12} // for security reasons disable TLS 1.0/1.1

	server := &http.Server{
		Addr:         WebListen,
		Handler:      Handler,
		ReadTimeout:  ReadTimeout,  // ReadTimeout is the maximum duration for reading the entire request, including the body.
		WriteTimeout: WriteTimeout, // WriteTimeout is the maximum duration before timing out writes of the response. This includes processing time and is therefore the max time any HTTP function may take.
		//IdleTimeout:  IdleTimeout,  // IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
		TLSConfig: tlsConfig,
	}

	if UseSSL {
		// HTTPS
		if err := server.ListenAndServeTLS(CertificateFile, CertificateKey); err != nil {
			backend.LogError("startWebServer", "Error listening on '%s': %v\n", WebListen, err)
		}
	} else {
		// HTTP
		if err := server.ListenAndServe(); err != nil {
			backend.LogError("startWebServer", "Error listening on '%s': %v\n", WebListen, err)
		}
	}
}

// parseDuration is the same as time.ParseDuration without returning an error. Valid units are ms, s, m, h. For example "10s".
func parseDuration(input string) (result time.Duration) {
	result, _ = time.ParseDuration(input)
	return
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

func webRedirect80(listen80 string) {
	// redirect HTTP -> HTTPS
	http.ListenAndServe(net.JoinHostPort(listen80, "80"), http.HandlerFunc(redirect))
}

func webGatewayHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}
