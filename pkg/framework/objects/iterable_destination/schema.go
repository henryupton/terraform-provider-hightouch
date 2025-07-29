package iterable_destination

import (
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var IterableDestinationResourceSchema = schema.Schema{
	Description: "Represents a Hightouch Iterable Destination, which is a connector to send data to Iterable for marketing campaigns.",
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the destination.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the destination.",
			Required:    true,
		},
		"slug": schema.StringAttribute{
			Description: "The slug of the destination.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the destination, 'iterable'.",
			Computed:    true,
			Default:     stringdefault.StaticString("iterable"),
		},
		"api_key": schema.StringAttribute{
			Description: "The Iterable API key for authentication.",
			Required:    true,
			Sensitive:   true, // Mark as sensitive to avoid logging
		},
		"data_center": schema.StringAttribute{
			Description: "The Iterable data center (US or EU).",
			Optional:    true,
			Default:     stringdefault.StaticString("US"),
		},
		"workspace_id": schema.Int64Attribute{
			Description: "The ID of the workspace that the destination belongs to.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "The timestamp when the destination was created.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			Description: "The timestamp when the destination was last updated.",
			Computed:    true,
		},
	},
}

var IterableDestinationDataSourceSchema = datasourceschema.Schema{
	Description: "Fetches information about a Hightouch Iterable Destination.",
	Attributes: map[string]datasourceschema.Attribute{
		"id": datasourceschema.Int64Attribute{
			Description: "The ID of the destination.",
			Required:    true,
		},
		"name": datasourceschema.StringAttribute{
			Description: "The name of the destination.",
			Computed:    true,
		},
		"slug": datasourceschema.StringAttribute{
			Description: "The slug of the destination.",
			Computed:    true,
		},
		"type": datasourceschema.StringAttribute{
			Description: "The type of the destination.",
			Computed:    true,
		},
		"api_key": datasourceschema.StringAttribute{
			Description: "The Iterable API key for authentication.",
			Computed:    true,
			Sensitive:   true,
		},
		"data_center": datasourceschema.StringAttribute{
			Description: "The Iterable data center (US or EU).",
			Computed:    true,
		},
		"workspace_id": datasourceschema.Int64Attribute{
			Description: "The ID of the workspace that the destination belongs to.",
			Computed:    true,
		},
		"created_at": datasourceschema.StringAttribute{
			Description: "The timestamp when the destination was created.",
			Computed:    true,
		},
		"updated_at": datasourceschema.StringAttribute{
			Description: "The timestamp when the destination was last updated.",
			Computed:    true,
		},
	},
}
