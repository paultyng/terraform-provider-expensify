package sdk

import (
	"fmt"
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

func New(userID, userSecret, serverURL string, c *http.Client) (*Client, error) {
	surl, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	if !surl.IsAbs() {
		return nil, fmt.Errorf("server URL must be absolute, not %q", serverURL)
	}

	return &Client{
		userID:     userID,
		userSecret: userSecret,
		serverURL:  surl,
		c:          c,
	}, nil
}

type Client struct {
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
