package get

import (
	"io"
	"os"
	"github.com/pkg/errors"
	"fmt"
	minio "github.com/minio/minio/cmd"
	"context"
)

type objectDownloader struct {
	FileCreator
	minio.ObjectLayer
}

func (d *objectDownloader) DownloadObject(ctx context.Context, lpath string, oi minio.ObjectInfo, opts minio.ObjectOptions) (err error) {
	file, err := d.Create(lpath)
	if err != nil {
		return
	}

	err = d.GetObject(ctx, oi.Bucket, oi.Name, 0, oi.Size, file, oi.ETag, opts)
	return
}

type ObjectInfoGetter interface {}

type FileCreator interface {
	Create(lpath string) (io.WriteCloser, error)
}

type baseFileCreator struct {

}

func (c *baseFileCreator) Create(lpath string) (io.WriteCloser, error) {
	file, err := os.Create(lpath)
	return file, err
}

type newFileCreator struct {
	f *baseFileCreator
}

func (c *newFileCreator) CreateFile(lpath string) (io.WriteCloser, error) {
	fi, err := os.Stat(lpath)
	if err == nil {
		return nil, errors.New(fmt.Sprintf("file %s allready exists", fi.Name()))
	}

	f, err := c.f.Create(lpath)
	return f, err
}

