package http

import (
	"crypto/tls"
	"net/http"
)

var secondLevelApiClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}
