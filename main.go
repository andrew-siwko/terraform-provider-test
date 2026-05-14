package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	fmt.Print("Starting terraform-provider-test\n")
	err := providerserver.Serve(context.Background(), New, providerserver.ServeOpts{
		Address: "registry.terraform.io/andrew/property-mirror",
	})

	if err != nil {
		fmt.Print("Ending terraform-provider-test\n")
		log.Fatal(err.Error())
	}
	fmt.Print("Ending terraform-provider-test\n")
}

