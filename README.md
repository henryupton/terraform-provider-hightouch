# Terraform Provider for Hightouch

A Terraform provider for managing [Hightouch](https://hightouch.com) resources. This provider allows you to configure
and manage Hightouch data sources and other resources through Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for development)

## Using the Provider

### Authentication

The provider requires a Hightouch API key. You can configure it in one of two ways:

#### Option 1: Provider Configuration

```hcl
provider "hightouch" {
  api_key      = "your-api-key-here"
  api_base_url = "https://api.hightouch.com/api/v1"  # Optional, defaults to this URL
}
```

#### Option 2: Environment Variables

```bash
export HIGHTOUCH_API_KEY="your-api-key-here"
export HIGHTOUCH_API_BASE_URL="https://api.hightouch.com/api/v1"  # Optional
```

### Example Usage

```hcl
terraform {
  required_providers {
    hightouch = {
      source = "local/henryupton/hightouch"
    }
  }
}

provider "hightouch" {
  # Configuration can be provided here or via environment variables
}

resource "hightouch_snowflake_source" "example" {
  # Snowflake source configuration
  # See resource documentation for available attributes
}
```

## How Hightouch Components Work Together

Hightouch follows a data pipeline architecture where components interact in a specific sequence:

```

[Source] → [Model] → [Destination] → [Sync]

```

### Component Relationships

1. **Sources** connect to your data warehouses (Snowflake, BigQuery, etc.)
   - Define connection credentials and configuration
   - Provide access to raw data tables and views
   - Example: `hightouch_snowflake_source`

2. **Models** define the data transformation and selection logic
   - Reference a specific **Source** via `source_id`
   - Contain SQL queries or dbt model references
   - Define primary keys for data synchronization
   - Act as the "what data to sync" layer

3. **Destinations** define where data should be sent
   - Configure external platforms (Iterable, Salesforce, etc.)
   - Store API credentials and connection settings
   - Example: `hightouch_iterable_destination`

4. **Syncs** orchestrate the data movement
   - Reference a **Model** via `model_id` (what data)
   - Reference a **Destination** via `destination_id` (where to send)
   - Define field mappings and transformation rules
   - Control scheduling and sync frequency
   - Handle the actual data transfer process

### Data Flow Example

```hcl
# 1. Connect to your data warehouse
resource "hightouch_snowflake_source" "warehouse" {
  name     = "Production Warehouse"
  slug     = "prod-warehouse"
  account  = "your-account"
  database = "analytics"
  # ... other config
}

# 2. Define what data to sync using SQL
resource "hightouch_model" "user_segments" {
  name       = "Active Users"
  slug       = "active-users"
  source_id  = hightouch_snowflake_source.warehouse.id
  sql        = "SELECT user_id, email, segment FROM users WHERE active = true"
  primary_key = "user_id"
}

# 3. Configure destination platform
resource "hightouch_iterable_destination" "marketing" {
  name        = "Iterable Marketing"
  slug        = "iterable-marketing"
  api_key     = var.iterable_api_key
  data_center = "US"
}

# 4. Create sync to move data from model to destination
resource "hightouch_sync" "user_sync" {
  name           = "Sync Active Users to Iterable"
  slug           = "sync-active-users"
  model_id       = hightouch_model.user_segments.id
  destination_id = hightouch_iterable_destination.marketing.id
  
  configuration = {
    # Field mappings and sync settings
  }
  
  schedule = {
    type = "interval"
    interval = 3600  # Every hour
  }
}
```

### Key Concepts

- **Sources** are reusable across multiple models
- **Models** can be used by multiple syncs
- **Destinations** can receive data from multiple syncs
- **Syncs** represent the actual data pipeline execution
- All components exist within a **Workspace** context

This architecture provides flexibility to:

- Share data connections across teams
- Create multiple data views from the same source
- Send the same data to multiple destinations
- Control sync timing and frequency independently

## Available Resources

- `hightouch_snowflake_source` - Manages Snowflake data sources in Hightouch
- `hightouch_iterable_destination` - Manages Iterable destinations in Hightouch

## Available Data Sources

- `data.hightouch_snowflake_source` - Fetches information about existing Snowflake sources
- `data.hightouch_iterable_destination` - Fetches information about existing Iterable destinations

## Development

### Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider:

```bash
go build -o terraform-provider-hightouch
```

### Installing for Local Development

To use the provider locally, you'll need to create a local provider registry. Create or modify your `~/.terraformrc`
file:

```hcl
provider_installation {
  dev_overrides {
    "local/henryupton/hightouch" = "/path/to/your/terraform-provider-hightouch"
  }
  direct {}
}
```

Replace `/path/to/your/terraform-provider-hightouch` with the actual path to your built provider binary.

### Running Tests

```bash
# Run the test lifecycle script
./scripts/test-lifecycle.sh
```

### Debugging

You can run the provider in debug mode:

```bash
go run main.go -debug
```

This will output instructions for configuring Terraform to connect to the debugging provider.

## API Documentation

This provider interacts with the Hightouch API. For more information about available endpoints and data models, refer to
the [Hightouch API documentation](https://hightouch.com/docs/api).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is released into the public domain under The Unlicense. See the [LICENSE](LICENSE) file for details.