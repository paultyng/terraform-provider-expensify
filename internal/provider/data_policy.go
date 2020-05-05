package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/terraform-provider-expensify/internal/sdk"
)

func dataPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sdk.Client)
	name := d.Get("name").(string)

	policies, err := c.PolicyList(ctx)
	if err != nil {
		return fromErr(err)
	}

	var policy *sdk.Policy
	for _, p := range policies {
		if p.Name != name {
			continue
		}
		policy = &p
		break
	}
	if policy == nil {
		return errorf("not found")
	}

	d.SetId(policy.ID)
	return nil
}
