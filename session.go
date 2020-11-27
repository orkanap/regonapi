package regonapi

import (
	"encoding/xml"
	"errors"
	"regexp"
)

// Login starts a new session
func (c *Client) Login() error {
	params := struct {
		XMLName xml.Name `xml:"ns:Zaloguj"`
		Key     string   `xml:"ns:pKluczUzytkownika"`
	}{
		Key: c.key,
	}

	body, err := xml.MarshalIndent(params, "", "    ")
	if err != nil {
		return err
	}

	b, err := c.call(publicEnvelope, "Zaloguj", string(body))
	if err != nil {
		return err
	}

	// MTOM/XOP encoded, use regex
	r := regexp.MustCompile("<ZalogujResult>(.*)</ZalogujResult>")
	s := r.FindStringSubmatch(string(b))
	if len(s) == 0 {
		return ErrInvalidKey
	}

	// First matching group
	c.sid = s[1]

	return nil
}

// Logout ends session
func (c *Client) Logout() error {
	if c.sid == "" {
		return ErrSessionNotStarted
	}

	params := struct {
		XMLName xml.Name `xml:"ns:Wyloguj"`
		SID     string   `xml:"ns:pIdentyfikatorSesji"`
	}{
		SID: c.sid,
	}

	body, err := xml.MarshalIndent(params, "", "    ")
	if err != nil {
		return err
	}

	b, err := c.call(publicEnvelope, "Wyloguj", string(body))
	if err != nil {
		return err
	}

	// MTOM/XOP encoded, use regex
	r := regexp.MustCompile("<WylogujResult>(.*)</WylogujResult>")
	s := r.FindStringSubmatch(string(b))
	if len(s) == 0 {
		return errors.New("session not active")
	}

	return nil
}

func (c *Client) getValue(paramName string) (string, error) {
	if c.sid == "" {
		return "", ErrSessionNotStarted
	}

	params := struct {
		XMLName   xml.Name `xml:"ns:GetValue"`
		ParamName string   `xml:"ns:pNazwaParametru"`
	}{
		ParamName: paramName,
	}

	body, err := xml.MarshalIndent(params, "", "    ")
	if err != nil {
		return "", err
	}

	b, err := c.call(privateEnvelope, "GetValue", string(body))
	if err != nil {
		return "", err
	}

	// MTOM/XOP encoded, use regex-a.
	r := regexp.MustCompile("<GetValueResult>(.*)</GetValueResult>")
	s := r.FindStringSubmatch(string(b))
	if len(s) == 0 {
		return "", ErrEmptyResult
	}

	return s[1], nil
}

// SessionStatus returns current session status: 1 = session active, 0 = session
// no longer active
func (c *Client) SessionStatus() (string, error) {
	return c.getValue("StatusSesji")
}
