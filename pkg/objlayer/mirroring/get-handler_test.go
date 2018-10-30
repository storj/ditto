package mirroring

import (
	"testing"
	tutils "storj.io/ditto/pkg/utils/testing_utils"
	"errors"
	minio "github.com/minio/minio/cmd"
	"context"
	"io"
	"bytes"
	"github.com/stretchr/testify/assert"
)

type getFunc func(ctx context.Context, bucket, object string, startOffset int64, length int64, writer io.Writer, etag string, opts minio.ObjectOptions) (err error)

func getObjectFuncFact(callback func(writer io.Writer, offset, length int64) error) getFunc {
	return func(ctx context.Context, bucket, object string, startOffset int64, length int64, writer io.Writer, etag string, opts minio.ObjectOptions) (err error) {

		return callback(writer, startOffset, length)
	}
}

func TestGetHandler(t *testing.T) {
	prime := tutils.NewProxyObjectLayer()
	alter := tutils.NewProxyObjectLayer()

	testError := errors.New("test error")
	testError2 := errors.New("test error 2")

	m := MirroringObjectLayer{
		Prime: prime,
		Alter: alter,
		Logger: nil,
	}

	getNoError := getObjectFuncFact(func(writer io.Writer, offset, length int64) error {
		return nil
	})

	opts := minio.ObjectOptions{}

	cases := []struct {
		testName string
		testFunc func(*testing.T)
	}{
		{
			"No error",
			func(t *testing.T) {
				prime.GetObjectFunc = getNoError
				alter.GetObjectFunc = getNoError

				ctx := context.Background()

				err := m.GetObject(ctx, "bucket", "object", 0, 0, nil, "etag", opts)
				assert.NoError(t, err)
			},
		},
		{
			"Error in both",
			func(t *testing.T) {
				prime.GetObjectFunc = getObjectFuncFact(func(writer io.Writer, offset, length int64) error {
					return testError
				})

				alter.GetObjectFunc = getObjectFuncFact(func(writer io.Writer, offset, length int64) error {
					return testError2
				})

				ctx := context.Background()

				err := m.GetObject(ctx, "bucket", "object", 0, 0, nil, "etag", opts)
				assert.Error(t, err)
				assert.Equal(t, testError2, err)
			},
		},
		{
			"Error while reading from main",
			func(t *testing.T) {
				obj1 := []byte("abc45678901234567890")
				obj2 := []byte("09876543210987654abc")

				prime.GetObjectFunc = getObjectFuncFact(func(writer io.Writer, offset, length int64) error {
					reader := bytes.NewReader(obj1)
					io.CopyN(writer, reader, length-10)
					return testError
				})

				alter.GetObjectFunc = getObjectFuncFact(func(writer io.Writer, offset, length int64) error {
					reader := bytes.NewReader(obj2[offset:])
					io.CopyN(writer, reader, length)
					return nil
				})

				ctx := context.Background()
				data := bytes.NewBuffer(nil)

				err := m.GetObject(ctx, "bucket", "object", 0, int64(len(obj1)), data, "etag", opts)
				assert.NoError(t, err)
				assert.Equal(t, append(obj1[:10], obj2[10:]...), data.Bytes())
			},
		},
	}

	for _, c := range cases {
		t.Run(c.testName, c.testFunc)
	}
}
