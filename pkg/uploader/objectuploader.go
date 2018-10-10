package uploader

import (
	"github.com/minio/minio/pkg/hash"
	"context"
	minio "github.com/minio/minio/cmd"
)

type UploadResult struct {
	Oi minio.ObjectInfo
	Err error
}

type ObjectUploader struct {
	ol minio.ObjectLayer
}

func (u *ObjectUploader) UploadObjectAsync(ctx context.Context, bucket, object string, data *hash.Reader) <-chan UploadResult {
	resc := make(chan UploadResult)

	utask := func(resc chan<- UploadResult) {
		oi, err := u.ol.PutObject(ctx, bucket, object, data, make(map[string]string), minio.ObjectOptions{})
		
		res := UploadResult{}
		res.Oi = oi
		res.Err = err
		resc<-res
	}

	go utask(resc)
	return resc
}