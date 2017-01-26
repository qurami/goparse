package goparse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestThatNewParseClientErrorsWhenAWrongAPIBaseURLIsPassed(t *testing.T) {
	_, err := NewParseClient(
		"invalid/url", // not absolute url!
		"mockApplicationID",
		"mockApplicationKey",
	)

	if err == nil {
		t.Fatal("Error expected, got nil")
	}
}

func TestThatNewParseClientReturnsAValidParseClient(t *testing.T) {
	sut, err := NewParseClient(
		"http://fake.parse.host",
		"mockApplicationID",
		"mockApplicationKey",
	)

	if err != nil || sut == nil {
		t.Fatalf("Error expected to be nil, got %s", err.Error())
	}

	mockURL, _ := url.Parse("http://fake.parse.host")
	mockClient := ParseClient{
		apiBaseURL:    mockURL,
		applicationID: "mockApplicationID",
		restAPIKey:    "mockApplicationKey",
	}

	if !reflect.DeepEqual(*sut, mockClient) {
		t.Error("ParseClient not set as expected")
	}
}

func TestThatNewParseClientStripsLastSlashOfTheAPIBaseURLIfPresent(t *testing.T) {
	sut, _ := NewParseClient(
		"http://fake.parse.host/",
		"mockApplicationID",
		"mockApplicationKey",
	)

	if sut.apiBaseURL.String() != "http://fake.parse.host" {
		t.Error("Slash was not stripped away")
	}
}

func TestThatSetMasterKeySetsAMasterKeyOnAClient(t *testing.T) {
	sut, _ := NewParseClient(
		"http://fake.parse.host/",
		"mockApplicationID",
		"mockApplicationKey",
	)

	sut.SetMasterKey("mockMasterKey")

	if sut.masterKey != "mockMasterKey" {
		t.Fatal("masterKey not set properly")
	}
}

func TestThatAnHTTPClientIsConfiguredWithTLSConfigurationIfAPIBaseURLIsHTTPS(t *testing.T) {
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{ "fake": "json data here" }`)
	}))

	sut, _ := NewParseClient(
		testServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	res, _ := sut.do("GET", "/fakeResource", nil)
	if res.StatusCode != 200 {
		t.Fatal("Expected HTTP 200")
	}
}

func TestThatResourceURIIsProperlySet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{ "fake": "json data here" }`)
	}))

	sut, _ := NewParseClient(
		testServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	res, _ := sut.do("GET", "fakeResource", nil)
	if res.Request.URL.String() != testServer.URL+"/fakeResource" {
		t.Fatal("Wrong request URI set")
	}
}

func TestThatProperHTTPRequestHeadersAreSet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{ "fake": "json data here" }`)
	}))

	sut, _ := NewParseClient(
		testServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	sut.SetMasterKey("mockMasterKey")

	res, _ := sut.do("GET", "/fakeResource", nil)
	if res.Request.Header.Get(contentTypeHeader) != "application/json" {
		t.Error("wrong Content-Type header")
	}
	if res.Request.Header.Get(applicationIDHeader) != "mockApplicationID" {
		t.Errorf("wrong %v header", applicationIDHeader)
	}
	if res.Request.Header.Get(restAPIKeyHeader) != "mockApplicationKey" {
		t.Errorf("wrong %v header", restAPIKeyHeader)
	}
	if res.Request.Header.Get(masterKeyHeader) != "mockMasterKey" {
		t.Errorf("wrong %v header", masterKeyHeader)
	}
}

func TestThatGetExecutesAnHTTPGetRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{ "fake": "json data here" }`)
	}))

	sut, _ := NewParseClient(
		testServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	res, _ := sut.Get("fakeResource")
	if res.Request.Method != "GET" {
		t.Fatalf("Expected a GET HTTP request, got %v", res.Request.Method)
	}
}

func TestThatPostExecutesAnHTTPPostRequest(t *testing.T) {
	echoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		defer r.Body.Close()

		b, _ := ioutil.ReadAll(r.Body)
		w.Write(b)
	}))

	sut, _ := NewParseClient(
		echoServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	mockStruct := struct {
		Name string `json:"name"`
	}{
		"Gianfranco",
	}

	res, _ := sut.Post("fakeResource", mockStruct)
	if res.Request.Method != "POST" {
		t.Fatalf("Expected a POST HTTP request, got %v", res.Request.Method)
	}

	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	if string(b) != "{\"name\":\"Gianfranco\"}" {
		t.Fatal("Wrong POST data")
	}
}

func TestThatPutExecutesAnHTTPPutRequest(t *testing.T) {
	echoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		defer r.Body.Close()

		b, _ := ioutil.ReadAll(r.Body)
		w.Write(b)
	}))

	sut, _ := NewParseClient(
		echoServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	mockStruct := struct {
		Name string `json:"name"`
	}{
		"Gianfranco",
	}

	res, _ := sut.Put("fakeResource", mockStruct)
	if res.Request.Method != "PUT" {
		t.Fatalf("Expected a PUT HTTP request, got %v", res.Request.Method)
	}

	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	if string(b) != "{\"name\":\"Gianfranco\"}" {
		t.Fatal("Wrong PUT data")
	}
}

func TestThatDeleteExecutesAnHTTPDeleteRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{ "fake": "json data here" }`)
	}))

	sut, _ := NewParseClient(
		testServer.URL,
		"mockApplicationID",
		"mockApplicationKey",
	)

	res, _ := sut.Delete("fakeResource")
	if res.Request.Method != "DELETE" {
		t.Fatalf("Expected a DELETE HTTP request, got %v", res.Request.Method)
	}
}
