package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	log.Print('Starting terraform-provider-test')
	err := providerserver.Serve(context.Background(), New, providerserver.ServeOpts{
		Address: "registry.terraform.io/andrew/property-mirror",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print('Ending terraform-provider-test')
}

