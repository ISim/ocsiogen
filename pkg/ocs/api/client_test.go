package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/rs/xid"
)

func TestCall(t *testing.T) {
	type formater func(f string, args ...interface{})
	tss := []struct {
		contentType string
		respFile    string
		code        int
		test        func(formater, *http.Response, error)
	}{
		{
			contentType: "application/json",
			respFile:    "response_OK.json",
			code:        http.StatusOK,
			test: func(f formater, resp *http.Response, err error) {
				if err != nil {
					f("unexpected error for response_OK.json file: %s", err)
					return
				}
				if resp == nil {
					f("response is nil for response_OK.json file?")
				}
			},
		},
		{
			contentType: "application/json",
			respFile:    "response_error.json",
			code:        http.StatusBadRequest,
			test: func(f formater, resp *http.Response, err error) {
				if err == nil {
					f("expected error, but call passed")
					return
				}
				e2, ok := err.(APIError)
				if !ok {
					f("expected APIError type, got %T: %s", err, err)
					return
				}
				r := e2.APIResponse()
				if rid := "bd46ee65-bc91-42ed-ac37-3ae9f6ed9090"; r.RequestID != rid {
					f("expected request ID '%s', got '%s'", rid, r.RequestID)
				}
				if e2.code != http.StatusBadRequest {
					f("expected status %d, got %d", http.StatusBadRequest, e2.code)
				}
			},
		},
		{
			contentType: "application/json",
			respFile:    "response_unparsable_json.txt",
			code:        http.StatusBadRequest,
			test: func(f formater, resp *http.Response, err error) {
				if err == nil {
					f("expected error, but call passed")
					return
				}
				_, ok := err.(HTTPError)
				if !ok {
					f("expected HTTPError type, got %T: %s", err, err)
					return
				}
			},
		},
		{
			contentType: "text/plain",
			respFile:    "response_plain.txt",
			code:        http.StatusInternalServerError,
			test: func(f formater, resp *http.Response, err error) {
				if err == nil {
					f("expected error, but call passed")
					return
				}
				er2, ok := err.(HTTPError)
				if !ok {
					f("expected HTTPError type, got %T: %s", err, err)
					return
				}
				exp := "Fatal error No. 234"
				body := er2.Body()
				if body == nil {
					f("expected body, got nil")
					return
				}
				if exp != string(body) {
					f("expected body '%s', got '%s'", string(exp), string(er2.Body()))
				}
			},
		},
	}

	for i, ts := range tss {
		jf := testFilePath("test-data", ts.respFile)
		f, err := os.Open(jf)
		if err != nil {
			t.Fatalf("can't open test file: '%s'", jf)
		}
		rawResponse, _ := ioutil.ReadAll(f)
		f.Close()
		r := make(chan *http.Request, 1)
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add("content-type", "application/json")
			rw.WriteHeader(ts.code)
			rw.Write(rawResponse)
			r <- req
		}))
		apiID, apiKey := xid.New().String(), xid.New().String()
		path := xid.New().String()
		targetURL, _ := url.Parse(server.URL + "/" + path) // arfiticial URL
		api, _ := NewClient(targetURL.String(), apiID, apiKey)

		res, err := api.call(api.prepareRequest(context.Background(), "POST", targetURL, nil))
		server.Close()

		ts.test(func(f string, args ...interface{}) {
			allargs := []interface{}{i + 1}
			allargs = append(allargs, args...)
			t.Errorf("Test case %d: "+f, allargs...)
		}, res, err)
	}

}
