package downloader

import (
	"context"
	minio "github.com/minio/minio/cmd"
	"storj.io/ditto/pkg/filesys"
)

type FileDownloader interface {
	DownloadFile(ctx context.Context, name string, info minio.ObjectInfo) (error)
}

type MinioDownloader interface {
	LayerHolder
	FileDownloader
}

func ForceFileDownloader(layer LayerObjectDownloader) (MinioDownloader) {
	return &fileDownloader{
		layer,
		filesys.ForceFileCreator(),
		filesys.FileRemover(),
	}
}

func NFileDownloader(layer LayerObjectDownloader) (MinioDownloader) {
	return &fileDownloader{
		layer,
		filesys.FileCreator(),
		filesys.FileRemover(),
	}
}

type fileDownloader struct {
	LayerObjectDownloader
	filesys.FsCreate
	filesys.FsRemove
}

func (d *fileDownloader) DownloadFile(ctx context.Context, name string, info minio.ObjectInfo) (error) {
	file, err := d.Create(name)
	if err != nil {
		return err
	}

	err = d.DownloadObject(ctx, file, info)
	if err != nil {
		d.Remove(name)
	}

	return nil
}
