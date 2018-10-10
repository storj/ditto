package uploader

import (
	"storj.io/ditto/pkg/context"
	minio "github.com/minio/minio/cmd"
	"time"
)

type MockFolderUploader struct {
 	Upltime int
}

func (f *MockFolderUploader) UploadFileAsync(ctx context.PutContext, bucket, lpath string) <-chan error {
	return emptyErrorChannel(f.Upltime)
}

func (f *MockFolderUploader) UploadFolderAsync(ctx context.PutContext, bucket, lpath string) <-chan error {
	return emptyErrorChannel(f.Upltime)
}

func (f *MockFolderUploader) SetObjLayer(layer minio.ObjectLayer) {

}

func emptyErrorChannel(dur int) <-chan error {
	errc := make(chan error, 1)

	go func() {
		time.Sleep(time.Second * time.Duration(dur))
		errc <- nil
	}()

	return errc
}
