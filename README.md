# Terraform Provider for Hightouch

A Terraform provider for managing [Hightouch](https://hightouch.com) resources. This provider allows you to configure and manage Hightouch data sources and other resources through Terraform.

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

## Available Resources

- `hightouch_snowflake_source` - Manages Snowflake data sources in Hightouch

## Development

### Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider:

```bash
go build -o terraform-provider-hightouch
```

### Installing for Local Development

To use the provider locally, you'll need to create a local provider registry. Create or modify your `~/.terraformrc` file:

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

This provider interacts with the Hightouch API. For more information about available endpoints and data models, refer to the [Hightouch API documentation](https://hightouch.com/docs/api).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is released into the public domain under The Unlicense. See the [LICENSE](LICENSE) file for details.