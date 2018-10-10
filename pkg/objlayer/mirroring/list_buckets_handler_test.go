// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"storj.io/ditto/pkg/config"
	"storj.io/ditto/pkg/utils"
	"testing"

	minio "github.com/minio/minio/cmd"
	test "storj.io/ditto/pkg/utils/testing_utils"
)

func TestListBucketsHandler(t *testing.T) {

	prime := test.NewProxyObjectLayer()
	alter := test.NewProxyObjectLayer()

	logger := &test.MockLogger{}

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
					testName: "withMerge: prime error, alter success",

					testFunc: func() {

						m.Config.ListOptions.Merge = true

						prime.ListBucketsFunc = func(ctx context.Context) (buckets []minio.BucketInfo, err error) {
							return nil, errors.New("prime error")
						}

						alterBuckets := []minio.BucketInfo{
							{
								Name: "bucket1",
							},
							{
								Name: "bucket2",
							},
						}

						alter.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return alterBuckets, nil
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						logErr, _ := logger.GetLastLogEParam()

						assert.Error(t, logErr)
						assert.Nil(t, processError)
						assert.NotNil(t, buckets)
						assert.Equal(t, len(buckets), len(alterBuckets))
						assert.Equal(t, buckets, alterBuckets)
					},
				},
				{
					testName: "withMerge: alter error, prime success",

					testFunc: func() {

						m.Config.ListOptions.Merge = true

						alter.ListBucketsFunc = func(ctx context.Context) (buckets []minio.BucketInfo, err error) {
							return nil, errors.New("prime error")
						}

						primeBuckets := []minio.BucketInfo{
							{
								Name: "bucket1",
							},
							{
								Name: "bucket2",
							},
						}

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return primeBuckets, nil
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						logErr, _ := logger.GetLastLogEParam()

						assert.Error(t, logErr)
						assert.Nil(t, processError)
						assert.NotNil(t, buckets)
						assert.Equal(t, len(buckets), len(primeBuckets))
						assert.Equal(t, len(buckets), len(primeBuckets))
					},
				},
				{
					testName: "withMerge: both errors",

					testFunc: func() {
						primeError := errors.New("prime error")
						alterError := errors.New("alter error")

						m.Config.ListOptions.Merge = true

						alter.ListBucketsFunc = func(ctx context.Context) (buckets []minio.BucketInfo, err error) {
							return nil, alterError
						}

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return nil, primeError
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						assert.NotNil(t, processError)
						assert.Equal(t, processError.Error(), utils.CombineErrors([]error{ alterError, primeError }).Error())
						assert.Nil(t, buckets)
					},
				},
				{
					testName: "withMerge: both success",

					testFunc: func() {

						primeBuckets := []minio.BucketInfo {
							{ Name: "pb1" },
							{ Name: "pb2" },
						}

						alterBuckets := []minio.BucketInfo {
							{ Name: "ab1" },
							{ Name: "ab2" },
						}

						combinedBuckets := utils.CombineBucketsDistinct(primeBuckets, alterBuckets)

						m.Config.ListOptions.Merge = true

						alter.ListBucketsFunc = func(ctx context.Context) (buckets []minio.BucketInfo, err error) {
							return alterBuckets, nil
						}

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return primeBuckets, nil
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						assert.Nil(t, processError)
						assert.NoError(t, processError)
						assert.NotNil(t, buckets)
						assert.Equal(t, len(buckets), len(combinedBuckets))

						for i := 0; i < len(buckets); i++ {
							assert.Equal(t, buckets[i].Name, combinedBuckets[i].Name)
						}

					},
				},
				{
					testName: "withoutMerge: prime success, no retry",

					testFunc: func() {

						m.Config.ListOptions.DefaultOptions.ThrowImmediately = true

						primeBuckets := []minio.BucketInfo {
							{ Name: "pb1" },
							{ Name: "pb2" },
						}

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return primeBuckets, nil
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						assert.Nil(t, processError)
						assert.NoError(t, processError)
						assert.NotNil(t, buckets)
						assert.Equal(t, len(buckets), len(primeBuckets))

						for i := 0; i < len(buckets); i++ {
							assert.Equal(t, buckets[i].Name, primeBuckets[i].Name)
						}

					},
				},
				{
					testName: "withoutMerge: prime error, no retry",

					testFunc: func() {

						m.Config.ListOptions.DefaultOptions.ThrowImmediately = true

						primeError := errors.New("primeError")

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return nil, primeError
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						assert.NotNil(t, processError)
						assert.Error(t, processError)
						assert.Equal(t, processError.Error(), primeError.Error())
						assert.Nil(t, buckets)
					},
				},
				{
					testName: "withoutMerge: prime success, retry",

					testFunc: func() {

						primeBuckets := []minio.BucketInfo {
							{ Name: "pb1" },
							{ Name: "pb2" },
						}

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return primeBuckets, nil
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						assert.Nil(t, processError)
						assert.NoError(t, processError)
						assert.NotNil(t, buckets)
						assert.Nil(t, h.alterBuckets)
						assert.Equal(t, len(buckets), len(primeBuckets))

						for i := 0; i < len(buckets); i++ {
							assert.Equal(t, buckets[i].Name, primeBuckets[i].Name)
						}
					},
				},
				{
					testName: "withoutMerge: prime error, alter success, retry",

					testFunc: func() {

						primeError := errors.New("prime error")

						alterBuckets := []minio.BucketInfo {
							{ Name: "pb1" },
							{ Name: "pb2" },
						}

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return nil, primeError
						}

						alter.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return alterBuckets, nil
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						loggerError, _ := logger.GetLastLogEParam()

						assert.Nil(t, processError)
						assert.NoError(t, processError)
						assert.NotNil(t, buckets)
						assert.Nil(t, h.primeBuckets)
						assert.NotNil(t, h.alterBuckets)
						assert.Equal(t, len(buckets), len(alterBuckets))
						assert.Error(t, loggerError)
						assert.NotNil(t, loggerError)
						assert.Equal(t, loggerError.Error(), primeError.Error())

						for i := 0; i < len(buckets); i++ {
							assert.Equal(t, buckets[i].Name, alterBuckets[i].Name)
						}
					},
				},
				{
					testName: "withoutMerge: prime error, alter error, retry",

					testFunc: func() {

						primeError := errors.New("prime error")
						alterError := errors.New("alter error")

						combinedError := utils.CombineErrors([]error{ alterError, primeError })

						prime.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return nil, primeError
						}

						alter.ListBucketsFunc = func(ctx context.Context) ([]minio.BucketInfo, error) {
							return nil, alterError
						}

						h := NewListBucketsHandler(&m, ctx)

						buckets, processError := h.Process()

						loggerError, _ := logger.GetLastLogEParam()

						assert.NotNil(t, processError)
						assert.Error(t, processError)
						assert.Equal(t, processError.Error(), combinedError.Error())
						assert.Nil(t, buckets)
						assert.Nil(t, h.primeBuckets)
						assert.Nil(t, h.alterBuckets)
						assert.Error(t, loggerError)
						assert.NotNil(t, loggerError)
						assert.Equal(t, loggerError.Error(), primeError.Error())
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