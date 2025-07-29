package snowflake_source

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hightouch/pkg/hightouch"
)

// SnowflakeSourceDataSource is the data source implementation.
type SnowflakeSourceDataSource struct {
	client *hightouch.Client
}

// NewSnowflakeSourceDataSource is a helper function to simplify data source server allocation.
func NewSnowflakeSourceDataSource() datasource.DataSource {
	return &SnowflakeSourceDataSource{}
}

// Metadata returns the data source type name.
func (d *SnowflakeSourceDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_snowflake_source"
}

// Schema defines the schema for the data source.
func (d *SnowflakeSourceDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SnowflakeSourceDataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *SnowflakeSourceDataSource) Configure(
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
func (d *SnowflakeSourceDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SnowflakeSourceResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get source from Hightouch API
	sourceID := int(config.ID.ValueInt64())
	if sourceID == 0 {
		resp.Diagnostics.AddError("Invalid Source ID", "The source ID must be provided.")
		return
	}

	source, err := d.client.GetSnowflakeSource(sourceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading source", "Could not read source, unexpected error: "+err.Error())
		return
	}

	// Map API response to Terraform state
	config.ID = types.Int64Value(int64(*source.ID))
	config.Name = types.StringValue(source.Name)
	config.Slug = types.StringValue(source.Slug)
	config.Type = types.StringValue(source.Type)
	config.WorkspaceID = types.Int64Value(int64(source.WorkspaceID))
	config.CreatedAt = types.StringValue(source.CreatedAt.String())
	config.UpdatedAt = types.StringValue(source.UpdatedAt.String())

	// Extract configuration fields
	if accountString, ok := source.Configuration["account"].(string); ok {
		config.Account = types.StringValue(accountString)
	}
	if portFloat, ok := source.Configuration["port"].(float64); ok {
		config.Port = types.Int64Value(int64(portFloat))
	}
	if usernameString, ok := source.Configuration["username"].(string); ok {
		config.Username = types.StringValue(usernameString)
	}
	if passwordString, ok := source.Configuration["password"].(string); ok {
		config.Password = types.StringValue(passwordString)
	}
	if databaseString, ok := source.Configuration["database"].(string); ok {
		config.Database = types.StringValue(databaseString)
	}
	if warehouseString, ok := source.Configuration["warehouse"].(string); ok {
		config.Warehouse = types.StringValue(warehouseString)
	}

	// Set state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
