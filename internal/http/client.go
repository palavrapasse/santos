package http

import (
	"crypto/tls"
	"net/http"
)

// #nosec
var secondLevelApiClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}
