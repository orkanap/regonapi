// Package regonapi provides interface for REGON API (Polish National Official
// Business Register). This is a thin wrapper for official SOAP webservice, BIR1
// version 1.1.
package regonapi

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	// ErrInvalidKey is returned by Login
	ErrInvalidKey = errors.New("invalid API key")

	// ErrSessionNotStarted is returned when session not started. Login to start
	// a new session.
	ErrSessionNotStarted = errors.New("session not started")

	// ErrNoDataFound is returned by search API when no data is found in the
	// REGON database
	ErrNoDataFound = errors.New("no data found for the specified search criteria")

	// ErrEmptyResult is returned when API call returns an empty body. Check
	// session status and if session has expired login again.
	ErrEmptyResult = errors.New("empty result")
)

// Client for REGON API webservice
type Client struct {
	endpoint string
	key      string
	sid      string
	ctx      context.Context
}

// NewClient returns new REGON API client. For an empty key client connects to
// test endpoint (with test key). Provide context to control the underlying HTTP
// requests.
func NewClient(ctx context.Context, key string) *Client {
	const (
		testEndpoint = "https://wyszukiwarkaregontest.stat.gov.pl/wsBIR/UslugaBIRzewnPubl.svc"
		testKey      = "abcde12345abcde12345"
		prodEndpoint = "https://wyszukiwarkaregon.stat.gov.pl/wsBIR/UslugaBIRzewnPubl.svc"
	)
	if ctx == nil {
		ctx = context.Background()
	}
	if key == "" {
		return &Client{
			endpoint: testEndpoint,
			key:      testKey,
			ctx:      ctx,
		}
	}
	return &Client{
		endpoint: prodEndpoint,
		key:      key,
		ctx:      ctx,
	}
}

const (
	publicEnvelope  = true
	privateEnvelope = false

	// An envelope template. It will be formatted to build SOAP request
	// envelope. Note (%s) verbs for XML namespace, action and body.
	envelope = `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" %s>
		<soap:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
			<wsa:To>https://wyszukiwarkaregontest.stat.gov.pl/wsBIR/UslugaBIRzewnPubl.svc</wsa:To>
			<wsa:Action>%s</wsa:Action>
		</soap:Header>
		<soap:Body>
		%s
		</soap:Body>
	</soap:Envelope>`
)

func buildEnvelope(public bool, action, body string) string {
	var ns, act string
	if public {
		ns = `xmlns:ns="http://CIS/BIR/PUBL/2014/07" xmlns:dat="http://CIS/BIR/PUBL/2014/07/DataContract"`
		act = `http://CIS/BIR/PUBL/2014/07/IUslugaBIRzewnPubl/` + action
	} else {
		ns = `xmlns:ns="http://CIS/BIR/2014/07"`
		act = `http://CIS/BIR/2014/07/IUslugaBIR/` + action
	}
	return fmt.Sprintf(envelope, ns, act, body)
}

func (c *Client) call(public bool, action, body string) (string, error) {
	request := buildEnvelope(public, action, body)

	req, err := http.NewRequestWithContext(c.ctx, "POST", c.endpoint, strings.NewReader(request))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	if c.sid != "" {
		req.Header.Set("sid", c.sid)
	}

	ret, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer ret.Body.Close()

	b, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
