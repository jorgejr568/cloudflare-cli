package config

import (
	"fmt"
	"github.com/jorgejr568/cloudflare-cli/cmd/messages"
	"github.com/jorgejr568/cloudflare-cli/internal/config"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
)

func cmdConfigSet(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "set [<key> <value>...]",
		Short: "Set cloudflare config",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args)%2 != 0 {
				cmd.PrintErr(messages.ErrorMessage(fmt.Errorf("invalid number of args")))
				return
			}

			localConfig, err := config.LoadLocalConfig()
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			for i := 0; i < len(args); i += 2 {
				key := args[i]
				value := args[i+1]

				err = setField(&localConfig, key, value)
				if err != nil {
					cmd.PrintErr(messages.ErrorMessage(err))
					return
				}
			}

			err = config.SaveLocalConfig(localConfig)
			if err != nil {
				cmd.PrintErr(messages.ErrorMessage(err))
				return
			}

			cmd.Print(messages.SuccessMessage("Config updated"))
		},
	}
	rootCmd.AddCommand(cmd)
	return nil
}

func parseArg(arg string) (string, string, error) {
	parts := strings.Split(arg, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid arg format: %s", arg)
	}

	return parts[0], parts[1], nil
}

func setField(localConfig *config.LocalConfig, key string, value string) error {
	reflection := reflect.ValueOf(localConfig)
	for i := 0; i < reflection.Elem().NumField(); i++ {
		field := reflection.Elem().Field(i)
		tag := reflection.Elem().Type().Field(i).Tag.Get("json")
		if tag == key {
			if !field.CanSet() {
				return fmt.Errorf("cannot set field %s", key)
			}

			field.SetString(value)
			return nil
		}
	}

	return fmt.Errorf("field %s not found", key)
}
