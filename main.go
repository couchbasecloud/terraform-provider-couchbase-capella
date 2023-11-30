// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/couchbasecloud/couchbase-capella",
		Debug:   debug,
	}

	err := providerserver.Serve(
		context.Background(),
		provider.New(),
		opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
