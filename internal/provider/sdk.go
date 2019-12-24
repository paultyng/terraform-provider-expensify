package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var relativeURL *url.URL

func init() {
	var err error
	relativeURL, err = url.Parse("/Integration-Server/ExpensifyIntegrations")
	if err != nil {
		panic(err)
	}
}

type client struct {
	userID     string
	userSecret string
	serverURL  *url.URL

	c *http.Client
}

type credentials struct {
	PartnerUserID     string `json:"partnerUserID"`
	PartnerUserSecret string `json:"partnerUserSecret"`
}

type requestJobDescription struct {
	Type        string      `json:"type"`
	Credentials credentials `json:"credentials"`
}

type OnReceive struct {
	ImmediateResponse []string `json:"immediateResponse"`
}

type InputSettingsFilters struct {
	ReportIDList string `json:"reportIDList,omitempty"`
}

type InputSettings struct {
	Type    string               `json:"type"`
	Filters InputSettingsFilters `json:"filters"`
}

type OutputSettings struct {
	FileExtension string `json:"fileExtension"`
}

type FileRequest struct {
	OnReceive      OnReceive      `json:"onReceive"`
	InputSettings  InputSettings  `json:"inputSettings"`
	OutputSettings OutputSettings `json:"outputSettings"`
}

func (c *client) File(ctx context.Context, fileReq FileRequest, template string) (string, error) {
	reqStructure := struct {
		requestJobDescription
		FileRequest
	}{
		requestJobDescription: requestJobDescription{
			Type: "file",
			Credentials: credentials{
				PartnerUserID:     c.userID,
				PartnerUserSecret: c.userSecret,
			},
		},
		FileRequest: fileReq,
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

func (c *client) Download(ctx context.Context, name, system string) (string, error) {
	reqStructure := struct {
		requestJobDescription
		FileName string `json:"fileName"`

		//optional
		FileSystem string `json:"fileSystem,omitempty"`
	}{
		requestJobDescription: requestJobDescription{
			Type: "download",
			Credentials: credentials{
				PartnerUserID:     c.userID,
				PartnerUserSecret: c.userSecret,
			},
		},
		FileName:   name,
		FileSystem: system,
	}
	reqJSON, err := json.Marshal(reqStructure)
	if err != nil {
		return "", fmt.Errorf("unable to marshal request body: %w", err)
	}

	data := bytes.NewBuffer([]byte("requestJobDescription="))
	data.Write(reqJSON)

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
