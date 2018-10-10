package put

import (
	"testing"
	tutils "storj.io/ditto/pkg/utils/testing_utils"
	minio "github.com/minio/minio/cmd"
	"context"
	"github.com/minio/minio/pkg/hash"
	"errors"
	"github.com/stretchr/testify/assert"
	"storj.io/ditto/pkg/uploader"
	fsystem "storj.io/ditto/pkg/filesys"
)

func TestExec(t *testing.T) {
	testError := errors.New("test error")
	fileNotFoundError := errors.New("file not found")
	getObjLayerError := errors.New("error retrieving obj layer")
	bucketNotFoundError := errors.New("bucket not found")

	putError := func(ctx context.Context, bucket, object string, data *hash.Reader, metadata map[string]string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
		return minio.ObjectInfo{}, testError
	}

	getBucketInfoError := func(ctx context.Context, bucket string) (minio.BucketInfo, error) {
		return minio.BucketInfo{}, bucketNotFoundError
	}

	cases := []struct {
		testName string
		testFunc func(t *testing.T)
	}{
		{
			"Error retrivieng object layer",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()
				mirr.PutObjectFunc = putError

				gw := &tutils.MockGateway{nil, getObjLayerError}
				lg := &tutils.MockLogger{}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(gw, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, getObjLayerError, err)
			},
		},
		{
			"Bucket not found error",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()
				mirr.PutObjectFunc = putError
				mirr.GetBucketInfoFunc = getBucketInfoError

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(gw, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, bucketNotFoundError, err)
			},
		},
		{
			"No error, folder",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()
				mirr.PutObjectFunc = putError

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return false, nil })

				exec := newPutExec(gw, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.NoError(t, err)
			},
		},
		{
			"No error, file",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()
				mirr.PutObjectFunc = putError

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(gw, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.NoError(t, err)
			},
		},
		{
			"File not found",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()
				mirr.PutObjectFunc = putError

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return false, fileNotFoundError })

				exec := newPutExec(gw, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, fileNotFoundError, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}
