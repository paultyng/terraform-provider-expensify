package provider

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"partner_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EXPENSIFY_USER_ID", ""),
			},
			"partner_user_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EXPENSIFY_SECRET", ""),
			},
			"server_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EXPENSIFY_SERVER", "https://integrations.expensify.com"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"expensify_report": resourceReport(),
		},
	}
	p.ConfigureFunc = configure(p)
	return p
}

func configure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		userID := d.Get("partner_user_id").(string)
		userSecret := d.Get("partner_user_secret").(string)
		serverURL, err := url.Parse(d.Get("server_url").(string))
		if err != nil {
			return nil, fmt.Errorf("unable to parse URL: %w", err)
		}

		c := &client{
			userID:     userID,
			userSecret: userSecret,
			serverURL:  serverURL,

			c: &http.Client{
				Transport: logging.NewTransport("expensify", http.DefaultTransport),
			},
		}

		return c, nil
	}
}
