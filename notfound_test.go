package custom404_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/theplant/custom404"

	"net/http/httptest"

	"io/ioutil"

	"strings"

	goji "goji.io"
	"goji.io/pat"
)

func TestMain(t *testing.T) {

	router := goji.NewMux()
	router.HandleFunc(pat.Get("/topics"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Topics")
	})

	mux := http.DefaultServeMux

	mux.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bar")
	})

	mux.Handle("/", router)

	newmux := custom404.WithCustom404(mux, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "My Not Found")
	})

	serv := httptest.NewServer(newmux)

	tests := [][]string{
		{"/topics", "200", "Topics"},
		{"/bar", "200", "Bar"},
		{"/abc", "404", "My Not Found"},
	}

	for _, ts := range tests {
		resp, err := http.Get(serv.URL + ts[0])
		if err != nil {
			panic(err)
		}
		bodyB, _ := ioutil.ReadAll(resp.Body)
		body := strings.TrimSpace(string(bodyB))
		if ts[2] != body {
			t.Error(body)
		}
		if ts[1] != fmt.Sprintf("%d", resp.StatusCode) {
			t.Error(resp.StatusCode)
		}

	}

}
