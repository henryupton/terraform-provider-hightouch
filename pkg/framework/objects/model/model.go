package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ModelResourceModel maps the resource schema data for a Hightouch model.
type ModelResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	SourceID    types.Int64  `tfsdk:"source_id"`
	SQL         types.String `tfsdk:"sql"`
	DBTable     types.String `tfsdk:"dbt_table"`
	QueryType   types.String `tfsdk:"query_type"`
	PrimaryKey  types.String `tfsdk:"primary_key"`
	IsSchema    types.Bool   `tfsdk:"is_schema"`
	WorkspaceID types.Int64  `tfsdk:"workspace_id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}
