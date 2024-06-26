package dns

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jorgejr568/cloudflare-cli/cmd/messages"
	"github.com/jorgejr568/cloudflare-cli/internal/clients/cloudflare"
	"github.com/jorgejr568/cloudflare-cli/internal/constants"
	"github.com/jorgejr568/cloudflare-cli/internal/utils"
	"github.com/spf13/cobra"
)

func cmdDnsAdd(rootCmd *cobra.Command, client cloudflare.CloudflareClient) error {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a DNS record",
		Run: func(cmd *cobra.Command, args []string) {
			domain := cmd.Flag(constants.FlagDomain).Value.String()
			ttl, err := cmd.Flags().GetInt(constants.FlagTTL)
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}
			if ttl < 1 || ttl > 86400 {
				cmd.PrintErr("TTL must be between 1 and 86400")
				return
			}

			proxied, err := cmd.Flags().GetBool(constants.FlagProxied)
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}
			tags, err := cmd.Flags().GetStringSlice(constants.FlagTags)
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			zoneType, err := cloudflare.ParseZoneType(cmd.Flag(constants.FlagType).Value.String())
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			record := cloudflare.ZoneRecordRequest{
				Type:    zoneType,
				Name:    acquireEntryFullName(domain, cmd.Flag(constants.FlagName).Value.String()),
				Content: cmd.Flag(constants.FlagContent).Value.String(),
				TTL:     ttl,
				Proxied: proxied,
				Tags:    tags,
				Comment: cmd.Flag(constants.FlagComment).Value.String(),
			}

			zone, err := client.GetZoneByDomain(cmd.Context(), cloudflare.GetZoneByDomainRequest{
				Domain: domain,
			})
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			response, err := client.AddZoneRecord(cmd.Context(), cloudflare.AddZoneRecordRequest{
				ZoneID: zone.ZoneID,
				Record: record,
			})
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"ID", "Type", "Name", "Content", "Proxied", "TTL", "Tags", "Comment"})
			t.AppendRow(table.Row{
				response.Record.ID,
				response.Record.Type,
				response.Record.Name,
				response.Record.Content,
				response.Record.Proxied,
				response.Record.TTL,
				response.Record.Tags,
				response.Record.Comment,
			})
			t.Render()
		},
	}

	cmd.Flags().String(constants.FlagName, "", "Name of the record")
	cmd.Flags().String(constants.FlagType, "", "Type of the record")
	cmd.Flags().String(constants.FlagContent, "", "Content of the record")
	cmd.Flags().Int(constants.FlagTTL, 1, "TTL of the record (1-86400)")
	cmd.Flags().Bool(constants.FlagProxied, false, "Proxied status of the record")
	cmd.Flags().StringSlice(constants.FlagTags, []string{}, "Tags of the record")
	cmd.Flags().String(constants.FlagComment, "", "Comment of the record")

	utils.LogFatalIfError(cmd.MarkFlagRequired(constants.FlagName))
	utils.LogFatalIfError(cmd.MarkFlagRequired(constants.FlagType))
	utils.LogFatalIfError(cmd.MarkFlagRequired(constants.FlagContent))

	rootCmd.AddCommand(cmd)
	return nil
}
