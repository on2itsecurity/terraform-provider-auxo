package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/on2itsecurity/terraform-provider-auxo/auxo"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		auxo.New,
		providerserver.ServeOpts{
			Debug:   debug,
			Address: "registry.terraform.io/on2itsecurity/auxo",
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
