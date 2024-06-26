package config

import (
	"github.com/spf13/cobra"
)

func CmdConfig(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}
	err := cmdConfigGet(cmd)
	if err != nil {
		return err
	}
	err = cmdConfigList(cmd)
	if err != nil {
		return err
	}
	err = cmdConfigSet(cmd)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(cmd)
	return nil
}
