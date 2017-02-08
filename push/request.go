//
//  request.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package push

import (
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/oxtoacart/bpool"
)

var requestBuffPool = bpool.NewBufferPool(4096)

type Request struct {
	method 	string
	url 	string
	headers map[string]string

	requestEntity interface{}

	sendBuf *bytes.Buffer
	readBuf *bytes.Buffer

	client	*http.Client
}

func NewRequest(client *http.Client) *Request {
	r := &Request{}
	r.headers = make(map[string]string)

	r.sendBuf = requestBuffPool.Get()
	r.readBuf = requestBuffPool.Get()

	r.client = client
	return r
}

func (r *Request) Close() {
	requestBuffPool.Put(r.sendBuf)
	requestBuffPool.Put(r.readBuf)
}

func (r *Request) URL(url string) *Request {
	r.url = url
	return r
}

func (r *Request) GET() *Request {
	r.method = "GET"
	return r
}

func (r *Request) POST() *Request {
	r.method = "POST"
	return r
}

func (r *Request) PUT() *Request {
	r.method = "PUT"
	return r
}

func (r *Request) DELETE() *Request {
	r.method = "DELETE"
	return r
}

func (r *Request) Headers(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) SendEntity(requestEntity interface{}) *Request {
	r.requestEntity = requestEntity
	return r
}

func (r *Request) Do(responseEntity interface{}) (int, error) {
	var buf *bytes.Buffer = nil

	// marshall request entity
	if r.requestEntity != nil {
		buf = r.sendBuf
		if err := json.NewEncoder(buf).Encode(r.requestEntity); err != nil {
			return 500, err
		}
	}

	req, err := http.NewRequest(r.method, r.url, buf)
	if err != nil {
		return 500, err
	}

	// add request headers
	if r.headers != nil {
		for key, value := range r.headers {
			req.Header.Add(key, value)
		}
	}

	// perform
	resp, err := r.client.Do(req)
	if err != nil {
		return 500, err
	}
	defer resp.Body.Close()

	// unmarshall response entity
	if responseEntity != nil {
		_, err := r.readBuf.ReadFrom(resp.Body)
		if err != nil {
			return 500, err
		}

		if err := json.NewDecoder(r.readBuf).Decode(responseEntity); err != nil {
			return 500, err
		}
	}
	return resp.StatusCode, nil
}
