package regonapi

import (
	"encoding/xml"
	"fmt"
	"html"
	"regexp"
)

// Report returns report of type report name for entity identified by REGON
func (c *Client) report(regon, reportName string) (string, error) {
	if c.sid == "" {
		return "", ErrSessionNotStarted
	}

	params := struct {
		XMLName    xml.Name `xml:"ns:DanePobierzPelnyRaport"`
		REGON      string   `xml:"ns:pRegon"`
		ReportName string   `xml:"ns:pNazwaRaportu"`
	}{
		REGON:      regon,
		ReportName: reportName,
	}

	body, err := xml.MarshalIndent(params, "", "    ")
	if err != nil {
		return "", err
	}

	b, err := c.call(publicEnvelope, "DanePobierzPelnyRaport", string(body))
	if err != nil {
		return "", err
	}

	// MTOM/XOP encoded, use regex-a.
	r := regexp.MustCompile("<DanePobierzPelnyRaportResult>((.|\\s)*)</DanePobierzPelnyRaportResult>")
	s := r.FindStringSubmatch(html.UnescapeString(b))
	if len(s) == 0 {
		return "", ErrEmptyResult
	}

	return s[1], nil
}

// NaturalPerson holds details of entities type (F)
type NaturalPerson struct {
	ErrorCode             int    `xml:"ErrorCode"`
	REGON9                string `xml:"fiz_regon9"`
	NIP                   string `xml:"fiz_nip"`
	StatusNIP             string `xml:"fiz_statusNip"` // "", Uchylony, Unieważniony
	LastName              string `xml:"fiz_nazwisko"`
	FirstName1            string `xml:"fiz_imie1"`
	FirstName2            string `xml:"fiz_imie2"`
	REGONRegistrationDate string `xml:"fiz_dataWpisuDoREGON"`
	LastUpdateDate        string `xml:"fiz_dataZaistnieniaZmiany"`
	RemovalDate           string `xml:"fiz_dataSkresleniaPodmiotuZRegon"`
	BasicLegalFormCode    string `xml:"fiz_podstawowaFormaPrawna_Symbol"`
	SpecificLegalFormCode string `xml:"fiz_szczegolnaFormaPrawna_Symbol"`
	OwnershipFormCode     string `xml:"fiz_formaWlasnosci_Symbol"`
	BasicLegalFormName    string `xml:"fiz_podstawowaFormaPrawna_Nazwa"`
	SpecificLegalFormName string `xml:"fiz_szczegolnaFormaPrawna_Nazwa"`
	OwnershipFormName     string `xml:"fiz_formaWlasnosci_Nazwa"`
	NumberOfLocalUnits    int    `xml:"fiz_liczbaJednLokalnych"`
}

// NaturalPersonDetails returns details of natural person running economic
// activity (type F) or nil and error. Entity is identified by REGON.
func (c *Client) NaturalPersonDetails(regon string) (*NaturalPerson, error) {
	rap, err := c.report(regon, "BIR11OsFizycznaDaneOgolne")
	if err != nil {
		return nil, err
	}

	var res struct {
		XMLName struct{}      `xml:"root"`
		Rep     NaturalPerson `xml:"dane"`
	}

	err = xml.Unmarshal([]byte(rap), &res)
	if err != nil {
		return nil, err
	}
	if res.Rep.ErrorCode == 4 {
		return nil, ErrNoDataFound
	}
	if res.Rep.ErrorCode > 0 {
		return nil, fmt.Errorf("report error %d", res.Rep.ErrorCode)
	}

	return &res.Rep, nil
}

// LegalPerson holds details of entities type (P)
type LegalPerson struct {
	ErrorCode             int    `xml:"ErrorCode"`
	REGON9                string `xml:"praw_regon9"`
	NIP                   string `xml:"praw_nip"`
	StatusNIP             string `xml:"praw_statusNip"` // "", Uchylony, Unieważniony
	Name                  string `xml:"praw_nazwa"`
	ShortName             string `xml:"praw_nazwaSkrocona"`
	RegistrationNumberReg string `xml:"praw_numerWRejestrzeEwidencji"`
	RegistrationDateReg   string `xml:"praw_dataWpisuDoRejestruEwidencji"`
	CreationDate          string `xml:"praw_dataPowstania"`
	StartDate             string `xml:"praw_dataRozpoczeciaDzialalnosci"`
	RegistrationDate      string `xml:"praw_dataWpisuDoRegon"`
	HoldDate              string `xml:"praw_dataZawieszeniaDzialalnosci"`
	RenevalDate           string `xml:"praw_dataWznowieniaDzialalnosci"`
	RemovalDate           string `xml:"praw_dataSkresleniaZRegon"`
	LastUpdateDate        string `xml:"praw_dataZaistnieniaZmiany"`
	EndDate               string `xml:"praw_dataZakonczeniaDzialalnosci"`
	Phone                 string `xml:"praw_numerTelefonu"`
	ExtPhone              string `xml:"praw_numerWewnetrznyTelefonu"`
	Fax                   string `xml:"praw_numerFaksu"`
	Email                 string `xml:"praw_adresEmail"`
	WWW                   string `xml:"praw_adresStronyinternetowej"`
	NumberOfLocalUnits    int    `xml:"praw_liczbaJednLokalnych"`
}

// LegalPersonDetails returns details of legal person (type P) or nil and error.
// Entity is identified by REGON.
func (c *Client) LegalPersonDetails(regon string) (*LegalPerson, error) {
	rap, err := c.report(regon, "BIR11OsPrawna")
	if err != nil {
		return nil, err
	}

	var res struct {
		XMLName struct{}    `xml:"root"`
		Rep     LegalPerson `xml:"dane"`
	}

	err = xml.Unmarshal([]byte(rap), &res)
	if err != nil {
		return nil, err
	}
	if res.Rep.ErrorCode == 4 {
		return nil, ErrNoDataFound
	}
	if res.Rep.ErrorCode > 0 {
		return nil, fmt.Errorf("report error %d", res.Rep.ErrorCode)
	}

	return &res.Rep, nil
}

// NaturalPersonPKD holds information of PKD (classification of business activity)
// for natural person
type NaturalPersonPKD struct {
	ErrorCode   int    `xml:"ErrorCode"`
	Code        string `xml:"fiz_pkd_Kod"`
	Name        string `xml:"fiz_pkd_Nazwa"`
	Primary     string `xml:"fiz_pkd_Przewazajace"`
	SilosID     int    `xml:"fiz_SilosID"`
	SilosCode   string `xml:"fiz_Silos_Symbol"`
	RemovalDate string `xml:"fiz_dataSkresleniaDzialalnosciZRegon"`
}

// NaturalPersonPKDList  returns list of PKD for natural person or nil and error.
// Entity is identified by REGON.
func (c *Client) NaturalPersonPKDList(regon string) ([]NaturalPersonPKD, error) {
	rap, err := c.report(regon, "BIR11OsFizycznaPkd")
	if err != nil {
		return nil, err
	}

	var res struct {
		XMLName struct{}           `xml:"root"`
		Rep     []NaturalPersonPKD `xml:"dane"`
	}

	err = xml.Unmarshal([]byte(rap), &res)
	if err != nil {
		return nil, err
	}
	if len(res.Rep) > 0 && res.Rep[0].ErrorCode > 0 {
		errorCode := res.Rep[0].ErrorCode
		if errorCode == 4 {
			return nil, ErrNoDataFound
		}
		return nil, fmt.Errorf("report error %d", errorCode)
	}

	return res.Rep, nil
}

// LegalPersonPKD holds information of PKD (classification of business activity)
// for legal person
type LegalPersonPKD struct {
	ErrorCode int    `xml:"ErrorCode"`
	Code      string `xml:"praw_pkdKod"`
	Name      string `xml:"praw_pkdNazwa"`
	Primary   string `xml:"praw_pkdPrzewazajace"`
}

// LegalPersonPKDList returns list of PKD for legal person or nil and error.
// Entity is identified by REGON.
func (c *Client) LegalPersonPKDList(regon string) ([]LegalPersonPKD, error) {
	rap, err := c.report(regon, "BIR11OsPrawnaPkd")
	if err != nil {
		return nil, err
	}

	var res struct {
		XMLName struct{}         `xml:"root"`
		Rep     []LegalPersonPKD `xml:"dane"`
	}

	err = xml.Unmarshal([]byte(rap), &res)
	if err != nil {
		return nil, err
	}
	if len(res.Rep) > 0 && res.Rep[0].ErrorCode > 0 {
		errorCode := res.Rep[0].ErrorCode
		if errorCode == 4 {
			return nil, ErrNoDataFound
		}
		return nil, fmt.Errorf("report error %d", errorCode)
	}

	return res.Rep, nil
}
