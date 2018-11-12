// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mb

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/pkg/s3utils"
	"github.com/spf13/cobra"
	"storj/ditto/cmd/utils"
	"storj/ditto/pkg/config"
)

var mirroring = utils.GetObjectLayer

var Cmd = &cobra.Command{
	Use:     "make-bucket [mb].",
	Aliases: []string{"mb"},

	Args:  validateArgs,
	Short: "Creates bucket with desired name.",
	Long:  "Creates bucket with desired name.",
	RunE:  exec,
}

func exec(cmd *cobra.Command, args []string) error {

	// TODO: add location flag in init func
	objLayer, err := mirroring()
	if err != nil {
		return err
	}

	err = objLayer.MakeBucketWithLocation(context.Background(), args[0], "")

	if err != nil {
		return err
	}

	fmt.Printf("Bucket %s created\n", args[0])

	return nil
}

func validateArgs(cmd *cobra.Command, args []string) error {
	switch len(args) {
		case 0:
			return errors.New("at least one argument required")

		case 1:
			err := s3utils.CheckValidBucketName(args[0])

			if err != nil {
				return errors.New(fmt.Sprintf("%s. \n Please refer to this bucket naming conventions "+
					"(https://docs.aws.amazon.com/AmazonS3/latest/dev/BucketRestrictions.html", err.Error()))
			}

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
}
