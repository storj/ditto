// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Ditto options setup",
	Long:  `Gives ability to get value set for particular key. List all options, or change it.`,
}
