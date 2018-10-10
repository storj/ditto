package put

import (
	l "storj.io/ditto/pkg/logger"
	futils "storj.io/ditto/cmd/utils"
	"fmt"
	"path"
	"os"
	"context"
	"github.com/minio/minio/pkg/auth"
	"github.com/spf13/cobra"
	minio "github.com/minio/minio/cmd"
	"storj.io/ditto/pkg/uploader"
	fsystem "storj.io/ditto/pkg/filesys"
	dcontext "storj.io/ditto/pkg/context"

	"storj.io/mirroring/utils"
)

type putExec struct {
	gw minio.Gateway
	uploader.ObjLayerAsyncUploader
	fsystem.DirChecker
	logger l.Logger
}

func NewPutExec(gw minio.Gateway, logger utils.Logger) putExec {
	uploader := uploader.NewFolderUploader(nil, fsystem.NewHashFileReader(), fsystem.NewDirReader(), logger)
	return newPutExec(gw, uploader, fsystem.BDirChecker(futils.CheckIfDir), logger)
}

func newPutExec(gw minio.Gateway, uploader uploader.ObjLayerAsyncUploader, dirChecker fsystem.DirChecker, logger utils.Logger) putExec {
	return putExec{gw, uploader, dirChecker, logger }
}

func (e putExec) logF(format string, params ...interface{}) {
	e.logger.Log(fmt.Sprintf(format, params...))
}

//Main function
func (e putExec) runE(cmd *cobra.Command, args []string) error {
	mirr, err := e.gw.NewGatewayLayer(auth.Credentials{})
	if err != nil {
		return err
	}

	e.SetObjLayer(mirr)

	bctx := context.Background()
	_, err = mirr.GetBucketInfo(bctx, args[0])
	if err != nil {
		return err
	}

	isDir, err := e.CheckIfDir(args[1])
	if err != nil {
		return err
	}

	ctx, cancelf := context.WithCancel(bctx)
	defer func() {
		cancelf()
	}()
	
	cwd, _ := os.Getwd()
	lpath := path.Join(cwd, args[1])

	ctxp := dcontext.NewPutCtx(
		ctx,
		frecursive,
		fforce,
		fprefix,
		fdelimiter)

	var errc <-chan error
	if isDir {
		errc = e.UploadFolderAsync(ctxp, args[0], lpath)
	} else {
		errc = e.UploadFileAsync(ctxp, args[0], lpath)
	}

	tnum := 1
	for i:= 0; i < tnum; i++ {
		select {
		case err = <-errc:
			e.logger.LogE(err)
		case sig := <-sigc:
			e.logF("Catched interrupt! %s\n", sig)
			cancelf()
			tnum++
		}
	}

	return err
}