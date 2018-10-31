// Copyright (C) 2018 Storj Labs, Inc.l "storj.io/ditto/pkg/logger"
// See LICENSE for copying information.

package get

import (
	"github.com/spf13/cobra"
	"storj.io/ditto/cmd/utils"
	l "storj.io/ditto/pkg/logger"
)

// getCmd represents the get command
var Cmd = &cobra.Command{
	Use:   "get [bucket name] [object name](opt) [OPTIONS]",
	Args: validateArgs,
	Short: "Download object and buckets",
	Aliases: []string{"g"},
	Long: ``,
	RunE: NewGetExec(utils.GetGateway, &l.StdOutLogger).runE,
}

var (
	minArg = 1
	maxArg = 2

	nameFlag string
	nameUsage = "Path of the file or folder to be downloaded.\n" +
		"If no objectname provided folder under that name will be created(Default: current folder)"

	prefixFlag string
	prefixUsage = "Used to download part of the bucket that contains specified prefix(Default: \"\")"

	delimiterFlag string
	delimiterUsage = "separates objnames from prefixes(Default: \"/\")"

	recursiveFlag bool
	recursiveUsage = "recursively download content from bucket or prefix"

	forceFlag bool
	forceUsage = "truncates a file if it exists"

	maxKeysFlag int
	maxKeysUsage = "max number of keys list objects returns(Default: 1000)"
)

func init() {
	Cmd.Flags().StringVarP(&nameFlag, "name", "n", "", nameUsage)
	Cmd.Flags().StringVarP(&prefixFlag, "prefix", "p", "", prefixUsage)
	Cmd.Flags().StringVarP(&delimiterFlag, "delimiter", "d", "/", prefixUsage)
	Cmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, recursiveUsage)
	Cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, forceUsage)
	Cmd.Flags().IntVarP(&maxKeysFlag, "keys", "k", 1000, forceUsage)
}