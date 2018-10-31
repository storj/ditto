package gateway

import (
	"errors"
	"fmt"
	"github.com/minio/cli"
	"github.com/minio/minio/pkg/auth"
	"storj.io/ditto/pkg/config"
	"storj.io/ditto/pkg/objlayer/mirroring"

	minio "github.com/minio/minio/cmd"
	l "storj.io/ditto/pkg/logger"
	s3 "storj.io/ditto/pkg/objlayer/s3compat"
)

func init() {
	err := minio.RegisterGatewayCommand(cli.Command{
		Name:            "mirroring",
		Usage:           "mirroring",
		Action:          mirroringGatewayMain,
		HideHelpCommand: true,
	})

	if err != nil {
	}
}

func mirroringGatewayMain(ctx *cli.Context) {
	err := config.ReadConfig(true)
	if err != nil {
		fmt.Println("error reading config")
	}
	config, err := config.ParseConfig()
	if err != nil {
		return
	}

	if err != nil {
		println("error while opening Gateway", err)
		return
	}

	minio.StartGateway(ctx, &Mirroring{Logger: &l.StdOutLogger, Config: config})
}

// Mirroring for mirroring service
type Mirroring struct {
	Config *config.Config
	Logger l.Logger
}

// Name implements minio.Gateway interface
func (gw *Mirroring) Name() string {
	return ""
}

// NewGatewayLayer implements minio.Gateway interface
func (gw *Mirroring) NewGatewayLayer(creds auth.Credentials) (objLayer minio.ObjectLayer, err error) {
	if gw.Config == nil {
		return nil, errors.New("configuration is not set")
	}

	s1Credentials := gw.Config.Server1
	prime, err := s3.NewS3Compat(s1Credentials.Endpoint, s1Credentials.AccessKey, s1Credentials.SecretKey)

	if err != nil {
		return nil, err
	}

	s2Credentials := gw.Config.Server2
	alter, err := s3.NewS3Compat(s2Credentials.Endpoint, s2Credentials.AccessKey, s2Credentials.SecretKey)

	if err != nil {
		return nil, err
	}

	objLayer = &mirroring.MirroringObjectLayer{
		Prime:  prime,
		Alter:  alter,
		Logger: gw.Logger,
		Config: gw.Config,
	}

	return objLayer, nil
}

// Production - both gateways are production ready.
func (gw *Mirroring) Production() bool {
	return false
}
