package put

import (
	"storj.io/ditto/pkg/logger"
	"testing"
	tutils "storj.io/ditto/pkg/utils/testing_utils"
	minio "github.com/minio/minio/cmd"
	"context"
		"errors"
	"github.com/stretchr/testify/assert"
	"storj.io/ditto/pkg/uploader"
	fsystem "storj.io/ditto/pkg/filesys"
	"os"
	"fmt"
)

func TestExec(t *testing.T) {
	//testError := errors.New("test error")
	fileNotFoundError := errors.New("file not found")
	getGwError := errors.New("error retrieving obj layer")
	getObjLayerError := errors.New("error retrieving obj layer")
	bucketNotFoundError := errors.New("bucket not found")

	getBucketInfoError := func(ctx context.Context, bucket string) (minio.BucketInfo, error) {
		return minio.BucketInfo{}, bucketNotFoundError
	}

	cases := []struct {
		testName string
		testFunc func(t *testing.T)
	}{
		{
			"Error retrievieng gateway",
			func(t *testing.T) {
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return nil, getGwError
				}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, getGwError, err)
			},
		},
		{
			"Error retrivieng object layer",
			func(t *testing.T) {
				gw := &tutils.MockGateway{nil, getObjLayerError}
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return gw, nil
				}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, getObjLayerError, err)
			},
		},
		{
			"Bucket not found error",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()
				mirr.GetBucketInfoFunc = getBucketInfoError

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return gw, nil
				}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, bucketNotFoundError, err)
			},
		},
		{
			"No error, folder",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return gw, nil
				}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return false, nil })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.NoError(t, err)
			},
		},
		{
			"No error, file",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return gw, nil
				}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.NoError(t, err)
			},
		},
		{
			"File not found",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return gw, nil
				}

				uploader := &uploader.MockFolderUploader{}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return false, fileNotFoundError })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.Error(t, err)
				assert.Equal(t, fileNotFoundError, err)
			},
		},
		//TODO: investigate worng duration
		{
			"Interrupt",
			func(t *testing.T) {
				mirr := tutils.NewProxyObjectLayer()

				gw := &tutils.MockGateway{mirr, nil}
				lg := &tutils.MockLogger{}

				resolver := func(logger.Logger) (minio.Gateway, error) {
					return gw, nil
				}

				sigc <- os.Interrupt

				uploader := &uploader.MockFolderUploader{2}
				dchecker := fsystem.MockDirChecker(func(string) (bool, error) { return true, nil })

				exec := newPutExec(resolver, uploader, dchecker, lg)
				err := exec.runE(nil, []string{"bucket", "localpath"})
				assert.NoError(t, err)

				intrplog, err := lg.GetLastLogParam()
				assert.NoError(t, err)
				assert.Equal(t, 1, lg.LogCount())
				assert.Equal(t, 0, lg.LogECount())
				assert.Equal(t, fmt.Sprintf("Catched interrupt! %s\n", os.Interrupt), intrplog)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}
