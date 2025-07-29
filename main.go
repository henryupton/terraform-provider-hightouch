package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-hightouch/pkg/provider"
)

// The 'version' variable is a placeholder that can be populated by build-time flags.
var version string = "dev"

func main() {
	var debug bool

	// The -debug flag allows you to run the provider in a debugger.
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with a debugger")
	flag.Parse()

	// The providerserver.Serve function is the main entrypoint for the provider.
	// It starts a gRPC server that Terraform Core communicates with.
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		// The 'Address' is a unique identifier for your provider in the local registry.
		// This should match the path you specify in your .terraformrc file.
		Address: "local/henryupton/hightouch",
		Debug:   debug,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
