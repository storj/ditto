package uploader

import (
	"storj.io/ditto/pkg/context"
	minio "github.com/minio/minio/cmd"
)

type MockFolderUploader struct {

}

func (f *MockFolderUploader) UploadFileAsync(ctx context.PutContext, bucket, lpath string) <-chan error {
	return emptyErrorChannel()
}

func (f *MockFolderUploader) UploadFolderAsync(ctx context.PutContext, bucket, lpath string) <-chan error {
	return emptyErrorChannel()
}

func (f *MockFolderUploader) SetObjLayer(layer minio.ObjectLayer) {

}

func emptyErrorChannel() <-chan error {
	errc := make(chan error, 1)
	errc <- nil
	return errc
}
