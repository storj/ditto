// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"storj.io/ditto/cmd/cp"
	"storj.io/ditto/cmd/get"
	"storj.io/ditto/cmd/list"
	"storj.io/ditto/cmd/make_bucket"
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
	Use:   "mirroring",
	Short: "A backup mirroring util",
	Long: `A backup mirroring util`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//rootCmd.
	addCommands()
	//TODO set persistent flags
	//TODO Bind flags with viper

	rootCmd.Execute()
	//if err := rootCmd.Execute(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
}
func addCommands() {
	rootCmd.AddCommand(make_bucket.Cmd)
	rootCmd.AddCommand(cp.Cmd)
	rootCmd.AddCommand(put.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(list.Cmd)
	//rootCmd.AddCommand(delete.Cmd)
	rootCmd.AddCommand(version.Cmd)
	//rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(server.Cmd)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mirroring/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.Set("configPath", cfgFile)
	}
}
