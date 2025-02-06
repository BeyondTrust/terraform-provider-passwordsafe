package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"

	ephemeralProvider "terraform-provider-passwordsafe/ephemeral_provider"
	provider "terraform-provider-passwordsafe/provider"

	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// terraform-plugin-framework provider (ephemeral resoruces)
		providerserver.NewProtocol5(
			ephemeralProvider.NewProvider(),
		),
		// terraform-plugin-sdk/v2 provider ()
		provider.Provider().GRPCProvider,
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	err = tf5server.Serve(
		"registry.terraform.io/providers/BeyondTrust/passwordsafe",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
