package model

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hightouch/pkg/hightouch"
)

// ModelDataSource is the data source implementation.
type ModelDataSource struct {
	client *hightouch.Client
}

// NewModelDataSource is a helper function to simplify data source server allocation.
func NewModelDataSource() datasource.DataSource {
	return &ModelDataSource{}
}

// Metadata returns the data source type name.
func (d *ModelDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

// Schema defines the schema for the data source.
func (d *ModelDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ModelDataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *ModelDataSource) Configure(
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
func (d *ModelDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ModelResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get model from Hightouch API
	modelID := int(config.ID.ValueInt64())
	if modelID == 0 {
		resp.Diagnostics.AddError("Invalid Model ID", "The model ID must be provided.")
		return
	}

	model, err := d.client.GetHightouchModel(modelID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading model", "Could not read model, unexpected error: "+err.Error())
		return
	}

	// Map API response to Terraform state
	config.ID = types.Int64Value(int64(*model.ID))
	config.Name = types.StringValue(model.Name)
	config.Slug = types.StringValue(model.Slug)
	config.SourceID = types.Int64Value(int64(model.SourceID))
	config.SQL = types.StringValue(model.SQL)
	config.DBTable = types.StringValue(model.DBTable)
	config.QueryType = types.StringValue(model.QueryType)
	config.PrimaryKey = types.StringValue(model.PrimaryKey)
	config.IsSchema = types.BoolValue(model.IsSchema)
	config.WorkspaceID = types.Int64Value(int64(model.WorkspaceID))
	config.CreatedAt = types.StringValue(model.CreatedAt.String())
	config.UpdatedAt = types.StringValue(model.UpdatedAt.String())

	// Set state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
