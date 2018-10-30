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
	Short: "Download files and buckets",
	Long: ``,
	RunE: NewGetExec(utils.GetGateway, &l.StdOutLogger).runE,
}

var (
	minArg = 1
	maxArg = 2

	nameFlag string
	nameUsage = "Path of the file or folder to be downloaded. A raw filename can be used to download to current directory under that name.\n" +
		"If no objectname provided folder under that name will be created"

	prefixFlag string
	prefixUsage = ""

	delimiterFlag string
	delimiterUsage = ""

	recursiveFlag bool
	recursiveUsage = ""

	forceFlag bool
	forceUsage = ""

	maxKeysFlag int
	maxKeysUsage = ""
)

func init() {
	Cmd.Flags().StringVarP(&nameFlag, "name", "n", "", nameUsage)
	Cmd.Flags().StringVarP(&prefixFlag, "prefix", "p", "", prefixUsage)
	Cmd.Flags().StringVarP(&delimiterFlag, "delimiter", "d", "/", prefixUsage)
	Cmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, recursiveUsage)
	Cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, forceUsage)
	Cmd.Flags().IntVarP(&maxKeysFlag, "keys", "k", 1000, forceUsage)
}