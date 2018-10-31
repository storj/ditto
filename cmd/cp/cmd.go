// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package cp

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/pkg/s3utils"
	minio "github.com/minio/minio/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdUtils "storj.io/ditto/cmd/utils"
	"storj.io/ditto/pkg/config"
	"storj.io/ditto/pkg/utils"
)

var mirroring = cmdUtils.GetObjectLayer
var missingArgsErrorMessage = "at least three arguments required."

var Cmd = &cobra.Command{
	Use:     "copy [cp] srcBucket, srcObj, dstBucket, dstObj(OPTIONAL).",
	Aliases: []string{"cp"},
	Args:    validateArgs,
	Short:   "Creates a copy of an object.",
	Long: "Creates a copy of an object that is already stored in a bucket. dstObj is optional. " +
		"If not specified, dstObj name will be as srcObj.",
	RunE: exec,
}

var (
	defaultSourceKey    = "default_source"
	throwImmediatelyKey = "throw_immediately"
)

var (
	defaultSourceFlag    string
	throwImmediatelyFlag bool
)

func exec(cmd *cobra.Command, args []string) error {

	ctx := context.Background()
	objectLayer, err := mirroring()

	if err != nil {
		return err
	}

	objectInfo, err := objectLayer.GetObjectInfo(ctx, args[0], args[1], minio.ObjectOptions{})
	if err != nil {
		return err
	}

	_, err = objectLayer.GetBucketInfo(ctx, args[2])
	if err != nil {
		return err
	}

	dstObj := args[1]

	if len(args) == 4 {
		dstObj = args[3]
	}

	// TODO: enable object options in future
	_, err = objectLayer.CopyObject(ctx, args[0], args[1], args[2], dstObj, objectInfo, minio.ObjectOptions{}, minio.ObjectOptions{})

	if err != nil {
		return err
	}

	fmt.Printf("Object %s/%s copied\n", args[0], args[1])

	return nil
}

func validateArgs(cmd *cobra.Command, args []string) error {

	switch len(args) {
		case 0, 1, 2:
			return errors.New(missingArgsErrorMessage)
		case 3, 4:
			srcBucketNameErr := s3utils.CheckValidBucketName(args[0])
			dstBucketNameErr := s3utils.CheckValidBucketName(args[2])
			srcObjectNameErr := s3utils.CheckValidObjectName(args[1])

			var dstObjectNameErr error = nil

			if len(args) == 4 {
				dstObjectNameErr = s3utils.CheckValidObjectName(args[3])
			}

			err := utils.CombineErrors([]error{
				utils.NewError(srcBucketNameErr, "srcBucket - "),
				utils.NewError(dstBucketNameErr, "dstBucket - "),
				utils.NewError(srcObjectNameErr, "srcObject - "),
				utils.NewError(dstObjectNameErr, "dstObject - "),
			})

			if err != nil {
				return err
			}

			return nil

		default:
			return errors.New("too many arguments")
	}

	return nil
}

func init() {
	err := config.ReadConfig(true)
	if err != nil {
		println("error while reading config file: ", err)
	}

	Cmd.Flags().StringVarP(&defaultSourceFlag, defaultSourceKey, "s", "server2", "Defines source server to display list")
	Cmd.Flags().BoolVarP(&throwImmediatelyFlag, throwImmediatelyKey, "t", false, "In case of error, throw error immediately, or retry from other server")

	viper.BindPFlag(config.COPY_DEFAULT_SOURCE, Cmd.Flags().Lookup(defaultSourceKey))
	viper.BindPFlag(config.COPY_THROW_IMMEDIATELY, Cmd.Flags().Lookup(throwImmediatelyKey))
}
