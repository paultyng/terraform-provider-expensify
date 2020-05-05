package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/terraform-provider-expensify/internal/sdk"
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
		DataSourcesMap: map[string]*schema.Resource{
			"expensify_policy": dataPolicy(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"expensify_report": resourceReport(),
		},
	}
	p.ConfigureContextFunc = configure(p)
	return p
}

func configure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userID := d.Get("partner_user_id").(string)
		userSecret := d.Get("partner_user_secret").(string)
		serverURL := d.Get("server_url").(string)

		c, err := sdk.New(userID, userSecret, serverURL, &http.Client{
			Transport: logging.NewTransport("expensify", http.DefaultTransport),
		})
		if err != nil {
			return nil, fromErr(err)
		}

		return c, nil
	}
}
