# goparse

[![Build Status](https://circleci.com/gh/qurami/goparse.svg?style=shield)](https://circleci.com/gh/qurami/goparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/qurami/goparse)](https://goreportcard.com/badge/github.com/qurami/goparse)

Simple Parse API client written in pure Go.


## Installation

`go get github.com/qurami/goparse`


## Usage example

```
import (
	"io/ioutil"
	"log"
	"net/url"

	"github.com/qurami/goparse"
)

func main() {
    // init ParseClient
	c, err := goparse.NewParseClient(
		"http://your-self-hosted-parse-api-url",
		"YOUR-APPLICATION-ID",
		"YOUR-REST-API-KEY",
	)
	if err != nil {
		log.Fatal(err)
	}

    // set urlParameters, if necessary
	urlParams := url.Values{}
	urlParams.Add("limit", "1")

    // run request
	res, err := c.Get("/classes/ServiceSatisfactionReport?" + urlParams.Encode())
	if err != nil {
		log.Println(err)
	}

    // read response
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.Status)
	log.Println(string(b))
}
```

You can also set the master key, if necessary, by running `c.SetMasterKey` method before executing the request.
