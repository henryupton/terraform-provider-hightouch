package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hightouch/pkg/hightouch"
)

// SyncDataSource is the data source implementation.
type SyncDataSource struct {
	client *hightouch.Client
}

// NewSyncDataSource is a helper function to simplify data source server allocation.
func NewSyncDataSource() datasource.DataSource {
	return &SyncDataSource{}
}

// Metadata returns the data source type name.
func (d *SyncDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sync"
}

// Schema defines the schema for the data source.
func (d *SyncDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SyncDataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *SyncDataSource) Configure(
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
func (d *SyncDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SyncResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get sync from Hightouch API
	syncID := int(config.ID.ValueInt64())
	if syncID == 0 {
		resp.Diagnostics.AddError("Invalid Sync ID", "The sync ID must be provided.")
		return
	}

	sync, err := d.client.GetHightouchSync(syncID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading sync", "Could not read sync, unexpected error: "+err.Error())
		return
	}

	// Convert configuration and schedule maps to JSON strings
	configJSON, err := json.Marshal(sync.Configuration)
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling configuration", "Could not marshal configuration to JSON: "+err.Error())
		return
	}

	scheduleJSON, err := json.Marshal(sync.Schedule)
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling schedule", "Could not marshal schedule to JSON: "+err.Error())
		return
	}

	// Map API response to Terraform state
	config.ID = types.Int64Value(int64(*sync.ID))
	config.Name = types.StringValue(sync.Name)
	config.Slug = types.StringValue(sync.Slug)
	config.DestinationID = types.Int64Value(int64(sync.DestinationID))
	config.ModelID = types.Int64Value(int64(sync.ModelID))
	config.Configuration = types.StringValue(string(configJSON))
	config.Schedule = types.StringValue(string(scheduleJSON))
	config.Status = types.StringValue(sync.Status)
	config.Disabled = types.BoolValue(sync.Disabled)
	config.WorkspaceID = types.Int64Value(int64(sync.WorkspaceID))
	config.CreatedAt = types.StringValue(sync.CreatedAt.String())
	config.UpdatedAt = types.StringValue(sync.UpdatedAt.String())

	// Set state
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
