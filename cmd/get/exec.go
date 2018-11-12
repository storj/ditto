package get

import (
	"context"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/pkg/auth"
	"github.com/spf13/cobra"
	"path/filepath"
	"storj/ditto/cmd/utils"
	dcontext "storj/ditto/pkg/context"
	"storj/ditto/pkg/downloader"
	"storj/ditto/pkg/filesys"
	"storj/ditto/pkg/logger"
)

func NewGetExec(resolver utils.GetwayResolver, lg logger.Logger) *getExec {
	return &getExec{resolver, lg}
}

type getExec struct {
	utils.GetwayResolver
	logger.Logger
}

func (e *getExec) runE(cmd *cobra.Command, args []string) (err error) {
	gw, err := e.GetwayResolver(e.Logger)
	if err != nil {
		return
	}

	mirr, err := gw.NewGatewayLayer(auth.Credentials{})
	if err != nil {
		return
	}

	ctx := context.Background()
	getCtx := &dcontext.GetContext{
		ctx,
		filepath.Clean(locationFlag),
		prefixFlag,
		delimiterFlag,
		recursiveFlag,
		forceFlag,
		maxKeysFlag,
	}

	var miniod downloader.MinioDownloader
	if forceFlag {
		miniod = downloader.ForceFileDownloader(downloader.NewObjectDownloader(mirr, minio.ObjectOptions{}))
	} else {
		miniod = downloader.NFileDownloader(downloader.NewObjectDownloader(mirr, minio.ObjectOptions{}))
	}

	argsLen := len(args)
	if argsLen == 2 {
		objname := utils.AppendObject(getCtx.Prefix, args[1], getCtx.Delimiter)
		return downloader.GetObject(getCtx, args[0], objname, miniod, e.Logger, filesys.DirMaker())
	}


	return downloader.GetPrefix(getCtx, args[0], prefixFlag, miniod, e.Logger, filesys.DirMaker())
}
