// Copyright (C) 2018 Storj Labs, Inc.l "storj.io/ditto/pkg/logger"
// See LICENSE for copying information.

package get

import (
	"errors"
	"fmt"
	"github.com/minio/minio/pkg/auth"
	"github.com/spf13/cobra"
	"storj.io/ditto/cmd/get/downloader"
	"storj.io/ditto/global"
	gw "storj.io/ditto/pkg/gateway"
	l "storj.io/ditto/pkg/logger"
)

// getCmd represents the get command
var Cmd = &cobra.Command{
	Use:   "get [bucket name] [object name](opt) [OPTIONS]",
	Args: validateArgs,
	Short: "Download files and buckets",
	Long: ``,
	RunE: newGetExec(&gw.Mirroring{nil, &l.StdOutLogger}, &l.StdOutLogger).runE,
}

func runE(cmd *cobra.Command, args []string) error {
	fmt.Println("get called")

	//for i := range args {
	//	fmt.Printf("Arg%d: %s\n", i, args[i])
	//}

	fmt.Printf("Filename: %s\n", nameFlag)
	fmt.Printf("Cwd: %s\n", global.Params.GetCwd())

	cwd := global.Params.GetCwd()
	if cwd == "" {
		return errors.New("unable to get current working directory")
	}

	var mirrGateway gw.Mirroring = gw.Mirroring{Logger: &l.StdOutLogger}

	var mirr, err =  mirrGateway.NewGatewayLayer(auth.Credentials{})
	if err != nil {
		return errors.New("unable to start mirroring service")
	}

	params := downloader.NewDefaultParams()
	params.SetPath(nameFlag)
	params.SetPrefix(prefixFlag)
	params.SetRecursive(true)

	dwn := downloader.NewDownloader(mirr, cwd, params)

	var bucketName string
	var objectName string
	bucketName = args[0]

	if len(args) == maxArg {
		objectName = args[1]
	}

	if objectName == "" {
		err = dwn.GetBucket(bucketName)
	} else {
		err = dwn.GetObject(bucketName, objectName)
	}

	if err != nil {
		fmt.Printf("Download error: %s\n", err)
		return nil
	}

	fmt.Println("The end")
	return nil
}

var (
	minArg = 1
	maxArg = 2

	nameFlag string
	nameUsage = "Path of the file or folder to be downloaded. A raw filename can be used to download to current directory under that name.\n" +
		"If no objectname provided folder under that name will be created"

	prefixFlag string
	prefixUsage = ""

	recursiveFlag bool
	recursiveUsage = ""
)

func init() {
	Cmd.Flags().StringVarP(&nameFlag, "name", "n", "", nameUsage)
	Cmd.Flags().StringVarP(&prefixFlag, "prefix", "p", "", prefixUsage)
	Cmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, recursiveUsage)
}