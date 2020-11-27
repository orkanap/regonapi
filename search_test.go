package regonapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestSearchOK(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<DaneSzukajPodmiotyResult>
		<root>
		<dane>
		  <Regon>xxxxxxxxx</Regon>
		  <Nip>nnnnnnnnnn</Nip>
		  <StatusNip />
		  <Nazwa>AAAAAAAA</Nazwa>
		  <Wojewodztwo>LUBELSKIE</Wojewodztwo>
		  <Powiat>kraśnicki</Powiat>
		  <Gmina>Kraśnik</Gmina>
		  <Miejscowosc>Kraśnik</Miejscowosc>
		  <KodPocztowy>23-200</KodPocztowy>
		  <Ulica>ul. Test-Wilcza</Ulica>
		  <NrNieruchomosci>yy</NrNieruchomosci>
		  <NrLokalu />
		  <Typ>F</Typ>
		  <SilosID>1</SilosID>
		  <DataZakonczeniaDzialalnosci />
		</dane>
		<dane>
		  <Regon>xxxxxxxxx</Regon>
		  <Nip>nnnnnnnnnn</Nip>
		  <StatusNip />
		  <Nazwa>GOSPODARSTWO ROLNE</Nazwa>
		  <Wojewodztwo>LUBELSKIE</Wojewodztwo>
		  <Powiat>kraśnicki</Powiat>
		  <Gmina>Zakrzówek</Gmina>
		  <Miejscowosc>Sulów</Miejscowosc>
		  <KodPocztowy>23-213</KodPocztowy>
		  <NrNieruchomosci>zz</NrNieruchomosci>
		  <NrLokalu />
		  <Typ>F</Typ>
		  <SilosID>2</SilosID>
		  <DataZakonczeniaDzialalnosci />
		</dane>
	  </root>
	  </DaneSzukajPodmiotyResult>`)
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "mock-sid"

	entites, err := regonsvc.SearchByNIP("n/a")

	is.NoErr(err)
	is.Equal(len(entites), 2)
	is.Equal(entites[0].Name, "AAAAAAAA")
	is.Equal(entites[1].Name, "GOSPODARSTWO ROLNE")
}

func TestSearchNotFound(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<DaneSzukajPodmiotyResult><root>
		<dane>
		  <ErrorCode>4</ErrorCode>
		  <ErrorMessagePl>Nie znaleziono podmiotu dla podanych kryteriów wyszukiwania.</ErrorMessagePl>
		  <ErrorMessageEn>No data found for the specified search criteria.</ErrorMessageEn>
		  <Nip>9999999999</Nip>
		</dane>
	  </root>></DaneSzukajPodmiotyResult>`)
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "mock-sid"

	entites, err := regonsvc.SearchByNIP("n/a")

	is.Equal(err, ErrNoDataFound)
	is.Equal(entites, nil)
}

func TestSearchOtherErr(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<DaneSzukajPodmiotyResult><root>
		<dane>
		  <ErrorCode>5</ErrorCode>
		  <ErrorMessagePl>AAA</ErrorMessagePl>
		  <ErrorMessageEn>BBB</ErrorMessageEn>
		  <Nip>9999999999</Nip>
		</dane>
	  </root>></DaneSzukajPodmiotyResult>`)
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "mock-sid"

	entites, err := regonsvc.SearchByNIP("n/a")

	is.True(err != nil)
	is.Equal(entites, nil)
}
