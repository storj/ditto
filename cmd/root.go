// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"github.com/spf13/cobra"
	"storj/ditto/cmd/config"
	"storj/ditto/cmd/cp"
	"storj/ditto/cmd/delete"
	"storj/ditto/cmd/get"
	"storj/ditto/cmd/list"
	"storj/ditto/cmd/mb"
	"storj/ditto/cmd/put"
	"storj/ditto/cmd/server"
	"storj/ditto/cmd/version"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ditto",
	Short: "A backup mirroring util",
	Long:  `A backup mirroring util`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	RootCmd.AddCommand(mb.Cmd)
	RootCmd.AddCommand(cp.Cmd)
	RootCmd.AddCommand(put.Cmd)
	RootCmd.AddCommand(get.Cmd)
	RootCmd.AddCommand(list.Cmd)
	RootCmd.AddCommand(delete.Cmd)
	RootCmd.AddCommand(version.Cmd)
	RootCmd.AddCommand(config.Cmd)
	RootCmd.AddCommand(server.Cmd)

	RootCmd.Execute()
}
