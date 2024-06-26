package config

import (
	"errors"
	"github.com/jorgejr568/cloudflare-cli/cmd/messages"
	"github.com/jorgejr568/cloudflare-cli/internal/config"
	"github.com/spf13/cobra"
	"reflect"
)

func cmdConfigGet(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "List cloudflare config",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			localConfig, err := config.LoadLocalConfig()
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			reflection := reflect.ValueOf(localConfig)
			for i := 0; i < reflection.NumField(); i++ {
				if reflection.Type().Field(i).Tag.Get("json") == args[0] {
					cmd.Println(reflection.Field(i).Interface())
					return
				}
			}

			cmd.PrintErr(messages.ErrorMessage(errors.New("key not found")))
		},
	}
	rootCmd.AddCommand(cmd)
	return nil
}
