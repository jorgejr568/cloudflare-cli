package dns

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jorgejr568/cloudflare-cli/cmd/messages"
	"github.com/jorgejr568/cloudflare-cli/internal/clients/cloudflare"
	"github.com/jorgejr568/cloudflare-cli/internal/constants"
	"github.com/spf13/cobra"
)

func cmdDnsList(rootCmd *cobra.Command, client cloudflare.CloudflareClient) error {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all DNS records",
		Run: func(cmd *cobra.Command, args []string) {
			domain := cmd.Flag(constants.FlagDomain).Value.String()
			zone, err := client.GetZoneByDomain(cmd.Context(), cloudflare.GetZoneByDomainRequest{
				Domain: domain,
			})
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			name := cmd.Flag(constants.FlagName).Value.String()
			if name != "" {
				name = acquireEntryFullName(domain, name)
			}
			zoneType := cloudflare.ZoneType("")
			if cmd.Flag(constants.FlagType).Changed {
				zoneType, err = cloudflare.ParseZoneType(cmd.Flag(constants.FlagType).Value.String())
				if err != nil {
					cmd.PrintErr(messages.ErrorMessage(err))
					return
				}
			}

			records, err := client.GetZoneRecords(cmd.Context(), cloudflare.GetZoneRecordsRequest{
				ZoneID: zone.ZoneID,
				Name:   name,
				Type:   zoneType,
			})
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"ID", "Type", "Name", "Content", "Proxied", "TTL", "Tags", "Comment"})
			for _, record := range records.Records {
				t.AppendRow(table.Row{
					record.ID,
					record.Type,
					record.Name,
					record.Content,
					record.Proxied,
					record.TTL,
					record.Tags,
					record.Comment,
				})
			}
			t.Render()
		},
	}

	cmd.Flags().String(constants.FlagName, "", "The name of the record to list")
	cmd.Flags().String(constants.FlagType, "", "The type of the record to list")

	rootCmd.AddCommand(cmd)
	return nil
}
