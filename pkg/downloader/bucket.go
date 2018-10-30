package downloader

import (
	"fmt"
	minio "github.com/minio/minio/cmd"
	"path/filepath"
	"storj.io/ditto/cmd/utils"
	"storj.io/ditto/pkg/context"
	"storj.io/ditto/pkg/filesys"
	"storj.io/ditto/pkg/logger"
	"strings"
)

func GetBucket(ctx *context.GetContext, bucket string, downloader MinioDownloader, logger logger.Logger, maker filesys.FsMkdir) (error) {
	return GetPrefix(ctx, bucket, "", downloader, logger, maker)
}

func GetObject(ctx *context.GetContext, bucket, object string, downloader MinioDownloader, logger logger.Logger, maker filesys.FsMkdir) (error) {
	info, err := downloader.Layer().GetObjectInfo(ctx, bucket, object, downloader.Options())
	if err != nil {
		return err
	}

	err = maker.Mkdir(ctx.Path)
	if err != nil {
		return err
	}

	return downloader.DownloadFile(ctx, filepath.Join(ctx.Path, utils.GetFileName(object, ctx.Delimiter)), info)
}

func GetPrefix(ctx *context.GetContext, bucket, prefix string, downloader MinioDownloader, logger logger.Logger, maker filesys.FsMkdir) (error) {
	info, err := GetPrefixInfo(ctx, bucket, prefix, downloader)
	if err != nil {
		return err
	}

	err = maker.Mkdir(ctx.Path)
	if err != nil {
		return err
	}

	for _, obj := range info.Objects {
		name := filepath.Join(ctx.Path, utils.GetFileName(obj.Name, ctx.Delimiter))
		err := downloader.DownloadFile(ctx, name, obj)
		if err != nil {
			logger.LogE(err)
			continue
		}

		logger.Log(fmt.Sprintf("%s downloaded successfully", name))
	}

	for _, pref := range info.Prefixes {
		if ctx.Recursive {
			nctx := context.Clone(ctx)
			nctx.Path = filepath.Join(ctx.Path, strings.TrimSuffix(utils.GetFileName(pref, ctx.Delimiter), ctx.Delimiter))
			nctx.Prefix = pref//strings.Join([]string{ctx.Prefix, pref}, ctx.Delimiter)[len(ctx.Delimiter):]

			err := GetPrefix(nctx, bucket, nctx.Prefix, downloader, logger, maker)
			if err != nil {
				logger.LogE(err)
			}

			continue
		}

		logger.Log(fmt.Sprintf("Found new prefix %s, set -r flag to download it", pref))
	}

	return nil
}

func GetPrefixInfo(ctx *context.GetContext, bucket, prefix string, holder LayerHolder) (minio.ListObjectsV2Info, error) {
	info, err := holder.Layer().ListObjectsV2(ctx, bucket, prefix, "", ctx.Delimiter, ctx.MaxKeys, false, "")
	if err != nil {
		return minio.ListObjectsV2Info{}, err
	}

	return info, nil
}
