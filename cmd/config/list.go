package config

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jorgejr568/cloudflare-cli/cmd/messages"
	"github.com/jorgejr568/cloudflare-cli/internal/config"
	"github.com/spf13/cobra"
	"reflect"
)

func cmdConfigList(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cloudflare config",
		Run: func(cmd *cobra.Command, args []string) {
			localConfig, err := config.LoadLocalConfig()
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			reflection := reflect.ValueOf(localConfig)
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Key", "Value"})
			for i := 0; i < reflection.NumField(); i++ {
				t.AppendRow(table.Row{
					reflection.Type().Field(i).Tag.Get("json"),
					reflection.Field(i).Interface(),
				})
			}

			t.Render()
		},
	}
	rootCmd.AddCommand(cmd)
	return nil
}
