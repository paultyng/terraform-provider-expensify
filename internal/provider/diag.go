package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func fromErr(err error) diag.Diagnostics {
	return diag.Diagnostics{diag.FromErr(err)}
}

func errorf(format string, arg ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		{
			Severity: diag.Error,
			Summary:  fmt.Sprintf(format, arg...),
		},
	}
}
