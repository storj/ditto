// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package list

import (
	"context"
	"errors"
	"fmt"
	minio "github.com/minio/minio/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
	"storj/ditto/cmd/utils"
	"storj/ditto/pkg/config"
)

// Function listed as var for testing purposes only
var mirroring = utils.GetObjectLayer

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Displays bucket list and all files at specified bucket",
	Long:  `Displays bucket list and all files at specified bucket`,
	Aliases: []string{"ls"},
	Args:  validateArgs,
	RunE:  exec,
}

// Error messages
var (
	prefixMissingMessage = "two arguments expected. ditto list -p bucketName prefix"
)

// Flag keys
var (
	prefixKey    = "prefix"
	recursiveKey = "recursive"
	delimiterKey = "delimiter"

	mergeKey            = "merge"
	defaultSourceKey    = "default_source"
	throwImmediatelyKey = "throw_immediately"
)

// Flags variables
var (
	prefixFlag    bool
	recursiveFlag bool
	delimiterFlag string

	mergeFlag            bool
	defaultSourceFlag    string
	throwImmediatelyFlag bool
)

func exec(cmd *cobra.Command, args []string) error {
	objLayer, err := mirroring()
	if err != nil {
		return err
	}

	argsLen := len(args)
	if prefixFlag {
		if argsLen != 2 {
			return errors.New(prefixMissingMessage)
		}

		if recursiveFlag {
			return listObjectsRecursive(context.Background(), objLayer, args[0], args[1], delimiterFlag)
		}

		return listObjects(objLayer, args[0], args[1])
	}

	switch len(args) {
	case 0:
		return listBuckets(objLayer)
	case 1:
		return listObjects(objLayer, args[0], "")
	default:
		return nil
	}
}

func validateArgs(cmd *cobra.Command, args []string) error {
	utils.LoadFlagValueFromViperIfNotSet(cmd, mergeKey, config.LIST_MERGE)

	argsLen := len(args)

	if prefixFlag {
		if argsLen != 2 {
			return errors.New(prefixMissingMessage)
		}

		return nil
	}

	if argsLen > 1 {
		return errors.New("too many arguments")
	}

	return nil
}

func listBuckets(layer minio.ObjectLayer) error {
	buckets, err := layer.ListBuckets(context.Background())

	if err != nil {
		return errors.New(fmt.Sprintf("Error while requesting bucket list: %s\n", err.Error()))
	}

	if len(buckets) == 0 {
		fmt.Println("No buckets found")
		return nil
	}

	if mergeFlag {
		fmt.Println("Merged bucket list:")
	} else {
		fmt.Println("Buckets:")
	}

	for _, bucket := range buckets {
		fmt.Println("\t", bucket.Name)
	}

	return nil
}

func listObjects(layer minio.ObjectLayer, bucketName string, prefix string) error {
	u, err := url.Parse(bucketName)
	if err != nil {

		return errors.New(fmt.Sprintf("Error while parsing bucketName: %s", err.Error()))
	}

	result, err := layer.ListObjectsV2(context.Background(), u.Path, prefix, "", delimiterFlag, 1000, false, "")

	if err != nil {

		return errors.New(fmt.Sprintf("Error while creating List Objects Layer for bucket %s,\nError: %s\n", bucketName, err.Error()))
	}

	printFilesList(result.Objects, bucketName)

	return nil
}

func listObjectsRecursive(ctx context.Context, layer minio.ObjectLayer, bucketName string, prefix string, delimiter string) error {
	prefixes := [] string{prefix}

	var fileList []minio.ObjectInfo
	for ; len(prefixes) > 0; {
		listObjInfo, err := layer.ListObjectsV2(ctx, bucketName, prefixes[0], "", delimiter, 1000, false, "")
		if err != nil {
			return err
		}
		fileList = append(fileList, listObjInfo.Objects...)

		if !recursiveFlag {
			break
		}

		// removing first prefix
		prefixes = append(prefixes[:0], prefixes[1:]...)

		// addition of new prefixed
		prefixes = append(prefixes, listObjInfo.Prefixes...)
	}

	printFilesList(fileList, bucketName)

	return nil
}

func printFilesList(list []minio.ObjectInfo, bucketName string) {
	if len(list) == 0 {
		fmt.Printf("bucket '%s' is empty", bucketName)
		return
	}

	if mergeFlag {
		fmt.Printf("Merged files list from %s\n", bucketName)
	} else {
		fmt.Printf("Files at %s\n", bucketName)
	}

	for _, object := range list {
		fmt.Println("\t", object.Name)
	}
}

func init() {
	err := config.ReadConfig(true)
	if err != nil {
		println("error while reading config file: ", err)
	}

	Cmd.Flags().BoolVarP(&prefixFlag, prefixKey, "p", false, "Folder simulation path")
	Cmd.Flags().BoolVarP(&recursiveFlag, recursiveKey, "r", false, "Shows all nested folder if prefix set")
	Cmd.Flags().StringVarP(&delimiterFlag, delimiterKey, "d", "/", "Char or char sequence that should be used as prefix delimiter")
	Cmd.Flags().BoolVarP(&mergeFlag, mergeKey, "m", false, "Display list from both servers and merge to single list")
	Cmd.Flags().StringVarP(&defaultSourceFlag, defaultSourceKey, "s", "server2", "Defines source server to display list")
	Cmd.Flags().BoolVarP(&throwImmediatelyFlag, throwImmediatelyKey, "t", false, "In case of error, throw error immediately, or retry from other server")

	viper.BindPFlag(config.LIST_MERGE, Cmd.Flags().Lookup(mergeKey))
	viper.BindPFlag(config.LIST_DEFAULT_SOURCE, Cmd.Flags().Lookup(defaultSourceKey))
	viper.BindPFlag(config.LIST_THROW_IMMEDIATELY, Cmd.Flags().Lookup(throwImmediatelyKey))

}
