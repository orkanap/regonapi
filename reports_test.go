package regonapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestRapLegalPersonNotFound(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<DanePobierzPelnyRaportResult><root>
		<dane>
		  <ErrorCode>4</ErrorCode>
		  <ErrorMessagePl>Nie znaleziono wpisu dla podanych kryteriów wyszukiwania.</ErrorMessagePl>
		  <ErrorMessageEn>No data found for the specified search criteria.</ErrorMessageEn>
		  <pRegon>999999999</pRegon>
		  <Typ_podmiotu />
		  <Raport>BIR11OsFizycznaDaneOgolne</Raport>
		</dane>
	  </root></DanePobierzPelnyRaportResult>`)
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "mock-sid"

	lp, err := regonsvc.LegalPersonDetails("n/a")
	is.Equal(err, ErrNoDataFound)
	is.Equal(lp, nil)
}

func TestRapLegalPersonOtherErr(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<DanePobierzPelnyRaportResult><root>
		<dane>
          <ErrorCode>5</ErrorCode>
          <ErrorMessagePl>Nieprawidłowa lub pusta nazwa raportu.</ErrorMessagePl>
          <ErrorMessageEn>Invalid or empty report name.</ErrorMessageEn>
          <pRegon>99999999900099</pRegon>
          <Typ_podmiotu />
          <Raport> BIR1.1JednLokalnaOsFizycznej</Raport>
		</dane>
	  </root></DanePobierzPelnyRaportResult>`)
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "mock-sid"

	lp, err := regonsvc.LegalPersonDetails("n/a")
	is.True(err != nil)
	is.Equal(lp, nil)
}

func TestRapLegalPersonListErr(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<DanePobierzPelnyRaportResult><root>
		<dane>
          <ErrorCode>5</ErrorCode>
          <ErrorMessagePl>Nieprawidłowa lub pusta nazwa raportu.</ErrorMessagePl>
          <ErrorMessageEn>Invalid or empty report name.</ErrorMessageEn>
          <pRegon>99999999900099</pRegon>
          <Typ_podmiotu />
          <Raport> BIR1.1JednLokalnaOsFizycznej</Raport>
		</dane>
	  </root></DanePobierzPelnyRaportResult>`)
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "mock-sid"

	lp, err := regonsvc.LegalPersonPKDList("n/a")
	is.True(err != nil)
	is.Equal(lp, nil)
}
