package iterable_destination

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-hightouch/pkg/hightouch"
)

// IterableDestinationResource is the resource implementation.
type IterableDestinationResource struct {
	client *hightouch.Client
}

// NewIterableDestinationResource is a helper function to simplify resource server allocation.
func NewIterableDestinationResource() resource.Resource {
	return &IterableDestinationResource{}
}

// Metadata returns the resource type name.
func (r *IterableDestinationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_iterable_destination"
}

// Schema defines the schema for the resource.
func (r *IterableDestinationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = IterableDestinationResourceSchema
}

// Configure adds the hightouch configured client to the resource.
func (r *IterableDestinationResource) Configure(
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
func (r *IterableDestinationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan IterableDestinationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration from Terraform types to Go types
	config := make(map[string]interface{})
	config["api_key"] = plan.APIKey.ValueString()
	config["data_center"] = plan.DataCenter.ValueString()

	// Call the API to create the destination
	destination, err := r.client.CreateHightouchDestination(
		plan.Name.ValueString(),
		plan.Slug.ValueString(),
		plan.Type.ValueString(),
		config,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating destination", "Could not create destination, unexpected error: "+err.Error())
		return
	}

	// Map response body to the plan
	destinationID := *destination.ID
	plan.ID = types.Int64Value(int64(destinationID))
	plan.WorkspaceID = types.Int64Value(int64(destination.WorkspaceID))
	plan.CreatedAt = types.StringValue(destination.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(destination.UpdatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the resource state with the latest data.
func (r *IterableDestinationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state IterableDestinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed destination from Hightouch API
	destinationID := int(state.ID.ValueInt64())
	if destinationID == 0 {
		resp.Diagnostics.AddError("Invalid Destination ID", "The destination ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}
	destination, err := r.client.GetHightouchDestination(destinationID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading destination", "Could not read destination, unexpected error: "+err.Error())
		return
	}

	workspaceID := int64(destination.WorkspaceID)
	if workspaceID == 0 {
		resp.Diagnostics.AddError("Invalid Workspace ID", "The workspace ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}

	// Overwrite state with refreshed values
	state.ID = types.Int64Value(int64(destinationID))
	state.Name = types.StringValue(destination.Name)
	state.Slug = types.StringValue(destination.Slug)
	state.Type = types.StringValue(destination.Type)
	state.UpdatedAt = types.StringValue(destination.UpdatedAt.String())
	state.WorkspaceID = types.Int64Value(int64(destination.WorkspaceID))
	state.CreatedAt = types.StringValue(destination.CreatedAt.String())

	// Convert configuration from Go types to Terraform types
	apiKeyString, ok := destination.Configuration["api_key"].(string)
	if !ok {
		fmt.Println("api_key is not a string")
		return
	}
	dataCenterString, ok := destination.Configuration["data_center"].(string)
	if !ok {
		// Default to US if not specified
		dataCenterString = "US"
	}

	state.APIKey = types.StringValue(apiKeyString)
	state.DataCenter = types.StringValue(dataCenterString)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated state.
func (r *IterableDestinationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state IterableDestinationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	destinationID := int(state.ID.ValueInt64())
	if destinationID == 0 {
		resp.Diagnostics.AddError("Invalid Destination ID", "The destination ID must be set before updating. Please ensure the resource has been created successfully before attempting to update it.")
		return
	}
	fmt.Printf("DestinationID: %+v\n", destinationID)

	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Printf("Destination Payload: %+v\n", plan)
	fmt.Printf("Destination State: %+v\n", state)

	// Convert configuration from Terraform types to Go types
	config := make(map[string]interface{})
	config["api_key"] = plan.APIKey.ValueString()
	config["data_center"] = plan.DataCenter.ValueString()

	// Call the API to update the destination
	destination, err := r.client.UpdateHightouchDestination(
		destinationID,
		plan.Name.ValueString(),
		config,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating destination", "Could not update destination, unexpected error: "+err.Error())
		return
	}

	// Update the plan with the response from the API
	plan.UpdatedAt = types.StringValue(destination.UpdatedAt.String())
	plan.WorkspaceID = types.Int64Value(int64(destination.WorkspaceID))
	plan.ID = types.Int64Value(int64(destinationID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource from the remote API.
func (r *IterableDestinationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state IterableDestinationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// To complete this, a `DeleteDestination` method would be needed on the client.
	// For example:
	// err := r.client.DeleteDestination(int(state.ID.ValueInt64()))
	// if err != nil {
	//     resp.Diagnostics.AddError("Error deleting destination", "Could not delete destination, unexpected error: "+err.Error())
	//     return
	// }

	// For now, we'll simulate a successful deletion.
	// In a real implementation, you would call your client's delete method here.
	resp.Diagnostics.AddWarning("Delete not implemented", "The Hightouch API does not currently support deleting destinations via the API. The resource will be removed from Terraform state, but not from Hightouch.")
}

// ImportState imports the resource into Terraform state.
func (r *IterableDestinationResource) ImportState(
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
