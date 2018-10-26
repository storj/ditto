package uploader

import (
	"path"
	"storj.io/ditto/cmd/utils"
	l "storj.io/ditto/pkg/logger"
	"storj.io/ditto/pkg/context"
	minio "github.com/minio/minio/cmd"
	"storj.io/ditto/pkg/filesys"
	"fmt"
)

type AsyncUploader interface {
	UploadFileAsync(ctx context.PutContext, bucket, lpath string) <-chan error
	UploadFolderAsync(ctx context.PutContext, bucket, lpath string) <-chan error
}

type ObjLayerAsyncUploader interface {
	AsyncUploader
	SetObjLayer(minio.ObjectLayer)
}

//----------------------------------------------------------------------------------------------------------------------
type folderUploader struct {
	fileUploader
	filesys.DirReader
	l.Logger
}

func NewFolderUploader(ol minio.ObjectLayer, hr filesys.HashFileReader, dr filesys.DirReader, lg l.Logger) ObjLayerAsyncUploader {
	return &folderUploader{fileUploader{ObjectUploader{ol}, hr}, dr, lg}
}

func (f *folderUploader) SetObjLayer(layer minio.ObjectLayer) {
	if f.ol != nil {
		return
	}

	f.ol = layer
}

func (f *folderUploader) uploadDir(ctx context.PutContext, bucket string, dir filesys.Dir) {
	files := dir.Files()
	for i := 0; i < len(files); i++ {
		item := files[i]

		res := <-f.fileUploader.UploadFileAsync(ctx, bucket, path.Join(dir.Path(), item.Name()))
		f.LogE(res.Err)
		if res.Err == nil {
			f.Log(fmt.Sprintf("Successfully uploaded %s", res.Oi.Name))
		}
	}

	if !ctx.Recursive(){
		return
	}

	dirs := dir.Dirs()
	for i := 0; i < len(dirs); i++ {
		item := dirs[i]

		ctxf := ctx.WithPrefixPath(
			utils.GetObjectName(item.Name(), ctx.Prefix(), ctx.Delimiter()),
			path.Join(ctx.Path(), item.Name()))

		f.Log(fmt.Sprintf("Recursively uploading folder %s", ctxf.Path()))
		_ = <-f.UploadFolderAsync(ctxf, bucket, path.Join(dir.Path(), item.Name()))

	}
}

func (f *folderUploader) UploadFileAsync(ctx context.PutContext, bucket, lpath string) <-chan error {
	derrc := make(chan error, 1)
	resc := f.fileUploader.UploadFileAsync(ctx, bucket, lpath)

	go func() {
		res := <-resc
		if res.Err == nil {
			//f.Log(fmt.Sprintf("Successfully uploaded %s", res.Oi.Name))
		}

		derrc <- res.Err
	}()

	return derrc
}

func (f *folderUploader) UploadFolderAsync(ctx context.PutContext, bucket, lpath string) <-chan error {
	derrc := make(chan error, 1)

	dir, err := f.ReadDir(lpath)
	if err != nil {
		derrc <- err
		return derrc
	}

	dirsLen := len(dir.Dirs())
	filesLen := len(dir.Files())
	if dirsLen == 0 && filesLen == 0 {
		f.Log(fmt.Sprintf("Folder %s is empty, no files to upload", lpath))
		derrc <- nil
		return derrc
	}

	if !ctx.Recursive() {
		if dirsLen > 0 {
			f.Log(fmt.Sprintf("Found %d folders, add -r flag to recursively download them", dirsLen))
		}
	}

	//upload
	go func() {
		defer func() {
			derrc <- nil
		}()

		f.uploadDir(ctx, bucket, dir)
	}()

	return derrc
}

//---------------------------------------------------------------------------------------------------------