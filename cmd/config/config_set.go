// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

var setSubCmd = &cobra.Command{
	Use:   "set",
	Short: "Mirroring options setup",
	Long:  `listSub`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list_cmd sub cmd called")
	},
}

func init() {
	Cmd.AddCommand(setSubCmd)
}
