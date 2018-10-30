// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package delete

import (
	"context"
	"errors"
	"github.com/minio/minio-go/pkg/s3utils"
	"github.com/spf13/cobra"

	minio "github.com/minio/minio/cmd"
	cmdUtils "storj.io/ditto/cmd/utils"
)

var mirroring = cmdUtils.GetObjectLayer

var Cmd = &cobra.Command{
	Use:   "delete ",
	Args:  validateArgs,
	Short: "Delete files and buckets.",
	Long: "to delete object - delete bucketName objectName.\n" +
		"to delete bucket - delete bucketName.\n",
	RunE: exec,
}

//Error messages ares
var (
	noArgsMessage 			 = "at least one argument required"
	prefixMissingMessage 	 = "prefix argument is missing"
	bucketNameInvalidMessage = "bucket name is not valid"
)

//Flags values area
var (
	forceFlag     bool
	recursiveFlag bool
	prefixFlag    bool
	delimeterFlag string
)

func exec(cmd *cobra.Command, args []string) (err error) {

	var bucketName, objectName, prefix, delimiter string
	bucketName = args[0]
	delimiter = "/"
	var ctx = context.Background()

	if prefixFlag {
		objectName, prefix = args[0], args[1]
	}

	m, err := mirroring()
	if err != nil {
		return err
	}

	rf := recursiveFlag && forceFlag
	noFlags := !(recursiveFlag && prefixFlag && forceFlag)

	switch {
		case prefixFlag:

			return deleteRecursive(ctx, m, bucketName, prefix, delimiter)

		case rf:
			err = deleteRecursive(ctx, m, bucketName, "", delimiter)

			if err != nil {
				return err
			}

			return m.DeleteBucket(ctx, bucketName)

		case noFlags:

			switch len(args) {
			case 1:
				err = m.DeleteBucket(ctx, bucketName)
			case 2:
				err = m.DeleteObject(ctx, bucketName, objectName)
			}
	}

	return
}

func validateArgs(cmd *cobra.Command, args []string) error {

	argsLength := len(args)

	if argsLength == 0 {
		return errors.New(noArgsMessage)
	}

	err := s3utils.CheckValidBucketName(args[0])

	if err != nil {
		return errors.New(bucketNameInvalidMessage)
	}

	noFlags := !(recursiveFlag && prefixFlag && forceFlag)

	switch {
		case prefixFlag:

			if argsLength < 2 {
				return errors.New(prefixMissingMessage)
			}

			err = s3utils.CheckValidObjectNamePrefix(args[1])

			if err != nil {
				return err
			}

		case noFlags:

			if argsLength == 1 {
				return nil
			}

			err = s3utils.CheckValidObjectName(args[1])

			if err != nil {
				return err
			}
	}

	return nil
}

func init() {
	Cmd.Flags().BoolVarP(&forceFlag, "force", "f", false,
		"if force flag applied - all files without prefixes in bucket will be removed.")
	Cmd.Flags().BoolVarP(&recursiveFlag, "recursive", "r", false, "User force flag to delete bucket")
	Cmd.Flags().BoolVarP(&prefixFlag, "prefix", "p", false, "User force flag to delete bucket")
	Cmd.Flags().StringVarP(&delimeterFlag, "delimiter", "d", "/", "delimiter usage")
}

func deleteRecursive(ctx context.Context, m minio.ObjectLayer, bucketName, prefix, delimiter string) error {
	prefixes  := []string{prefix}

	for ; len(prefixes) > 0 ;  {

		listObjInfo, err := m.ListObjectsV2(ctx, bucketName, prefixes[0], "", delimiter, 1000,false, "")

		if err != nil {
			return err
		}

		objCount := len(listObjInfo.Objects)

		// deleting files from current prefix
		for i := 0; i < objCount; i++ {
			err = m.DeleteObject(ctx, bucketName, listObjInfo.Objects[i].Name)

			if err != nil {
				return err
			}
		}

		// check for recursion
		if !recursiveFlag {
			break
		}

		// removing first prefix
		prefixes = append(prefixes[:0], prefixes[1:]...)

		// addition of new prefixed
		prefixes = append(prefixes, listObjInfo.Prefixes...)
	}

	return nil
}