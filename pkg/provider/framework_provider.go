package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"terraform-provider-hightouch/pkg/framework/objects/snowflake_source"
	"terraform-provider-hightouch/pkg/hightouch"
)

type hightouchProvider struct {
	// version can be set during provider build time
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hightouchProvider{
			version: version,
		}
	}
}

type hightouchProviderModel struct {
	APIKey     types.String `tfsdk:"api_key"`
	APIBaseURL types.String `tfsdk:"api_base_url"`
}

func (p *hightouchProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hightouch"
	resp.Version = p.version
}

func (p *hightouchProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for interacting with the Hightouch API.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The API key for Hightouch.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_base_url": schema.StringAttribute{
				Description: "The base URL for the Hightouch API. Defaults to https://api.hightouch.com/api/v1",
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (p *hightouchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config hightouchProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := ""
	// If the API key is unknown, we can't configure the client.
	// This can happen during planning if the key is from a dynamic source.
	if config.APIKey.ValueString() != "" {
		apiKey = config.APIKey.ValueString()
	} else if os.Getenv("HIGHTOUCH_API_KEY") != "" {
		apiKey = os.Getenv("HIGHTOUCH_API_KEY")
	} else {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The Hightouch API key must be provided either in the configuration or as an environment variable.",
		)
		return
	}

	apiBaseUrl := ""
	// Get API base URL from config or use default
	if config.APIBaseURL.ValueString() != "" {
		apiBaseUrl = config.APIBaseURL.ValueString()
	} else if os.Getenv("HIGHTOUCH_API_BASE_URL") != "" {
		apiBaseUrl = os.Getenv("HIGHTOUCH_API_BASE_URL")
	} else {
		apiBaseUrl = "https://api.hightouch.com/api/v1"
	}

	// Create a new client and make it available to all resources
	client := hightouch.NewClient(apiKey, apiBaseUrl) // Assuming NewClient is in this package or imported

	// This is the crucial step: we pass the configured client
	// to all resources and data sources.
	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *hightouchProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		snowflake_source.NewHightouchSnowflakeSourceResource,
	}
}

// DataSources You would also have a DataSources method, which might be empty for now
func (p *hightouchProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
