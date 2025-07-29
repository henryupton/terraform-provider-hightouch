package iterable_destination

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hightouch/pkg/hightouch"
)

// IterableDestinationDataSource is the data source implementation.
type IterableDestinationDataSource struct {
	client *hightouch.Client
}

// NewIterableDestinationDataSource is a helper function to simplify data source server allocation.
func NewIterableDestinationDataSource() datasource.DataSource {
	return &IterableDestinationDataSource{}
}

// Metadata returns the data source type name.
func (d *IterableDestinationDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_iterable_destination"
}

// Schema defines the schema for the data source.
func (d *IterableDestinationDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = IterableDestinationDataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *IterableDestinationDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hightouch.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *IterableDestinationDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config IterableDestinationResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get destination from Hightouch API
	destinationID := int(config.ID.ValueInt64())
	if destinationID == 0 {
		resp.Diagnostics.AddError("Invalid Destination ID", "The destination ID must be provided.")
		return
	}

	destination, err := d.client.GetHightouchDestination(destinationID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading destination", "Could not read destination, unexpected error: "+err.Error())
		return
	}

	// Map API response to Terraform state
	config.ID = types.Int64Value(int64(*destination.ID))
	config.Name = types.StringValue(destination.Name)
	config.Slug = types.StringValue(destination.Slug)
	config.Type = types.StringValue(destination.Type)
	config.WorkspaceID = types.Int64Value(int64(destination.WorkspaceID))
	config.CreatedAt = types.StringValue(destination.CreatedAt.String())
	config.UpdatedAt = types.StringValue(destination.UpdatedAt.String())

	// Extract configuration fields
	if apiKeyString, ok := destination.Configuration["api_key"].(string); ok {
		config.APIKey = types.StringValue(apiKeyString)
	}
	if dataCenterString, ok := destination.Configuration["data_center"].(string); ok {
		config.DataCenter = types.StringValue(dataCenterString)
	} else {
		// Default to US if not specified
		config.DataCenter = types.StringValue("US")
	}

	// Set state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
