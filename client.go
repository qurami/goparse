package goparse

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	contentTypeHeader   = "Content-Type"
	masterKeyHeader     = "X-Parse-Master-Key"
	restAPIKeyHeader    = "X-Parse-REST-API-Key"
	applicationIDHeader = "X-Parse-Application-ID"
)

// ParseClient provides access to the Parse API.
type ParseClient struct {
	apiBaseURL    *url.URL
	applicationID string
	restAPIKey    string
	masterKey     string
}

// NewParseClient returns a new ParseClient pointing to the given
// rawAPIBaseURL, set with the applicationID and apiKey.
func NewParseClient(rawAPIBaseURL, applicationID, apiKey string) (*ParseClient, error) {
	apiBaseURL, err := url.ParseRequestURI(
		strings.TrimSuffix(
			rawAPIBaseURL,
			"/",
		),
	)
	if err != nil {
		return nil, err
	}

	client := ParseClient{
		apiBaseURL:    apiBaseURL,
		applicationID: applicationID,
		restAPIKey:    apiKey,
	}

	return &client, nil
}

// SetMasterKey sets a master key, where necessary.
func (c *ParseClient) SetMasterKey(masterKey string) {
	c.masterKey = masterKey
}

// Get performs an http GET method call on the given resource uri.
func (c *ParseClient) Get(resourceURI string) (*http.Response, error) {
	return c.do(http.MethodGet, resourceURI, nil)
}

// Post performs an http POST method call on the given resource uri with the given body.
func (c *ParseClient) Post(resourceURI string, body interface{}) (*http.Response, error) {
	return c.do(http.MethodPost, resourceURI, body)
}

// Put performs an http PUT method call on the given resource uri with the given body.
func (c *ParseClient) Put(resourceURI string, body interface{}) (*http.Response, error) {
	return c.do(http.MethodPut, resourceURI, body)
}

// Delete performs an http DELETE method call on the given resource uri and unmarshals
// response into result.
func (c *ParseClient) Delete(resourceURI string) (*http.Response, error) {
	return c.do(http.MethodDelete, resourceURI, nil)
}

// do build and executes the HTTP request for the given method, resource URI and body.
func (c *ParseClient) do(method string, resourceURI string, body interface{}) (*http.Response, error) {
	client := c.prepareHTTPClient()

	req, err := c.prepareHTTPRequest(method, resourceURI, body)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

// prepareHTTPClient builds and returns a new HTTP client according to the API
// base URL Scheme, using a TLS configuration if using HTTPS.
func (c *ParseClient) prepareHTTPClient() *http.Client {
	if c.apiBaseURL.Scheme == "https" {
		tlscfg := &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		tlscfg.InsecureSkipVerify = true

		tr := &http.Transport{
			TLSClientConfig: tlscfg,
		}

		return &http.Client{Transport: tr}
	}

	return &http.Client{}
}

// prepareHTTPRequest builds and returns a new HTTP request with given
// method, for the given resource URI, with the given request body.
func (c *ParseClient) prepareHTTPRequest(method string, resourceURI string, body interface{}) (*http.Request, error) {
	var req *http.Request
	var err error

	requestURL := fmt.Sprintf(
		"%s/%s",
		c.apiBaseURL,
		strings.TrimPrefix(
			resourceURI,
			"/",
		),
	)

	if body == nil {
		req, err = http.NewRequest(method, requestURL, nil)
		if err != nil {
			return nil, err
		}
	} else {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		requestBodyReader := bytes.NewReader(jsonBody)
		req, err = http.NewRequest(method, requestURL, requestBodyReader)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set(contentTypeHeader, "application/json")
	req.Header.Set(applicationIDHeader, c.applicationID)
	req.Header.Set(restAPIKeyHeader, c.restAPIKey)

	if c.masterKey != "" {
		req.Header.Set(masterKeyHeader, c.masterKey)
	}

	return req, nil
}
