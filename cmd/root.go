// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"storj.io/ditto/cmd/config"
	"storj.io/ditto/cmd/cp"
	"storj.io/ditto/cmd/delete"
	"storj.io/ditto/cmd/get"
	"storj.io/ditto/cmd/list"
	"storj.io/ditto/cmd/mb"
	"storj.io/ditto/cmd/put"
	"storj.io/ditto/cmd/server"
	"storj.io/ditto/cmd/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "storj.io/ditto/pkg/gateway"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ditto",
	Short: "A backup mirroring util",
	Long:  `A backup mirroring util`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	addCommands()
	rootCmd.Execute()
}

func addCommands() {
	rootCmd.AddCommand(mb.Cmd)
	rootCmd.AddCommand(cp.Cmd)
	rootCmd.AddCommand(put.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(delete.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(server.Cmd)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mirroring/config.json)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.Set("configPath", cfgFile)
	}
}
