package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Policy struct {
	OutputCurrency string `json:"outputCurrency"`
	Owner          string `json:"owner"`
	Role           string `json:"role"`
	Name           string `json:"name"`
	ID             string `json:"id"`
	Type           string `json:"type"`
}

func (c *Client) PolicyList(ctx context.Context) ([]Policy, error) {
	reqStructure := struct {
		requestJobDescription
		InputSettings struct {
			Type string `json:"type"`
		} `json:"inputSettings"`
	}{
		requestJobDescription: requestJobDescription{
			Type: "get",
			Credentials: credentials{
				PartnerUserID:     c.userID,
				PartnerUserSecret: c.userSecret,
			},
		},
		InputSettings: struct {
			Type string "json:\"type\""
		}{
			Type: "policyList",
		},
	}
	reqJSON, err := json.Marshal(reqStructure)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request body: %w", err)
	}

	data := bytes.NewBuffer([]byte("requestJobDescription="))
	data.Write(reqJSON)

	url := c.serverURL.ResolveReference(relativeURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), data)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %w", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	respStructure := &struct {
		PolicyList []Policy `json:"policyList"`
	}{}

	err = json.Unmarshal(bodyBytes, &respStructure)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return respStructure.PolicyList, nil
}
