#!/bin/bash

# =================================================================
# Bash Script to Build and Install a Local Terraform Provider
# =================================================================
# This script automates the steps outlined in the guide to compile
# a Go-based Terraform provider and set it up for local testing.

#set -e # Exit immediately if a command exits with a non-zero status.

# --- Configuration ---
# You can change these variables to match your provider's details.
PROVIDER_NAME="hightouch"
PROVIDER_NAMESPACE="henryupton"
PROVIDER_SOURCE="local/${PROVIDER_NAMESPACE}/${PROVIDER_NAME}"
PROVIDER_VERSION="1.0.0"
BINARY_NAME="terraform-provider-${PROVIDER_NAME}"
MAIN_GO_FILE="main.go" # The entrypoint Go file for your provider.
TEST_DIR="terraform-test-hightouch" # A directory to store the test .tf file.

# --- 1. Check for Go and Go file ---
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go to build the provider."
    exit 1
fi

if [ ! -f "$MAIN_GO_FILE" ]; then
    echo "Error: Main provider file '$MAIN_GO_FILE' not found."
    echo "Please run this script in the same directory as your provider's Go source code."
    exit 1
fi


# --- 2. Initialize Go Module ---
echo "--> Initializing Go module..."
go mod init "terraform-provider-${PROVIDER_NAME}"
go mod tidy


# --- 3. Build the Provider Binary ---
echo "--> Building the provider binary..."
go build -o "${BINARY_NAME}"


# --- 4. Determine OS and Architecture ---
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
echo "--> Detected OS: ${OS}, Arch: ${ARCH}"


# --- 5. Create Local Plugin Directory ---
INSTALL_PATH="${HOME}/.terraform.d/plugins/${PROVIDER_SOURCE}/${PROVIDER_VERSION}/${OS}_${ARCH}"
echo "--> Creating installation directory: ${INSTALL_PATH}"
mkdir -p "${INSTALL_PATH}"


# --- 6. Move Binary to Plugin Directory ---
echo "--> Installing provider binary..."
mv "${BINARY_NAME}" "${INSTALL_PATH}/"


# --- 7. Create Example Terraform Configuration ---
echo "--> Creating example test configuration in ./${TEST_DIR}/"
mkdir -p "${TEST_DIR}"
cat > "${TEST_DIR}/main.tf" << EOF
terraform {
  required_providers {
    ${PROVIDER_NAME} = {
      source  = "${PROVIDER_SOURCE}"
      version = "${PROVIDER_VERSION}"
    }
  }
}

# Configure the provider.
provider "${PROVIDER_NAME}" {
  api_key = "your-api-key-here"
}

resource "${PROVIDER_NAME}_source" "this" {
  name    = "My Test Server"
  slug    = "my-test-server"
  type = "snowflake"

  configuration {
    host = "abc12345.us-east-1.snowflakecomputing.com"
    port = 443
    user = "your_username"
    password = "your_password"
    database = "your_database"
  }
}

output "source_name" {
  value = ${PROVIDER_NAME}_source.this.name
}
EOF

# --- Success Message ---
echo ""
echo "✅ Success! The Terraform provider has been built and installed locally."
echo ""
echo "Testing with 'terraform plan':"
cd "${TEST_DIR}"
terraform init
terraform plan -out plan.tfplan
terraform apply plan.tfplan
cd ..
rm -rf "${TEST_DIR}" # Clean up the test directory after testing
