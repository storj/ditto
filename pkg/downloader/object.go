package downloader

import (
	"context"
	minio "github.com/minio/minio/cmd"
	"io"
)

type LayerHolder interface {
	Layer() minio.ObjectLayer
	Options() minio.ObjectOptions
}

type ObjectDownloader interface {
	DownloadObject(ctx context.Context, writer io.Writer, info minio.ObjectInfo) error
}

type LayerObjectDownloader interface {
	LayerHolder
	ObjectDownloader
}

func NewObjectDownloader(layer minio.ObjectLayer, opts minio.ObjectOptions) LayerObjectDownloader {
	return &objectDownloader{layer, opts}
}

type objectDownloader struct {
	minio.ObjectLayer
	opts minio.ObjectOptions
}

func (d *objectDownloader) Layer() minio.ObjectLayer {
	return d.ObjectLayer
}

func (d *objectDownloader) Options() minio.ObjectOptions {
	return d.opts
}

func (d *objectDownloader) DownloadObject(ctx context.Context, writer io.Writer, info minio.ObjectInfo) error {
	return d.GetObject(ctx, info.Bucket, info.Name, 0, info.Size, writer, info.ETag, d.opts)
}




