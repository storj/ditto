// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mb

import (
	"context"
	"github.com/minio/minio/cmd"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"

	test "storj.io/ditto/pkg/utils/testing_utils"
)

func TestValidateArgs(t *testing.T) {

	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "Nil slice arguments",

			testFunc: func() {
				err := validateArgs(nil, nil)

				assert.Error(t, err)
				assert.Equal(t, "at least one argument required", err.Error())
			},
		},
		{
			testName: "Empty slice arguments",

			testFunc: func() {
				err := validateArgs(nil, []string{})

				assert.Error(t, err)
				assert.Equal(t, "at least one argument required", err.Error())
			},
		},
		{
			testName: "Invalid bucket name",

			testFunc: func() {
				err := validateArgs(nil, []string{"."})

				assert.Error(t, err)
			},
		},
		{
			testName: "Too many arguments",

			testFunc: func() {
				err := validateArgs(nil, []string{".", "."})

				assert.Error(t, err)
				assert.Equal(t, "too many arguments", err.Error())
			},
		},
		{
			testName: "Valid",

			testFunc: func() {
				err := validateArgs(nil, []string{"a.y.e.bucket"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}

func TestExec(t *testing.T) {

	prime := test.NewProxyObjectLayer

	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "MakeBucketWithLocation error",

			testFunc: func() {

				makeBucketErrorString := "MakeBucketWithLocation failed"
				mirroring = func() (cmd.ObjectLayer, error) {
					proxyObj := prime()
					proxyObj.MakeBucketWithLocationFunc = func(ctx context.Context, bucket string, location string) (err error) {
						return errors.New(makeBucketErrorString)
					}
					return proxyObj, nil
				}

				err := exec(nil, []string{"-flag"})

				assert.Error(t, err)
				assert.Equal(t, err.Error(), makeBucketErrorString)
			},
		},
		{
			testName: "MakeBucketWithLocation success",

			testFunc: func() {

				mirroring = func() (cmd.ObjectLayer, error) {
					proxyObj := prime()
					proxyObj.MakeBucketWithLocationFunc = func(ctx context.Context, bucket string, location string) (err error) {
						return nil
					}

					return proxyObj, nil
				}

				err := exec(nil, []string{"-flag"})

				assert.Nil(t, err)
				assert.NoError(t, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}

}
