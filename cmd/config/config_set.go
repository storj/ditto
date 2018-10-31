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

var setSubCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Change value at saves it to config file",
	Long:  "Change value at saves it to config file",
	RunE:  executeSetCmd,
	Args:  validateSetArgs,
}

var writeConfigMethod = writeConfig

func executeSetCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Two arguments expected")
	}
	possibleValues, exist := argsMap[args[0]]
	if !exist {

		return errors.New("Key is not exist")
	}

	if len(possibleValues) == 0 {
		return writeConfigMethod(args[0], args[1])
	} else {
		return setPredefinedOptions(args[0], args[1], possibleValues)
	}
}

func validateSetArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Two arguments expected")
	}

	return nil
}

func setPredefinedOptions(key string, value string, possibleValues []string) error {
	if contains(possibleValues, value) {
		return writeConfigMethod(key, value)
	} else {
		return errors.New(fmt.Sprintf("Only these arguments accepted: %s", possibleValues))
	}
}

func writeConfig(key string, value string) error {
	config.ParseConfig()

	viper.Set(key, value)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil

}

func init() {
	Cmd.AddCommand(setSubCmd)
}
