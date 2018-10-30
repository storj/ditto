// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"

	minio "github.com/minio/minio/cmd"
	test "storj.io/ditto/pkg/utils/testing_utils"
)

func TestGetObjectInfoHandler(t *testing.T) {

	prime := test.NewProxyObjectLayer()
	alter := test.NewProxyObjectLayer()

	logger := &test.MockLogger{}

	m := MirroringObjectLayer{
		Prime: prime,
		Alter: alter,
		Logger: logger,
	}

	cases := []struct {
		testName, address string
		testFunc          func()
	}{
		{
			testName: "CopyObjectHandler: prime success, alter is not called",

			testFunc: func() {
				isAlterCalled := false

				h := NewGetObjectInfoHandler(&m, context.Background(), "bucket", "object", minio.ObjectOptions{})

				alter.GetObjectInfoFunc = func(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
					isAlterCalled = true

					return minio.ObjectInfo{},nil
				}

				h.Process()

				assert.Nil(t, h.primeErr)
				assert.Nil(t, h.alterErr)
				assert.Equal(t, false, isAlterCalled)
			},
		},
		{
			testName: "CopyObjectHandler: prime error, alter called",

			testFunc: func() {
				isAlterCalled := false

				h := NewGetObjectInfoHandler(&m, context.Background(), "bucket", "object", minio.ObjectOptions{})

				prime.GetObjectInfoFunc = func(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
					return minio.ObjectInfo{}, errors.New("prime failed")
				}

				alter.GetObjectInfoFunc = func(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
					isAlterCalled = true
					return minio.ObjectInfo{},nil
				}

				h.Process()

				logPrimeErr, _ := logger.GetLastLogEParam()

				assert.NotNil(t, h.primeErr)
				assert.NotNil(t, h.primeErr)
				assert.Nil(t, h.alterErr)
				assert.Nil(t, logPrimeErr)
				assert.NoError(t, logPrimeErr)
				//assert.Equal(t, "prime failed", logPrimeErr.Error())
				assert.Equal(t, true, isAlterCalled)
			},
		},
		{
			testName: "CopyObjectHandler: both error",

			testFunc: func() {
				isAlterCalled := false

				h := NewGetBucketInfoHandler(&m, context.Background(), "src_bucket")

				prime.GetBucketInfoFunc = func(ctx context.Context, bucket string) (objInfo minio.BucketInfo, err error) {
					return minio.BucketInfo{}, errors.New("prime failed")
				}

				alter.GetBucketInfoFunc = func(ctx context.Context, bucket string) (objInfo minio.BucketInfo, err error) {
					isAlterCalled = true
					return minio.BucketInfo{}, errors.New("alter failed")
				}

				h.Process()

				logErr, _ := logger.GetLastLogEParam()

				assert.NotNil(t, h.primeErr)
				assert.NotNil(t, h.alterErr)
				assert.NotNil(t, logErr)
				assert.Equal(t, "alter failed", logErr.Error())
				assert.Equal(t, true, isAlterCalled)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}

