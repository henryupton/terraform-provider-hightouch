package sync

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SyncResourceModel maps the resource schema data for a Hightouch sync.
type SyncResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Slug          types.String `tfsdk:"slug"`
	DestinationID types.Int64  `tfsdk:"destination_id"`
	ModelID       types.Int64  `tfsdk:"model_id"`
	SourceID      types.Int64  `tfsdk:"source_id"`
	Configuration types.String `tfsdk:"configuration"`
	Schedule      types.String `tfsdk:"schedule"`
	Status        types.String `tfsdk:"status"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	WorkspaceID   types.Int64  `tfsdk:"workspace_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}
