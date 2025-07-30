package model

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-hightouch/pkg/hightouch"
)

// ModelResource is the resource implementation.
type ModelResource struct {
	client *hightouch.Client
}

// NewModelResource is a helper function to simplify resource server allocation.
func NewModelResource() resource.Resource {
	return &ModelResource{}
}

// Metadata returns the resource type name.
func (r *ModelResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

// Schema defines the schema for the resource.
func (r *ModelResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ModelResourceSchema
}

// Configure adds the hightouch configured client to the resource.
func (r *ModelResource) Configure(
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
func (r *ModelResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan ModelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to create the model
	model, err := r.client.CreateHightouchModel(
		plan.Name.ValueString(),
		plan.Slug.ValueString(),
		int(plan.SourceID.ValueInt64()),
		plan.SQL.ValueString(),
		plan.QueryType.ValueString(),
		plan.PrimaryKey.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating model", "Could not create model, unexpected error: "+err.Error())
		return
	}

	// Map response body to the plan
	modelID := *model.ID
	plan.ID = types.Int64Value(int64(modelID))
	plan.WorkspaceID = types.Int64Value(int64(model.WorkspaceID))
	plan.CreatedAt = types.StringValue(model.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(model.UpdatedAt.String())
	plan.DBTable = types.StringValue(model.DBTable)
	plan.IsSchema = types.BoolValue(model.IsSchema)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the resource state with the latest data.
func (r *ModelResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ModelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed model from Hightouch API
	modelID := int(state.ID.ValueInt64())
	if modelID == 0 {
		resp.Diagnostics.AddError("Invalid Model ID", "The model ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}
	model, err := r.client.GetHightouchModel(modelID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading model", "Could not read model, unexpected error: "+err.Error())
		return
	}

	workspaceID := int64(model.WorkspaceID)
	if workspaceID == 0 {
		resp.Diagnostics.AddError("Invalid Workspace ID", "The workspace ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}

	// Overwrite state with refreshed values
	state.ID = types.Int64Value(int64(modelID))
	state.Name = types.StringValue(model.Name)
	state.Slug = types.StringValue(model.Slug)
	state.SourceID = types.Int64Value(int64(model.SourceID))
	state.SQL = types.StringValue(model.SQL)
	state.DBTable = types.StringValue(model.DBTable)
	state.QueryType = types.StringValue(model.QueryType)
	state.PrimaryKey = types.StringValue(model.PrimaryKey)
	state.IsSchema = types.BoolValue(model.IsSchema)
	state.UpdatedAt = types.StringValue(model.UpdatedAt.String())
	state.WorkspaceID = types.Int64Value(int64(model.WorkspaceID))
	state.CreatedAt = types.StringValue(model.CreatedAt.String())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated state.
func (r *ModelResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ModelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	modelID := int(state.ID.ValueInt64())
	if modelID == 0 {
		resp.Diagnostics.AddError("Invalid Model ID", "The model ID must be set before updating. Please ensure the resource has been created successfully before attempting to update it.")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to update the model
	model, err := r.client.UpdateHightouchModel(
		modelID,
		plan.Name.ValueString(),
		plan.SQL.ValueString(),
		plan.PrimaryKey.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating model", "Could not update model, unexpected error: "+err.Error())
		return
	}

	// Update the plan with the response from the API
	plan.UpdatedAt = types.StringValue(model.UpdatedAt.String())
	plan.WorkspaceID = types.Int64Value(int64(model.WorkspaceID))
	plan.ID = types.Int64Value(int64(modelID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource from the remote API.
func (r *ModelResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state ModelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For now, we'll simulate a successful deletion since delete is not implemented in the SDK
	resp.Diagnostics.AddWarning("Delete not implemented", "The Hightouch API does not currently support deleting models via the API. The resource will be removed from Terraform state, but not from Hightouch.")
}

// ImportState imports the resource into Terraform state.
func (r *ModelResource) ImportState(
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
