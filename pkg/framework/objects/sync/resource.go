package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-hightouch/pkg/hightouch"
)

// SyncResource is the resource implementation.
type SyncResource struct {
	client *hightouch.Client
}

// NewSyncResource is a helper function to simplify resource server allocation.
func NewSyncResource() resource.Resource {
	return &SyncResource{}
}

// Metadata returns the resource type name.
func (r *SyncResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sync"
}

// Schema defines the schema for the resource.
func (r *SyncResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = SyncResourceSchema
}

// Configure adds the hightouch configured client to the resource.
func (r *SyncResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
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
	r.client = client
}

// Create creates the resource and sets the initial state.
func (r *SyncResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan SyncResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse configuration and schedule JSON strings
	var configuration map[string]interface{}
	if err := json.Unmarshal([]byte(plan.Configuration.ValueString()), &configuration); err != nil {
		resp.Diagnostics.AddError("Invalid Configuration JSON", "Could not parse configuration JSON: "+err.Error())
		return
	}

	var schedule map[string]interface{}
	if err := json.Unmarshal([]byte(plan.Schedule.ValueString()), &schedule); err != nil {
		resp.Diagnostics.AddError("Invalid Schedule JSON", "Could not parse schedule JSON: "+err.Error())
		return
	}

	// Call the API to create the sync
	sync, err := r.client.CreateHightouchSync(
		plan.Name.ValueString(),
		plan.Slug.ValueString(),
		int(plan.SourceID.ValueInt64()),
		int(plan.DestinationID.ValueInt64()),
		int(plan.ModelID.ValueInt64()),
		configuration,
		schedule,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating sync", "Could not create sync, unexpected error: "+err.Error())
		return
	}

	// Map response body to the plan
	syncID := *sync.ID
	plan.ID = types.Int64Value(int64(syncID))
	plan.WorkspaceID = types.Int64Value(int64(sync.WorkspaceID))
	plan.Status = types.StringValue(sync.Status)
	plan.Disabled = types.BoolValue(sync.Disabled)
	plan.CreatedAt = types.StringValue(sync.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(sync.UpdatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the resource state with the latest data.
func (r *SyncResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state SyncResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed sync from Hightouch API
	syncID := int(state.ID.ValueInt64())
	if syncID == 0 {
		resp.Diagnostics.AddError("Invalid Sync ID", "The sync ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}
	sync, err := r.client.GetHightouchSync(syncID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading sync", "Could not read sync, unexpected error: "+err.Error())
		return
	}

	workspaceID := int64(sync.WorkspaceID)
	if workspaceID == 0 {
		resp.Diagnostics.AddError("Invalid Workspace ID", "The workspace ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
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

	// Overwrite state with refreshed values
	state.ID = types.Int64Value(int64(syncID))
	state.Name = types.StringValue(sync.Name)
	state.Slug = types.StringValue(sync.Slug)
	state.DestinationID = types.Int64Value(int64(sync.DestinationID))
	state.ModelID = types.Int64Value(int64(sync.ModelID))
	state.Configuration = types.StringValue(string(configJSON))
	state.Schedule = types.StringValue(string(scheduleJSON))
	state.Status = types.StringValue(sync.Status)
	state.Disabled = types.BoolValue(sync.Disabled)
	state.UpdatedAt = types.StringValue(sync.UpdatedAt.String())
	state.WorkspaceID = types.Int64Value(int64(sync.WorkspaceID))
	state.CreatedAt = types.StringValue(sync.CreatedAt.String())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated state.
func (r *SyncResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state SyncResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	syncID := int(state.ID.ValueInt64())
	if syncID == 0 {
		resp.Diagnostics.AddError("Invalid Sync ID", "The sync ID must be set before updating. Please ensure the resource has been created successfully before attempting to update it.")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse configuration and schedule JSON strings
	var configuration map[string]interface{}
	if err := json.Unmarshal([]byte(plan.Configuration.ValueString()), &configuration); err != nil {
		resp.Diagnostics.AddError("Invalid Configuration JSON", "Could not parse configuration JSON: "+err.Error())
		return
	}

	var schedule map[string]interface{}
	if err := json.Unmarshal([]byte(plan.Schedule.ValueString()), &schedule); err != nil {
		resp.Diagnostics.AddError("Invalid Schedule JSON", "Could not parse schedule JSON: "+err.Error())
		return
	}

	// Call the API to update the sync
	sync, err := r.client.UpdateHightouchSync(
		syncID,
		plan.Name.ValueString(),
		configuration,
		schedule,
		plan.Disabled.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating sync", "Could not update sync, unexpected error: "+err.Error())
		return
	}

	// Update the plan with the response from the API
	plan.UpdatedAt = types.StringValue(sync.UpdatedAt.String())
	plan.WorkspaceID = types.Int64Value(int64(sync.WorkspaceID))
	plan.Status = types.StringValue(sync.Status)
	plan.ID = types.Int64Value(int64(syncID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource from the remote API.
func (r *SyncResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state SyncResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For now, we'll simulate a successful deletion since delete is not implemented in the SDK
	resp.Diagnostics.AddWarning("Delete not implemented", "The Hightouch API does not currently support deleting syncs via the API. The resource will be removed from Terraform state, but not from Hightouch.")
}

// ImportState imports the resource into Terraform state.
func (r *SyncResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID for Import", "ID must be a valid integer.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
