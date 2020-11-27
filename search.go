package regonapi

import (
	"encoding/xml"
	"fmt"
	"html"
	"regexp"
)

// Entity record in the REGON database describes business entity
type Entity struct {
	ErrorCode       int    `xml:"ErrorCode"`
	REGON           string `xml:"Regon"`
	NIP             string `xml:"Nip"`
	StatusNIP       string `xml:"StatusNip"` // "", Uchylony, Unieważniony
	Name            string `xml:"Nazwa"`
	Province        string `xml:"Wojewodztwo"`
	District        string `xml:"Powiat"`
	Community       string `xml:"Gmina"`
	City            string `xml:"Miejscowosc"`
	PostalCode      string `xml:"KodPocztowy"`
	Street          string `xml:"Ulica"`
	PropertyNumber  string `xml:"NrNieruchomosci"`
	ApartmentNumber string `xml:"NrLokalu"`

	// (P) legal entity (F) natural person running economic activity (LP) local
	// unit of the legal entity (LF) local unit of a natural person
	Type            string `xml:"Typ"`
	SilosID         string `xml:"SilosID"`
	ActivityEndDate string `xml:"DataZakonczeniaDzialalnosci"`
	PostCity        string `xml:"MiejscowoscPoczty"`
}

// Search searches the REGON database for records that match the specified
// search criteria. Returns list of entities or nil and error.
func (c *Client) search(krs, nip, regon string) ([]Entity, error) {
	if c.sid == "" {
		return nil, ErrSessionNotStarted
	}

	params := struct {
		XMLName xml.Name `xml:"ns:DaneSzukajPodmioty"`
		KRS     string   `xml:"ns:pParametryWyszukiwania>dat:Krs,omitempty"`
		NIP     string   `xml:"ns:pParametryWyszukiwania>dat:Nip,omitempty"`
		REGON   string   `xml:"ns:pParametryWyszukiwania>dat:Regon,omitempty"`
	}{
		KRS:   krs,
		NIP:   nip,
		REGON: regon,
	}

	body, err := xml.MarshalIndent(params, "", "    ")
	if err != nil {
		return nil, err
	}

	b, err := c.call(publicEnvelope, "DaneSzukajPodmioty", string(body))
	if err != nil {
		return nil, err
	}

	// MTOM/XOP encoded, use regex-a.
	r := regexp.MustCompile("<DaneSzukajPodmiotyResult>((.|\\s)*)</DaneSzukajPodmiotyResult>")
	s := r.FindStringSubmatch(html.UnescapeString(b))
	if len(s) == 0 {
		return nil, ErrEmptyResult
	}

	var res struct {
		XMLName struct{} `xml:"root"`
		Data    []Entity `xml:"dane"`
	}

	err = xml.Unmarshal([]byte(s[1]), &res)
	if err != nil {
		return nil, err
	}
	if len(res.Data) > 0 && res.Data[0].ErrorCode > 0 {
		errorCode := res.Data[0].ErrorCode
		if errorCode == 4 {
			return nil, ErrNoDataFound
		}
		return nil, fmt.Errorf("search error %d", errorCode)
	}

	return res.Data, nil
}

// SearchByNIP searches the REGON database by NIP. Returns list of entities or
// nil and error. Input parameter 10 digits string.
func (c *Client) SearchByNIP(nip string) ([]Entity, error) {
	return c.search("", nip, "")
}

// SearchByREGON searches the database by REGON. Returns list of entities or nil
// and error. Input parameter must be normalized: 9 or 14 digits.
func (c *Client) SearchByREGON(regon string) ([]Entity, error) {
	return c.search("", "", regon)
}

// SearchByKRS searches the REGON database by KRS. Returns list of entities or
// nil and error. Input parameter 10 digits string.
func (c *Client) SearchByKRS(krs string) ([]Entity, error) {
	return c.search(krs, "", "")
}
