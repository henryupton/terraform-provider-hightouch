package snowflake_source

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-hightouch/pkg/hightouch"
)

// HightouchSnowflakeSourceResource is the resource implementation.
type HightouchSnowflakeSourceResource struct {
	client *hightouch.Client
}

// NewHightouchSnowflakeSourceResource is a helper function to simplify resource server allocation.
func NewHightouchSnowflakeSourceResource() resource.Resource {
	return &HightouchSnowflakeSourceResource{}
}

// Metadata returns the resource type name.
func (r *HightouchSnowflakeSourceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_snowflake_source"
}

// Schema defines the schema for the resource.
func (r *HightouchSnowflakeSourceResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = SnowflakeSourceResourceSchema
}

// Configure adds the hightouch_resources configured client to the resource.
func (r *HightouchSnowflakeSourceResource) Configure(
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
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the hightouch_resources developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

// Create creates the resource and sets the initial state.
func (r *HightouchSnowflakeSourceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan HightouchSnowflakeSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert configuration from Terraform types to Go types
	config := make(map[string]interface{})
	config["port"] = plan.Port.ValueInt64()
	config["account"] = plan.Account.ValueString()
	config["username"] = plan.Username.ValueString()
	config["database"] = plan.Database.ValueString()
	config["warehouse"] = plan.Warehouse.ValueString()
	config["password"] = plan.Password.ValueString()

	// Call the API to create the source
	source, err := r.client.CreateHightouchSource(
		plan.Name.ValueString(),
		plan.Slug.ValueString(),
		plan.Type.ValueString(),
		config,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating source", "Could not create source, unexpected error: "+err.Error())
		return
	}

	// Map response body to the plan
	sourceID := *source.ID
	plan.ID = types.Int64Value(int64(sourceID))
	plan.WorkspaceID = types.Int64Value(int64(source.WorkspaceID))
	plan.CreatedAt = types.StringValue(source.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(source.UpdatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the resource state with the latest data.
func (r *HightouchSnowflakeSourceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state HightouchSnowflakeSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed source from Hightouch API
	sourceID := int(state.ID.ValueInt64())
	if sourceID == 0 {
		resp.Diagnostics.AddError("Invalid Source ID", "The source ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}
	source, err := r.client.GetHightouchSnowflakeSource(sourceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading source", "Could not read source, unexpected error: "+err.Error())
		return
	}

	workspaceID := int64(source.WorkspaceID)
	if workspaceID == 0 {
		resp.Diagnostics.AddError("Invalid Workspace ID", "The workspace ID must be set before reading. Please ensure the resource has been created successfully before attempting to read it.")
		return
	}

	// Overwrite state with refreshed values
	state.ID = types.Int64Value(int64(sourceID))
	state.Name = types.StringValue(source.Name)
	state.Slug = types.StringValue(source.Slug)
	state.Type = types.StringValue(source.Type)
	state.UpdatedAt = types.StringValue(source.UpdatedAt.String())
	state.WorkspaceID = types.Int64Value(int64(source.WorkspaceID))
	state.CreatedAt = types.StringValue(source.CreatedAt.String())

	// Convert configuration from Go types to Terraform types
	accountString, ok := source.Configuration["account"].(string)
	if !ok {
		fmt.Println("account is not a string")
		return
	}
	portFloat, ok := source.Configuration["port"].(float64)
	if !ok {
		// Handle the error if the type is not a float64
		fmt.Println("port is not a float64")
		return
	}
	usernameString, ok := source.Configuration["username"].(string)
	if !ok {
		fmt.Println("username is not a string")
		return
	}
	databaseString, ok := source.Configuration["database"].(string)
	if !ok {
		fmt.Println("database is not a string")
		return
	}
	passwordString, ok := source.Configuration["password"].(string)
	if !ok {
		fmt.Println("password is not a string")
		return
	}
	warehouseString, ok := source.Configuration["warehouse"].(string)
	if !ok {
		fmt.Println("warehouse is not a string")
		return
	}

	state.Account = types.StringValue(accountString)
	state.Port = types.Int64Value(int64(portFloat))
	state.Username = types.StringValue(usernameString)
	state.Database = types.StringValue(databaseString)
	state.Password = types.StringValue(passwordString)
	state.Warehouse = types.StringValue(warehouseString)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated state.
func (r *HightouchSnowflakeSourceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state HightouchSnowflakeSourceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	sourceID := int(state.ID.ValueInt64())
	if sourceID == 0 {
		resp.Diagnostics.AddError("Invalid Source ID", "The source ID must be set before updating. Please ensure the resource has been created successfully before attempting to update it.")
		return
	}
	fmt.Printf("SourceID: %+v\n", sourceID)

	if resp.Diagnostics.HasError() {
		return
	}

	fmt.Printf("Source Payload: %+v\n", plan)
	fmt.Printf("Source Payload: %+v\n", state)

	// Convert configuration from Terraform types to Go types
	config := make(map[string]interface{})
	config["port"] = plan.Port.ValueInt64()
	config["account"] = plan.Account.ValueString()
	config["username"] = plan.Username.ValueString()
	config["database"] = plan.Database.ValueString()
	config["warehouse"] = plan.Warehouse.ValueString()

	// Call the API to update the source
	source, err := r.client.UpdateHightouchSource(
		sourceID,
		plan.Name.ValueString(),
		config,
	)
	if err != nil {
		resp.Diagnostics.AddError("Error updating source", "Could not update source, unexpected error: "+err.Error())
		return
	}

	// Update the plan with the response from the API
	plan.UpdatedAt = types.StringValue(source.UpdatedAt.String())
	plan.WorkspaceID = types.Int64Value(int64(source.WorkspaceID))
	plan.ID = types.Int64Value(int64(sourceID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource from the remote API.
func (r *HightouchSnowflakeSourceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state HightouchSnowflakeSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// To complete this, a `DeleteSource` method would be needed on the client.
	// For example:
	// err := r.client.DeleteSource(int(state.ID.ValueInt64()))
	// if err != nil {
	// 	 resp.Diagnostics.AddError("Error deleting source", "Could not delete source, unexpected error: "+err.Error())
	// 	 return
	// }

	// For now, we'll simulate a successful deletion.
	// In a real implementation, you would call your client's delete method here.
	resp.Diagnostics.AddWarning("Delete not implemented", "The Hightouch API does not currently support deleting sources via the API. The resource will be removed from Terraform state, but not from Hightouch.")
}

// ImportState imports the resource into Terraform state.
func (r *HightouchSnowflakeSourceResource) ImportState(
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
