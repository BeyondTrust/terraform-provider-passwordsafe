package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"

	providerFramework "terraform-provider-passwordsafe/providers/provider_framework" // new version of provider, ephemeral resources implemented. (terraform-plugin-framework)
	providerSdkv2 "terraform-provider-passwordsafe/providers/provider_sdkv2"         // first version of provider. (terraform-plugin-sdk/v2)

	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	ctx := context.Background()
	providers := []func() tfprotov5.ProviderServer{
		// terraform-plugin-framework provider (ephemeral resoruces), new provider
		providerserver.NewProtocol5(
			providerFramework.NewProvider(),
		),

		// terraform-plugin-sdk/v2 provider, first provider we developed.
		providerSdkv2.Provider().GRPCProvider,
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
