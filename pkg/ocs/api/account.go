package api

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type AccountState struct {
	State       string    `json:"state"`
	StateReason string    `json:"stateReason"`
	ValidFrom   time.Time `json:"stateValidReason"`
}

type Account struct {
	RequestID      string            `json:"requestId"`
	AccountProfile string            `json:"accountProfile,omitempty"`
	ExternalID     string            `json:"accountExternalId,omitempty"`
	State          AccountState      `json:"state"`
	Name           string            `json:"customName"`
	Attributes     map[string]string `json:"customAttributes"`
}

func (c *Client) CreateAccount(ctx context.Context, a Account) (*Response, error) {
	data, _ := json.Marshal(a)
	buff := bytes.NewBuffer(data)
	r := c.prepareRequest(ctx, "POST", c.buildURL(apiCreateCustomer), buff)

	resp, err := c.call(r)
	if err != nil {
		return nil, errors.Wrapf(err, "create account")
	}
	apiResp, err := c.parseAPIResponse(resp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}
