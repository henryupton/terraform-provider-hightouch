package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-hightouch/provider"
	"terraform-provider-hightouch/hightouch"
)

// main is the entry point for the provider.
func main() {


	plugin.Serve(&plugin.ServeOpts{
		// ProviderFunc returns the provider's schema and is the only required
		// field. By putting all our code in the `main` package, we can
		// reference the `Provider` function directly.
		ProviderFunc: Provider,
	})
}

// Provider defines the provider's schema, resources, and configuration.
func Provider() *schema.Provider {
	return &schema.Provider{
		// Schema defines provider-level configuration options, like API keys.
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HIGHTOUCH_API_KEY", nil),
				Description: "The API key for the Hightouch REST API.",
			},
			"api_base_url": {
			    Type:        schema.TypeString,
                Optional:    true,
                Sensitive:   false,
                DefaultFunc: schema.EnvDefaultFunc("HIGHTOUCH_API_BASE_URL", "https://api.hightouch.com"),
                Description: "The base URL for the Hightouch API.",
            },
		},

		// ResourcesMap maps the resource names in Terraform configurations
		// to their corresponding schema and CRUD functions.
		ResourcesMap: map[string]*schema.Resource{
			"hightouch_source": provider.ResourceHightouchSource(),
		},

		// DataSourcesMap is for read-only data sources.
		DataSourcesMap: map[string]*schema.Resource{},

		// ConfigureContextFunc is used to configure the provider, for example,
		// by setting up an API client.
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure processes the provider configuration and returns a configured API client.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	if apiKey == "" {
		return nil, diag.Errorf("API key is required")
	}

	client := hightouch.NewClient(apiKey)
	return client, nil
}
