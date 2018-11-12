package uploader

import (
	"fmt"
	"storj/ditto/cmd/utils"
	"context"
	"errors"
	minio "github.com/minio/minio/cmd"
	dcontext "storj/ditto/pkg/context"
	"storj/ditto/pkg/filesys"
)

func checkObj(ctx context.Context, ol minio.ObjectLayer, bucket, object string) error {
	//TODO: check opts
	_, err := ol.GetObjectInfo(ctx, bucket, object, minio.ObjectOptions{})
	if err == nil {
		return errors.New(fmt.Sprintf("Object allready exists %s", object))
	}

	return nil
}

type fileUploader struct {
	ObjectUploader
	filesys.HashFileReader
}

func (u *fileUploader) UploadFileAsync(ctx dcontext.PutContext, bucket, lpath string) <-chan UploadResult {
	dresc := make(chan UploadResult, 1) // delayed result chanel for error handling
	res := UploadResult{}

	hfreader, err := u.ReadFileH(lpath)
	if err != nil {
		res.Err = err
		dresc <- res
		return dresc
	}

	object := utils.GetObjectName(hfreader.FileInfo().Name(), ctx.Prefix(), ctx.Delimiter())

	if !ctx.Force() {
		err = checkObj(ctx, u.ol, bucket, object)
		if err != nil {
			res.Err = err
			dresc <- res
			return dresc
		}
	}

	return u.UploadObjectAsync(ctx, bucket, object, hfreader.HashReader())
}
