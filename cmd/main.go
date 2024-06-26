package cmd

import (
	cmdconfig "github.com/jorgejr568/cloudflare-cli/cmd/config"
	"github.com/jorgejr568/cloudflare-cli/cmd/dns"
	"github.com/jorgejr568/cloudflare-cli/internal/clients/cloudflare"
	"github.com/jorgejr568/cloudflare-cli/internal/config"
	"github.com/jorgejr568/cloudflare-cli/internal/constants"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"log"
	"net/http"
)

func Run() {
	container := dig.New()

	err := container.Provide(func() (config.Composite, error) {
		return config.NewComposite()
	})
	if err != nil {
		log.Fatalf("failed to load config composite: %v", err)
	}
	err = container.Provide(func() (*http.Client, error) {
		return &http.Client{}, nil
	})
	if err != nil {
		log.Fatalf("failed to load http client: %v", err)
	}
	err = container.Provide(func(client *http.Client, config config.Composite) (cloudflare.CloudflareClient, error) {
		return cloudflare.NewHttpCloudflareClient(
			client,
			config.CloudflareAPIKey(),
			constants.CloudflareAPIBaseURL,
		), nil
	})
	if err != nil {
		log.Fatalf("failed to load cloudflare client: %v", err)
	}

	err = container.Provide(cmdRoot)
	if err != nil {
		log.Fatalf("failed to load root command: %v", err)
	}

	err = container.Invoke(dns.CmdDns)
	if err != nil {
		log.Fatalf("failed to load dns commands: %v", err)
	}
	err = container.Invoke(cmdconfig.CmdConfig)
	if err != nil {
		log.Fatalf("failed to load config commands: %v", err)
	}

	_ = container.Invoke(func(cmd *cobra.Command) {
		if err := cmd.Execute(); err != nil {
			log.Fatal(err)
		}
	})
}
