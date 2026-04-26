package httpclient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// New returns a shared-style HTTP client with bounded timeouts and TLS 1.2+.
func New() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 2 * time.Minute,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return &http.Client{
		Transport: transport,
		Timeout:   3 * time.Minute, // end-to-end (connect + full body read for long streams)
	}
}

// Default is the process-wide default client (timeouts/TLS in New()).
var Default = New()
