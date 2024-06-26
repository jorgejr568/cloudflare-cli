package dns

import (
	"github.com/jorgejr568/cloudflare-cli/internal/clients/cloudflare"
	"github.com/jorgejr568/cloudflare-cli/internal/constants"
	"github.com/jorgejr568/cloudflare-cli/internal/utils"
	"github.com/spf13/cobra"
)

func CmdDns(rootCmd *cobra.Command, client cloudflare.CloudflareClient) error {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Manage DNS records",
	}
	cmd.PersistentFlags().StringP(constants.FlagDomain, "d", "", "The domain to list DNS records for")
	utils.LogFatalIfError(cmd.MarkPersistentFlagRequired(constants.FlagDomain))
	err := cmdDnsList(cmd, client)
	if err != nil {
		return err
	}

	err = cmdDnsAdd(cmd, client)
	if err != nil {
		return err
	}

	err = cmdDnsDelete(cmd, client)
	if err != nil {
		return err
	}
	rootCmd.AddCommand(cmd)
	return nil
}
