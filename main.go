package main // import "github.com/paultyng/terraform-provider-expensify"

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/paultyng/terraform-provider-expensify/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
