package snowflake_source

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var datasourceSchema = schema.Schema{
	Description: "Represents a Hightouch Source, which is a connector to a data warehouse, database, or other data platform.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the source.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the source.",
			Required:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug of the source.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the source, 'snowflake'.",
			Computed:    true,
			Default:     stringdefault.StaticString("snowflake"),
		},
		"account": schema.StringAttribute{
			Description: "Source account.",
			Required:    true,
		},
		"port": schema.Int64Attribute{
			Description: "Source port.",
			//Default:     int64default.StaticInt64(443),
			Optional: true,
			//Computed:    true,
		},
		"username": schema.StringAttribute{
			Description: "Username.",
			Required:    true,
		},
		"password": schema.StringAttribute{
			Description: "Password.",
			Required:    true,
			Sensitive:   true, // Mark as sensitive to avoid logging
		},
		"database": schema.StringAttribute{
			Description: "Database name.",
			Required:    true,
		},
		"warehouse": schema.StringAttribute{
			Description: "Warehouse name (if applicable, e.g., for Snowflake).",
			Required:    true, // Optional if not applicable to all source types
		},
		"workspace_id": schema.Int64Attribute{
			Description: "The ID of the workspace that the source belongs to.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "The timestamp when the source was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			Description: "The timestamp when the source was last updated.",
			Computed:    true,
		},
	},
}
