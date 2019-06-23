package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Client is a caller of OCS.io API
type Client struct {
	url    *url.URL
	appID  string
	appKey string
}

// NewClient is a contructor of OCS.io api client
func NewClient(baseurl, appID, appKey string) (*Client, error) {
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create api client (URL: %s)", baseurl)
	}
	return &Client{
		url:    u,
		appID:  appID,
		appKey: appKey,
	}, nil
}

func (c *Client) buildURL(action string) *url.URL {
	partURL, _ := url.Parse(url4actions[action])
	return c.url.ResolveReference(partURL)
}

// prepare request with proper header fields
func (c *Client) prepareRequest(ctx context.Context, method string, URL *url.URL, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, URL.String(), body)
	r.Header.Add("Content-type", "application/json; charset=utf-8")
	r.Header.Add(headerOCSAppID, c.appID)
	r.Header.Add(headerOCSAppKey, c.appKey)
	return r.WithContext(ctx)
}

func (c *Client) call(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Wrapf(err, "transport error, URL: %s", r.URL.String())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read http response body")
	}

	herr := HTTPError{
		url:         r.URL.String(),
		code:        resp.StatusCode,
		contentType: r.Header.Get(http.CanonicalHeaderKey("content-type")),
		body:        b,
	}
	if strings.HasPrefix(herr.contentType, "application/json") {
		apiResp := Response{}
		err = json.Unmarshal(b, &apiResp)
		if err == nil {
			return nil, APIError{
				HTTPError:   herr,
				apiResponse: apiResp,
			}
		}
	}

	return nil, herr
}

func (c *Client) parseAPIResponse(resp *http.Response) (*Response, error) {
	var apiResp Response
	defer resp.Body.Close()

	contentType := resp.Request.Header.Get(http.CanonicalHeaderKey("content-type"))

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read http response body")
	}
	err = json.Unmarshal(b, &apiResp)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing JSON failed")
	}

	if apiResp.Result.Status != apiResultStatusOK {
		return nil, APIError{
			apiResponse: apiResp,
			HTTPError: HTTPError{
				url:         resp.Request.URL.String(),
				code:        resp.StatusCode,
				contentType: contentType,
				body:        b,
			},
		}
	}
	return &apiResp, nil
}

type APIError struct {
	apiResponse Response
	HTTPError
}

func (a APIError) APIResponse() Response {
	return a.apiResponse
}

type HTTPError struct {
	url         string
	code        int
	body        []byte
	contentType string
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("http call '%s' status: %d", h.url, h.code)
}

func (h HTTPError) HttpStatus() int {
	return h.code
}

func (h HTTPError) Body() []byte {
	return h.body
}

func (h HTTPError) ContentType() string {
	return h.contentType
}
