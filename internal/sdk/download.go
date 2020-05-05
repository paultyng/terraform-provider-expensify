package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) Download(ctx context.Context, name, system string) (string, error) {
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
