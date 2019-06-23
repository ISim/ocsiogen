package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/rs/xid"
)

func TestAccountCreate(t *testing.T) {
	cc := Account{
		RequestID: xid.New().String(),
		Name:      "INEC s.r.o",
	}
	apiKey := "huberokororo"
	apiID := "Rumburak"

	r := make(chan *http.Request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("content-type", "application/json")
		res := Response{
			RequestID: cc.RequestID,
			Data:      []ResponseData{},
		}
		res.Result.Status = "OK"
		r <- req
		json.NewEncoder(rw).Encode(res)

	}))
	defer server.Close()
	api, _ := NewClient(server.URL+"/v223", apiID, apiKey)

	res, err := api.CreateAccount(context.Background(), cc)

	req := <-r

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if tmp := res.Result.Status; tmp != "OK" {
		t.Errorf("expected OK, got '%s'", tmp)
	}

	expURL := url4actions[apiCreateCustomer]
	if !regexp.MustCompile(expURL + `$`).MatchString(req.URL.String()) {
		t.Errorf("unexpected URL calles '%s', expected something like `.../v223/%s",
			req.URL.String(), expURL)
	}
}
