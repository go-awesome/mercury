// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webpush

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"errors"
	"net/url"
	"encoding/base64"
)

const publicKeyLength = 65;
const privateKeyLength = 32;

const gcmEndpoint = "https://android.googleapis.com/gcm/send"

type vapid struct {
	subject    string
	privateKey string
	publicKey  string
}

type Push struct {
	currentGcmApiKey string
	currentVapid 	 *vapid
}

func (p *Push) SetGcmApiKey(gcmApiKey string) {
	p.currentGcmApiKey = gcmApiKey
}

func (p *Push) SetVAPID(subject, privateKey, publicKey string) {
	p.currentVapid = &vapid{subject: subject, privateKey: privateKey, publicKey: publicKey}
}

func (p *Push) Do(client *http.Client, sub *Subscription, message string) (*http.Response, error) {
	req, err := http.NewRequest("POST", sub.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("TTL", "604800")

	var cryptoHeaderKey string

	if message != "" {
		payload, err := Encrypt(sub, message)
		if err != nil {
			return nil, err
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(payload.Ciphertext))
		req.ContentLength = int64(len(payload.Ciphertext))
		req.Header.Add("Encryption", headerField("salt", payload.Salt))
		req.Header.Add("Content-Encoding", "aesgcm")

		cryptoHeaderKey = headerField("dh", payload.ServerPublicKey)
	}

	isGcm := strings.HasPrefix(sub.Endpoint, gcmEndpoint)
	if isGcm {
		if p.currentGcmApiKey == "" {
			return nil, errors.New("The GCM API Key should be a non-empty string")
		}
		req.Header.Add("Authorization", fmt.Sprintf(`key=%s`, p.currentGcmApiKey))

	} else if p.currentVapid != nil {

		subEndpointURL, err := url.Parse(sub.Endpoint)
		if err != nil { return nil, err }

		audience := subEndpointURL.Scheme + "://" + subEndpointURL.Host

		auth, criptoKey, err := getVapidHeaders(audience, p.currentVapid.subject, p.currentVapid.privateKey, p.currentVapid.publicKey)
		if err != nil { return nil, err }

		req.Header.Add("Authorization", auth)

		if cryptoHeaderKey != "" {
			cryptoHeaderKey = cryptoHeaderKey + ";" + criptoKey
		} else {
			cryptoHeaderKey = criptoKey
		}
	}
	req.Header.Add("Crypto-Key", cryptoHeaderKey)

	return client.Do(req)
}

func getVapidHeaders(audience string, subject string, privateKey string, publicKey string) (string, string, error) {

	if subject == "" { return "", "", errors.New("vapid: you must provide a subject that is either a mailto: or a URL")}

	b64 := base64.URLEncoding.WithPadding(base64.NoPadding)

	privateKeyBytes, err := b64.DecodeString(privateKey)
	if err != nil { return "", "", err }

	publicKeyBytes, err := b64.DecodeString(publicKey)
	if err != nil { return "", "", err }

	if len(privateKeyBytes) != privateKeyLength { return "", "", errors.New("push: private key should be 32 bytes long when decoded")}
	if len(publicKeyBytes)  != publicKeyLength  { return "", "", errors.New("push: public key should be 65 bytes long when decoded")}

	vapid := NewVapid(privateKeyBytes)
	vapid.Sub = subject

	token, err := vapid.Token(audience)
	if err != nil { return "", "", err }

	auth      := "WebPush " + token
	cryptoKey := "p256ecdsa=" + publicKey

	return auth, cryptoKey, nil
}

// A helper for creating the value part of the HTTP encryption headers
func headerField(headerType string, value []byte) string {
	return fmt.Sprintf(`%s=%s`, headerType, strings.TrimRight(base64.URLEncoding.EncodeToString(value), "="))
}
