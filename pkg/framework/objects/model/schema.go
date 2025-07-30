package model

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var ModelResourceSchema = schema.Schema{
	Description: "Represents a Hightouch Model, which defines how data is selected and transformed from a source.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the model.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the model.",
			Required:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug of the model.",
			Required:    true,
		},
		"source_id": schema.Int64Attribute{
			Description: "The ID of the source this model queries from.",
			Required:    true,
		},
		"sql": schema.StringAttribute{
			Description: "The SQL query that defines the model.",
			Required:    true,
		},
		"dbt_table": schema.StringAttribute{
			Description: "The dbt table name if using dbt.",
			Optional:    true,
		},
		"query_type": schema.StringAttribute{
			Description: "The type of query (e.g., 'sql', 'dbt').",
			Optional:    true,
			Default:     stringdefault.StaticString("sql"),
		},
		"primary_key": schema.StringAttribute{
			Description: "The primary key column for the model.",
			Required:    true,
		},
		"is_schema": schema.BoolAttribute{
			Description: "Whether this model represents a schema.",
			Optional:    true,
		},
		"workspace_id": schema.Int64Attribute{
			Description: "The ID of the workspace that the model belongs to.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "The timestamp when the model was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			Description: "The timestamp when the model was last updated.",
			Computed:    true,
		},
	},
}

var ModelDataSourceSchema = datasourceschema.Schema{
	Description: "Fetches information about a Hightouch Model.",
	Attributes: map[string]datasourceschema.Attribute{
		"id": datasourceschema.Int64Attribute{
			Description: "The ID of the model.",
			Required:    true,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The name of the model.",
			Computed:    true,
		},
		"slug": datasourceschema.StringAttribute{
			Description: "The slug of the model.",
			Computed:    true,
		},
		"source_id": datasourceschema.Int64Attribute{
			Description: "The ID of the source this model queries from.",
			Computed:    true,
		},
		"sql": datasourceschema.StringAttribute{
			Description: "The SQL query that defines the model.",
			Computed:    true,
		},
		"dbt_table": datasourceschema.StringAttribute{
			Description: "The dbt table name if using dbt.",
			Computed:    true,
		},
		"query_type": datasourceschema.StringAttribute{
			Description: "The type of query (e.g., 'sql', 'dbt').",
			Computed:    true,
		},
		"primary_key": datasourceschema.StringAttribute{
			Description: "The primary key column for the model.",
			Computed:    true,
		},
		"is_schema": datasourceschema.BoolAttribute{
			Description: "Whether this model represents a schema.",
			Computed:    true,
		},
		"workspace_id": datasourceschema.Int64Attribute{
			Description: "The ID of the workspace that the model belongs to.",
			Computed:    true,
		},
		"created_at": datasourceschema.StringAttribute{
			Description: "The timestamp when the model was created.",
			Computed:    true,
		},
		"updated_at": datasourceschema.StringAttribute{
			Description: "The timestamp when the model was last updated.",
			Computed:    true,
		},
	},
}
