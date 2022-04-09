package network

import (
	"crypto/tls"
	"net"
	"net/http"
)

//NewCustomHttp Return custom http client - essentially unlimited timeouts and close connections after request done
func NewCustomHttp() (client *http.Client) {
	//Setting timeout custom transport layer
	tr := http.Transport{
		Proxy:       nil,
		DialContext: nil,
		DialTLSContext: (&net.Dialer{
			Timeout: 0,
		}).DialContext,
		TLSClientConfig:        &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout:    0,
		DisableKeepAlives:      true,
		DisableCompression:     false,
		MaxIdleConns:           0,
		MaxIdleConnsPerHost:    0,
		MaxConnsPerHost:        0,
		IdleConnTimeout:        0,
		ResponseHeaderTimeout:  0,
		ExpectContinueTimeout:  0,
		MaxResponseHeaderBytes: 0,
		ForceAttemptHTTP2:      false,
	}

	client = &http.Client{
		Transport: &tr,
		Timeout:   0,
	}
	return client
}
