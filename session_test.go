package regonapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestLoginSID(t *testing.T) {
	is := is.New(t)

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<ZalogujResult>resp-sid</ZalogujResult>")
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL

	err := regonsvc.Login()

	is.NoErr(err)
	is.Equal(regonsvc.sid, "resp-sid")
}
