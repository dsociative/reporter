package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type HTTPClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

type ClickhouseRequester struct {
	client HTTPClient
	addr   string
}

func NewClickhouseRequester(client HTTPClient, addr string) ClickhouseRequester {
	return ClickhouseRequester{
		client: client,
		addr:   addr,
	}
}

type Query struct {
	ClientID  string
	Network   string
	StartDate string
	StopDate  string
}

func (c ClickhouseRequester) Do(client, network, startDate, stopDate string) (string, error) {
	var data string
	var b *bytes.Buffer
	var err error

	if b, err = render(client, network, startDate, stopDate); err != nil {
		return data, err
	}

	resp, err := c.client.Post(c.addr, "text", b)
	if err == nil && resp.Body != nil {
		var r []byte
		r, err = ioutil.ReadAll(resp.Body)
		data = string(r)
	}
	return data, err
}
