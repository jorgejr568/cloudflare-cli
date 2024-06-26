package cmd

import "github.com/spf13/cobra"

func cmdRoot() (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "cloudflare-cli",
		Short: "A CLI for interacting with the Cloudflare API",
	}, nil
}
