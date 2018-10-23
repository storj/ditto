// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"storj.io/ditto/pkg/config"
)

var getSubCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Displays value set for requested key",
	Long:  "Displays value set for requested key",
	RunE:  executeGetCmd,
	Args:  validateGetArgs,
}

// Method reference for unit testing
var readConfigMethod = config.ReadConfig

func executeGetCmd(cmd *cobra.Command, args []string) error {
	arg := args[0]

	readConfigMethod(false)
	if containsKey(arg) {
		fmt.Sprintf("\t%s\n", getValueFromConfigFile(arg))
		return nil
	} else {
		return errors.New("Key unsupported")
	}
}

func validateGetArgs(cmd *cobra.Command, args [] string) error {
	if len(args) != 1 {
		return errors.New("Only one argument accepted")
	}
	return nil

}

func getValueFromConfigFile(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	} else {
		return "Key is not set"
	}
}

func init() {
	Cmd.AddCommand(getSubCmd)
}
