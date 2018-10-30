// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package delete

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateArgs(t *testing.T) {

	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "FAILS command - 'delete' ",

			testFunc: func() {
				err := validateArgs(nil, nil)

				assert.Error(t, err)
				assert.Equal(t, noArgsMessage, err.Error())
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName' ",

			testFunc: func() {
				err := validateArgs(nil, []string{"bucketName"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},


		{
			testName: "FAILS command - 'delete bucketName -p' ",

			testFunc: func() {

				prefixFlag = true

				err := validateArgs(nil, []string{"bucketName"})

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, prefixMissingMessage, err.Error())
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName prefix -p' ",

			testFunc: func() {

				prefixFlag = true

				err := validateArgs(nil, []string{"bucketName", "prefix"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName prefix1 prefix2 prefix3 prefix4 -p' ",

			testFunc: func() {

				prefixFlag = true

				err := validateArgs(nil, []string{"bucketName", "prefix1", "prefix2", "prefix3", "prefix4", "prefix5"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},


		{
			testName: "FAILS command - 'delete bucketName -pr' ",

			testFunc: func() {

				prefixFlag = true
				recursiveFlag = true

				err := validateArgs(nil, []string{"bucketName"})

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, prefixMissingMessage, err.Error())
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName prefix -pr' ",

			testFunc: func() {

				prefixFlag = true
				recursiveFlag = true

				err := validateArgs(nil, []string{"bucketName", "prefix"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName prefix1 prefix2 prefix3 prefix4 -pr' ",

			testFunc: func() {

				prefixFlag = true
				recursiveFlag = true

				err := validateArgs(nil, []string{"bucketName", "prefix1", "prefix2", "prefix3", "prefix4", "prefix5"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},


		{
			testName: "FAILS command - 'delete -rf' (nil slice)",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true

				err := validateArgs(nil, nil)

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, noArgsMessage, err.Error())
			},
		},
		{
			testName: "FAILS command - 'delete -rf' (empty slice)",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true

				err := validateArgs(nil, []string{})

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, noArgsMessage, err.Error())
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName -rf' ",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true

				err := validateArgs(nil, []string{"bucketName"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName some text to test args -rf' ",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true

				err := validateArgs(nil, []string{"bucketName", "some", "text", "to", "test", "args"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},


		{
			testName: "FAILS command - 'delete -prf' ",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true
				prefixFlag = true

				err := validateArgs(nil, []string{})

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, noArgsMessage, err.Error())
			},
		},
		{
			testName: "FAILS command - 'delete bucket -prf' ",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true
				prefixFlag = true

				err := validateArgs(nil, []string{"bucket"})

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, prefixMissingMessage, err.Error())
			},
		},
		{
			testName: "SUCCESS command - 'delete bucket . -prf' ",

			testFunc: func() {

				forceFlag = true
				recursiveFlag = true
				prefixFlag = true

				err := validateArgs(nil, []string{"bucket", "."})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},


		{
			testName: "SUCCESS command - 'delete bucketName fileName' ",

			testFunc: func() {
				err := validateArgs(nil, []string{"bucketName", "fileName"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
		{
			testName: "SUCCESS command - 'delete . fileName' ",

			testFunc: func() {
				err := validateArgs(nil, []string{".", "fileName"})

				assert.Error(t, err)
				assert.NotNil(t, err)
				assert.Equal(t, bucketNameInvalidMessage, err.Error())
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName .' ",

			testFunc: func() {
				err := validateArgs(nil, []string{"bucketName", "."})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
		{
			testName: "SUCCESS command - 'delete bucketName fileName1 fileName2 fileName3' ",

			testFunc: func() {
				err := validateArgs(nil, []string{"bucketName", "fileName"})

				assert.NoError(t, err)
				assert.Nil(t, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()

			forceFlag     = false
			recursiveFlag = false
			prefixFlag    = false
		})
	}
}