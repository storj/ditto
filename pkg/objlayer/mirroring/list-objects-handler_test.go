// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"storj/ditto/pkg/config"
	"storj/ditto/pkg/utils"
	"testing"

	minio "github.com/minio/minio/cmd"
	test "storj/ditto/pkg/utils/testing_utils"
)

func TestListObjectsHandler(t *testing.T) {

	prime := test.NewProxyObjectLayer()
	alter := test.NewProxyObjectLayer()

	logger := &test.MockDiffLogger{}

	m := MirroringObjectLayer{
		Prime: prime,
		Alter: alter,
		Logger: logger,
		Config: &config.Config{
			ListOptions: &config.ListOptions{
				DefaultOptions: &config.DefaultOptions{},
			},
		},
	}

	ctx := context.Background()

	cases := []struct {
		testName, address string
		testFunc          func()
	}{
		{
			testName: "merge: prime error, alter success",

			testFunc: func() {

				m.Config.ListOptions.Merge = true

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, errors.New("prime error")
				}


				alterObjects := []minio.ObjectInfo {
					{
						Name: "bucket1",
					},
					{
						Name: "bucket2",
					},
				}

				alter.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {

					result = minio.ListObjectsInfo {
						Objects: alterObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objectInfo, processError := h.Process()

				logErr, _ := logger.GetLastLogEParam()

				assert.Error(t, logErr)
				assert.Nil(t, processError)
				assert.NotNil(t, objectInfo)
				assert.NotNil(t, objectInfo.Objects)
				assert.True(t, len(objectInfo.Objects) > 0)
				assert.Equal(t, len(objectInfo.Objects), len(alterObjects))
				assert.Equal(t, objectInfo.Objects, alterObjects)
			},
		},
		{
			testName: "merge: alter error, prime success",

			testFunc: func() {

				m.Config.ListOptions.Merge = true

				alter.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, errors.New("alter error")
				}

				primeObjects := []minio.ObjectInfo {
					{
						Name: "bucket1",
					},
					{
						Name: "bucket2",
					},
				}

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {

					result = minio.ListObjectsInfo {
						Objects: primeObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objectInfo, processError := h.Process()

				logErr, _ := logger.GetLastLogEParam()

				assert.Error(t, logErr)
				assert.Nil(t, processError)
				assert.NotNil(t, objectInfo)
				assert.NotNil(t, objectInfo.Objects)
				assert.Equal(t, len(objectInfo.Objects), len(primeObjects))
			},
		},
		{
			testName: "merge: both errors",

			testFunc: func() {
				primeError := errors.New("prime error")
				alterError := errors.New("alter error")

				m.Config.ListOptions.Merge = true

				alter.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, alterError
				}

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, primeError
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objectInfo, processError := h.Process()

				assert.NotNil(t, processError)
				assert.Error(t, processError)
				assert.True(t, len(objectInfo.Objects) == 0)
				assert.Equal(t, processError.Error(), utils.CombineErrors([]error{ alterError, primeError }).Error())
			},
		},
		{
			testName: "merge: both success",

			testFunc: func() {

				primeObjects := []minio.ObjectInfo {
					{ Name: "pb1" },
					{ Name: "pb2" },
				}

				alterObjects := []minio.ObjectInfo {
					{ Name: "ab1" },
					{ Name: "ab2" },
				}

				combinedBuckets := utils.CombineObjectsDistinct(primeObjects, alterObjects)
				expectedDiff    := utils.ListObjectsWithDifference(primeObjects, alterObjects)

				m.Config.ListOptions.Merge = true

				alter.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					result = minio.ListObjectsInfo {
						Objects: alterObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					result = minio.ListObjectsInfo {
						Objects: primeObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objInfo, processError := h.Process()
				loggerDiff := logger.GetDiff()

				assert.Nil(t, processError)
				assert.NoError(t, processError)
				assert.NotNil(t, objInfo.Objects)
				assert.Equal(t, len(objInfo.Objects), len(combinedBuckets))
				assert.Equal(t, len(loggerDiff), len(expectedDiff))

				for i := 0; i < len(objInfo.Objects); i++ {
					assert.Equal(t, objInfo.Objects[i].Name, combinedBuckets[i].Name)
				}
				for i := 0; i < len(loggerDiff); i++ {
					assert.Equal(t, expectedDiff[i].Name, loggerDiff[i].Name)
					assert.Equal(t, expectedDiff[i].Diff, loggerDiff[i].Diff)
				}

			},
		},
		{
			testName: "withoutMerge: prime success, no retry",

			testFunc: func() {

				m.Config.ListOptions.DefaultOptions.ThrowImmediately = true

				primeObjects := []minio.ObjectInfo {
					{ Name: "pb1" },
					{ Name: "pb2" },
				}

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					result = minio.ListObjectsInfo {
						Objects: primeObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objInfo, processError := h.Process()

				assert.Nil(t, processError)
				assert.NoError(t, processError)
				assert.NotNil(t, objInfo.Objects)
				assert.Equal(t, len(objInfo.Objects), len(primeObjects))

				for i := 0; i < len(objInfo.Objects); i++ {
					assert.Equal(t, objInfo.Objects[i].Name, primeObjects[i].Name)
				}

			},
		},
		{
			testName: "withoutMerge: prime error, no retry",

			testFunc: func() {

				m.Config.ListOptions.DefaultOptions.ThrowImmediately = true

				primeError := errors.New("primeError")

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, primeError
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objInfo, processError := h.Process()

				assert.NotNil(t, processError)
				assert.Error(t, processError)
				assert.Equal(t, processError.Error(), primeError.Error())
				assert.True(t, len(objInfo.Objects) == 0)
			},
		},
		{
			testName: "withoutMerge: prime success, retry",

			testFunc: func() {

				primeObjects := []minio.ObjectInfo {
					{ Name: "pb1" },
					{ Name: "pb2" },
				}

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					result = minio.ListObjectsInfo {
						Objects: primeObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objInfo, processError := h.Process()

				assert.Nil(t, processError)
				assert.NoError(t, processError)
				assert.Equal(t, len(objInfo.Objects), len(primeObjects))

				for i := 0; i < len(objInfo.Objects); i++ {
					assert.Equal(t, objInfo.Objects[i].Name, primeObjects[i].Name)
				}
			},
		},
		{
			testName: "withoutMerge: prime error, alter success, retry",

			testFunc: func() {

				primeError := errors.New("prime error")

				alterObjects := []minio.ObjectInfo {
					{ Name: "pb1" },
					{ Name: "pb2" },
				}

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, primeError
				}

				alter.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					result = minio.ListObjectsInfo {
						Objects: alterObjects,
						NextMarker: "next",
						IsTruncated: false,
						Prefixes: nil,
					}

					return
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objInfo, processError := h.Process()


				loggerError, _ := logger.GetLastLogEParam()

				assert.Nil(t, processError)
				assert.NoError(t, processError)
				assert.Equal(t, len(objInfo.Objects), len(alterObjects))
				assert.Error(t, loggerError)
				assert.NotNil(t, loggerError)
				assert.Equal(t, loggerError.Error(), primeError.Error())

				for i := 0; i < len(objInfo.Objects); i++ {
					assert.Equal(t, objInfo.Objects[i].Name, alterObjects[i].Name)
				}
			},
		},
		{
			testName: "withoutMerge: prime error, alter error, retry",

			testFunc: func() {

				primeError := errors.New("prime error")
				alterError := errors.New("alter error")

				combinedError := utils.CombineErrors([]error{ alterError, primeError })

				prime.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, primeError
				}

				alter.ListObjectsFunc = func(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
					return result, alterError
				}

				h := NewListObjectsHandler(&m, ctx, "", "", "", "", 1)

				objInfo, processError := h.Process()

				loggerError, _ := logger.GetLastLogEParam()

				assert.NotNil(t, processError)
				assert.Error(t, processError)
				assert.Equal(t, processError.Error(), combinedError.Error())
				assert.Error(t, loggerError)
				assert.NotNil(t, loggerError)
				assert.Equal(t, loggerError.Error(), primeError.Error())
				assert.True(t, len(objInfo.Objects) == 0)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()

			m.Config.ListOptions.Merge = false
			m.Config.ListOptions.DefaultOptions.ThrowImmediately = false
			logger.Clear()
		})
	}
}