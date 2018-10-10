// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"storj.io/ditto/cmd/config/config_utils"
	"storj.io/ditto/pkg/config"
)

var getSubCmd = &cobra.Command{
	Use:   "get",
	Short: "Mirroring options setup",
	Long:  `getSub`,

	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]

		config.ReadDefaultConfig(false)
		if config_utils.ContainsKey(arg) {
			getValueFromConfigFile(arg)
		} else {
			fmt.Println("\tKey unsupported")
		}
	},

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Only one argument accepted")
		}
		return nil
	},
}

func getValueFromConfigFile(key string) {
	if viper.IsSet(key) {
		fmt.Printf("\t%s\n", viper.GetString(key))
	} else {
		fmt.Println("\tUNDEFINED")
	}
}

func init() {
	Cmd.AddCommand(getSubCmd)
}
