// +build integration

package regonapi_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/orkanap/regonapi"
)

func TestSessionLogin(t *testing.T) {
	is := is.New(t)

	// login
	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// session should be active
	status, err := regonsvc.SessionStatus()
	is.NoErr(err)
	is.True(status == "1")

	err = regonsvc.Logout()
	is.NoErr(err)
}

func TestSessionInvalidKey(t *testing.T) {
	is := is.New(t)

	regonsvc := regonapi.NewClient(context.Background(), "invalid-key")
	err := regonsvc.Login()
	is.True(err == regonapi.ErrInvalidKey)
}

func TestSearch(t *testing.T) {
	is := is.New(t)

	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// entity does not exist, wrong NIP
	entities, err := regonsvc.SearchByNIP("9999999999")
	is.Equal(err, regonapi.ErrNoDataFound)
	is.Equal(entities, nil)

	// entity exists
	entities, err = regonsvc.SearchByKRS("0000359584")
	is.NoErr(err)
	is.Equal(len(entities), 1)

	// entity exists, check name returned
	entities, err = regonsvc.SearchByREGON("000331501")
	is.NoErr(err)
	is.Equal(len(entities), 1)
	is.Equal(entities[0].Name, "GŁÓWNY URZĄD STATYSTYCZNY")

	err = regonsvc.Logout()
	is.NoErr(err)
}

func TestSearchInactiveSession(t *testing.T) {
	is := is.New(t)

	// login
	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// logout
	err = regonsvc.Logout()
	is.NoErr(err)

	// session should be inactive
	_, err = regonsvc.SearchByKRS("0000359584")
	is.True(err == regonapi.ErrEmptyResult)
}

func TestRepNaturalPersonDetails(t *testing.T) {
	is := is.New(t)

	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// entity exists
	_, err = regonsvc.NaturalPersonDetails("092343955")
	is.NoErr(err)

	// entity not exists
	_, err = regonsvc.NaturalPersonDetails("123456789")
	is.True(err == regonapi.ErrNoDataFound)

	// entity exists but not a natural person
	_, err = regonsvc.NaturalPersonDetails("000331501")
	is.True(err == regonapi.ErrNoDataFound)

	err = regonsvc.Logout()
	is.NoErr(err)
}

func TestRepNaturalPersonPKDList(t *testing.T) {
	is := is.New(t)

	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// entity exists
	pkds, err := regonsvc.NaturalPersonPKDList("092343955")
	is.NoErr(err)
	is.True(len(pkds) > 0)

	// entity not exists
	_, err = regonsvc.NaturalPersonPKDList("123456789")
	is.True(err == regonapi.ErrNoDataFound)

	// entity exists but not a natural person
	//
	// NO ERR! This is different than NaturalPersonDetails("000331501") where
	// ErrNoDataFound is returned.
	pkds, err = regonsvc.NaturalPersonPKDList("000331501")
	is.NoErr(err)

	err = regonsvc.Logout()
	is.NoErr(err)
}

func TestRepLegalPersonDetails(t *testing.T) {
	is := is.New(t)

	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// entity exists
	rep, err := regonsvc.LegalPersonDetails("000331501")
	is.NoErr(err)
	is.Equal(rep.ShortName, "GUS")

	// entity not exists
	_, err = regonsvc.LegalPersonDetails("123456789")
	is.True(err == regonapi.ErrNoDataFound)

	err = regonsvc.Logout()
	is.NoErr(err)
}

func TestRepLegalPersonPKDList(t *testing.T) {
	is := is.New(t)

	regonsvc := regonapi.NewClient(context.Background(), "")
	err := regonsvc.Login()
	is.NoErr(err)

	// entity exists
	pkds, err := regonsvc.LegalPersonPKDList("000331501")
	is.NoErr(err)
	is.True(len(pkds) > 0)

	// entity not exists
	_, err = regonsvc.LegalPersonPKDList("123456789")
	is.True(err == regonapi.ErrNoDataFound)

	// entity exists but not a natural person
	//
	// NO ERR! This is different than NaturalPersonDetails("092343955") where
	// ErrNoDataFound is returned.
	pkds, err = regonsvc.LegalPersonPKDList("092343955")
	is.NoErr(err)

	err = regonsvc.Logout()
	is.NoErr(err)
}
