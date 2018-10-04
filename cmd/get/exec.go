package get

import (
	"github.com/spf13/cobra"
	minio "github.com/minio/minio/cmd"
	"storj.io/ditto/utils"
	"github.com/minio/minio/pkg/auth"
	"context"
	"os"
	)

func newGetExec(gw minio.Gateway, lg utils.Logger) *getExec {
	return &getExec{gw, lg}
}

type getExec struct {
	minio.Gateway
	utils.Logger
}

func (e *getExec) runE(cmd *cobra.Command, args []string) (err error) {
	mirr, err := e.NewGatewayLayer(auth.Credentials{})
	if err != nil {
		return
	}

	ctx := context.Background()

	file, err := os.Create(args[1])
	if err != nil {
		return
	}

	_, err = file.Stat()
	if err != nil {
		return
	}

	err = mirr.GetObject(ctx, args[0], args[1], 0, int64(3000000), file, "", minio.ObjectOptions{})
	return err
}