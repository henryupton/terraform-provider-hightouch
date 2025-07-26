package provider

import (
    "fmt"
    "strconv"
    "context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-hightouch/hightouch"
)

// resourceHightouchSource defines the schema for the Hightouch source resource.
func ResourceHightouchSource() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Hightouch source.",
		CreateContext: resourceHightouchSourceCreate,
		ReadContext:   resourceHightouchSourceRead,
		UpdateContext: resourceHightouchSourceUpdate,
		DeleteContext: resourceHightouchSourceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the source.",
			},
			"slug": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the source.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the source (e.g., 'postgres').",
			},
			// To represent the nested JSON object, we use a TypeList with MaxItems: 1.
			// This creates a nested block in the Terraform configuration.
// 			"configuration": {
// 				Type:        schema.TypeList,
// 				Required:    true,
// 				MaxItems:    1,
// 				Description: "Connection details for the source.",
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"host": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Description: "Database host.",
// 						},
// 						"port": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Description: "Database port.",
// 						},
// 						"user": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Description: "Database user.",
// 						},
// 						"database": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Description: "Database name.",
// 						},
// 						"password": {
// 							Type:        schema.TypeString,
// 							Required:    true,
// 							Sensitive:   true, // Mark as sensitive to hide from output.
// 							Description: "Database password.",
// 						},
// 					},
// 				},
// 			},
		},
	}
}

// Placeholder CRUD functions for the hightouch_source resource.
func resourceHightouchSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// In a real provider, you would make an API call to Hightouch here.
	// We'll use the slug as the resource ID.
	slug := d.Get("slug").(string)
	d.SetId(slug)

	// Example payload for creating a PostgreSQL source.
	// The 'configuration' will vary depending on the source 'type'.
	// Refer to the Hightouch API documentation for the required fields for your specific source.
	payload := &hightouch.CreateSourcePayload{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
		Type: d.Get("type").(string),
		Configuration: map[string]interface{}{
			"host":     d.Get("host").(string),
			"port":     d.Get("port").(int),
			"user":     d.Get("user").(string),
			"database": d.Get("database").(string),
			"password": d.Get("password").(string),
// 		},
	}

    fmt.Println(payload)

	fmt.Println("Attempting to create a new source...")
	client := hightouch.NewClient(meta.(string)) // Assuming meta contains the API key
	source, err := client.CreateSource(payload)
	fmt.Println(source, err)

	return resourceHightouchSourceRead(ctx, d, meta)
}

func resourceHightouchSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    i, err := strconv.Atoi(d.Id())
    if err != nil {
        panic(err)
    }

    client := hightouch.NewClient(meta.(string)) // Assuming meta contains the API key
	source, err := client.GetSource(i)
	fmt.Println(source)

	return nil
}

func resourceHightouchSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Update logic would go here.
	return resourceHightouchSourceRead(ctx, d, meta)
}

func resourceHightouchSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Deletion logic would go here.
	return nil
}
