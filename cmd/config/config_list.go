// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"storj/ditto/pkg/config"
)

var listSubCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays list of possible to change options",
	Long:  "Displays list of possible to change options",
	Run:   executeListCmd,
}

func executeListCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Options, which can be set via `config set`:")
	for _, value := range config.GetKeysArray() {
		fmt.Printf("\t%s\n", value)
	}
}

func init() {
	Cmd.AddCommand(listSubCmd)
}
