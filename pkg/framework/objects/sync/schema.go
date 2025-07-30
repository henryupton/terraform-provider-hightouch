package sync

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var SyncResourceSchema = schema.Schema{
	Description: "Represents a Hightouch Sync, which connects a model to a destination and defines how data flows between them.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the sync.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the sync.",
			Required:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug of the sync.",
			Required:    true,
		},
		"destination_id": schema.Int64Attribute{
			Description: "The ID of the destination for this sync.",
			Required:    true,
		},
		"model_id": schema.Int64Attribute{
			Description: "The ID of the model this sync uses as a data source.",
			Required:    true,
		},
		"source_id": schema.Int64Attribute{
			Description: "The ID of the source (usually inherited from model).",
			Optional:    true,
		},
		"configuration": schema.StringAttribute{
			Description: "JSON configuration for the sync (field mappings, etc.).",
			Required:    true,
		},
		"schedule": schema.StringAttribute{
			Description: "JSON schedule configuration for the sync.",
			Optional:    true,
			Default:     stringdefault.StaticString("{}"),
		},
		"status": schema.StringAttribute{
			Description: "The current status of the sync.",
			Computed:    true,
		},
		"disabled": schema.BoolAttribute{
			Description: "Whether the sync is disabled.",
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"workspace_id": schema.Int64Attribute{
			Description: "The ID of the workspace that the sync belongs to.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "The timestamp when the sync was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			Description: "The timestamp when the sync was last updated.",
			Computed:    true,
		},
	},
}

var SyncDataSourceSchema = datasourceschema.Schema{
	Description: "Fetches information about a Hightouch Sync.",
	Attributes: map[string]datasourceschema.Attribute{
		"id": datasourceschema.Int64Attribute{
			Description: "The ID of the sync.",
			Required:    true,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The name of the sync.",
			Computed:    true,
		},
		"slug": datasourceschema.StringAttribute{
			Description: "The slug of the sync.",
			Computed:    true,
		},
		"destination_id": datasourceschema.Int64Attribute{
			Description: "The ID of the destination for this sync.",
			Computed:    true,
		},
		"model_id": datasourceschema.Int64Attribute{
			Description: "The ID of the model this sync uses as a data source.",
			Computed:    true,
		},
		"source_id": datasourceschema.Int64Attribute{
			Description: "The ID of the source (usually inherited from model).",
			Computed:    true,
		},
		"configuration": datasourceschema.StringAttribute{
			Description: "JSON configuration for the sync (field mappings, etc.).",
			Computed:    true,
		},
		"schedule": datasourceschema.StringAttribute{
			Description: "JSON schedule configuration for the sync.",
			Computed:    true,
		},
		"status": datasourceschema.StringAttribute{
			Description: "The current status of the sync.",
			Computed:    true,
		},
		"disabled": datasourceschema.BoolAttribute{
			Description: "Whether the sync is disabled.",
			Computed:    true,
		},
		"workspace_id": datasourceschema.Int64Attribute{
			Description: "The ID of the workspace that the sync belongs to.",
			Computed:    true,
		},
		"created_at": datasourceschema.StringAttribute{
			Description: "The timestamp when the sync was created.",
			Computed:    true,
		},
		"updated_at": datasourceschema.StringAttribute{
			Description: "The timestamp when the sync was last updated.",
			Computed:    true,
		},
	},
}
