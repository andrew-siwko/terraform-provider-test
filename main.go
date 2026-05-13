package main

import (
	"context"
	"log"

	// "github.com/hashicorp/terraform-plugin-framework/datasource"
	// "github.com/hashicorp/terraform-plugin-framework/provider"
	// provschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	// "github.com/hashicorp/terraform-plugin-framework/providerserver"
	// "github.com/hashicorp/terraform-plugin-framework/resource"
	// resschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/types"
)

func main() {
	// The Serve function and Opts live in the providerserver package in recent versions
	err := providerserver.Serve(context.Background(), New, providerserver.ServeOpts{
		Address: "registry.terraform.io/andrew/property-mirror",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

