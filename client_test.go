package regonapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestCallOK(t *testing.T) {
	is := is.New(t)

	var sid string

	// create mock handler and server
	h := func(w http.ResponseWriter, r *http.Request) {
		sid = r.Header.Get("sid")
		fmt.Fprint(w, "test-response")
	}
	mocks := httptest.NewServer(http.HandlerFunc(h))
	defer mocks.Close()

	// new client and set endpoint to mock server
	regonsvc := NewClient(context.Background(), "")
	regonsvc.endpoint = mocks.URL
	regonsvc.sid = "test-sid"

	resp, err := regonsvc.call(publicEnvelope, "action", "body")

	is.NoErr(err)
	is.Equal(sid, "test-sid")
	is.Equal(resp, "test-response")
}
