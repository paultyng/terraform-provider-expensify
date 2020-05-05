package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OnReceive struct {
	ImmediateResponse []string `json:"immediateResponse"`
}

type InputSettingsFilters struct {
	ReportIDList string `json:"reportIDList,omitempty"`
}

// TODO: move this just in to the func
type InputSettings struct {
	Type    string               `json:"type"`
	Filters InputSettingsFilters `json:"filters"`
}

// TODO: move this just in to the func
type OutputSettings struct {
	FileExtension string `json:"fileExtension"`
}

type FileRequest struct {
	OnReceive      OnReceive      `json:"onReceive"`
	InputSettings  InputSettings  `json:"inputSettings"`
	OutputSettings OutputSettings `json:"outputSettings"`
}

func (c *Client) File(ctx context.Context, fileReq FileRequest, template string) (string, error) {

	reqStructure := struct {
		requestJobDescription
		FileRequest
		Test bool `json:"test"`
	}{
		requestJobDescription: requestJobDescription{
			Type: "file",
			Credentials: credentials{
				PartnerUserID:     c.userID,
				PartnerUserSecret: c.userSecret,
			},
		},
		FileRequest: fileReq,
		Test:        true,
	}
	reqJSON, err := json.Marshal(reqStructure)
	if err != nil {
		return "", fmt.Errorf("unable to marshal request body: %w", err)
	}

	data := bytes.NewBuffer([]byte("requestJobDescription="))
	data.Write(reqJSON)
	data.WriteString("&template=")
	data.WriteString(url.QueryEscape(template))

	url := c.serverURL.ResolveReference(relativeURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), data)
	if err != nil {
		return "", fmt.Errorf("unable to build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to perform request: %w", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read body: %w", err)
	}

	return string(bodyBytes), nil
}
