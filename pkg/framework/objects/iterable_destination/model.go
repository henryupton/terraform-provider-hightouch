package iterable_destination

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IterableDestinationResourceModel maps the resource schema data for an Iterable destination in Hightouch.
type IterableDestinationResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Type        types.String `tfsdk:"type"`
	APIKey      types.String `tfsdk:"api_key"`
	DataCenter  types.String `tfsdk:"data_center"`
	WorkspaceID types.Int64  `tfsdk:"workspace_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}
