package snowflake_source

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SnowflakeSourceResourceModel maps the resource schema data for a Snowflake source in Hightouch.
type SnowflakeSourceResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Type        types.String `tfsdk:"type"`
	Account     types.String `tfsdk:"account"`
	Port        types.Int64  `tfsdk:"port"`
	Username    types.String `tfsdk:"username"`
	Database    types.String `tfsdk:"database"`
	Password    types.String `tfsdk:"password"`
	Warehouse   types.String `tfsdk:"warehouse"`
	WorkspaceID types.Int64  `tfsdk:"workspace_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}
