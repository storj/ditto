package uploader

import (
	"testing"
	tutils "storj.io/ditto/pkg/utils/testing_utils"
	"github.com/minio/minio/pkg/hash"
	minio "github.com/minio/minio/cmd"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
)

func TestObjectUploader(t *testing.T) {
	cases := []struct {
		testName string
		testFunc func(t *testing.T)
	}{
		{
			"Main",
			func(t *testing.T) {
				bucket := "bucket"
				object := "object"

				testError := errors.New("test error")

				mirr := tutils.NewProxyObjectLayer()
				mirr.PutObjectFunc = func(ctx context.Context, bucket, object string, data *hash.Reader, metadata map[string]string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
					return minio.ObjectInfo{Bucket:bucket, Name:object}, testError
				}

				objUploader := ObjectUploader{mirr}
				resch := objUploader.UploadObjectAsync(nil, bucket, object, nil)
				res := <-resch
				assert.Error(t, res.Err)
				assert.Equal(t, testError, res.Err)
				assert.Equal(t, bucket, res.Oi.Bucket)
				assert.Equal(t, object, res.Oi.Name)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}
