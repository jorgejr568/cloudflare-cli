package dns

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jorgejr568/cloudflare-cli/cmd/messages"
	"github.com/jorgejr568/cloudflare-cli/internal/clients/cloudflare"
	"github.com/jorgejr568/cloudflare-cli/internal/constants"
	"github.com/spf13/cobra"
)

func cmdDnsDelete(rootCmd *cobra.Command, client cloudflare.CloudflareClient) error {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a DNS record",
		Run: func(cmd *cobra.Command, args []string) {
			domain := cmd.Flag(constants.FlagDomain).Value.String()
			zone, err := client.GetZoneByDomain(cmd.Context(), cloudflare.GetZoneByDomainRequest{
				Domain: domain,
			})
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			recordID := cmd.Flag(constants.FlagID).Value.String()
			if recordID == "" {
				if cmd.Flag(constants.FlagName).Changed == false || cmd.Flag(constants.FlagType).Changed == false {
					cmd.PrintErr(messages.ErrorMessage(fmt.Errorf("either --id OR (--name AND --type) must be specified")))
					return
				}

				name := acquireEntryFullName(domain, cmd.Flag(constants.FlagName).Value.String())
				zoneType, err := cloudflare.ParseZoneType(cmd.Flag(constants.FlagType).Value.String())
				if err != nil {
					cmd.PrintErr(messages.ErrorMessage(err))
					return
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

				if len(records.Records) == 0 {
					cmd.PrintErr(messages.ErrorMessage(fmt.Errorf("No records found")))
					return
				}

				if len(records.Records) > 1 {
					cmd.Print(messages.WarningMessage("Multiple records found, please specify the record utilizing the --id flag"))
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
					cmd.Print(messages.WarningMessage("Please specify the record utilizing the --id flag"))
					return
				}

				recordID = records.Records[0].ID
			}

			err = client.DeleteZoneRecord(cmd.Context(), cloudflare.DeleteZoneRecordRequest{
				ZoneID:   zone.ZoneID,
				RecordID: recordID,
			})
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			cmd.Print(messages.SuccessMessage(fmt.Sprintf("Record %s deleted", recordID)))
		},
	}

	cmd.Flags().String(constants.FlagName, "", "The name of the record to delete")
	cmd.Flags().String(constants.FlagType, "", "The type of the record to delete")
	cmd.Flags().String(constants.FlagID, "", "The ID of the record to delete")
	rootCmd.AddCommand(cmd)

	return nil
}
