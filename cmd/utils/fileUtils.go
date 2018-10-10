package utils

import (
	"github.com/minio/minio/pkg/auth"
	"os"
	"strings"

	"storj.io/ditto/pkg/config"
	"storj.io/ditto/pkg/gateway"

	l "storj.io/ditto/pkg/logger"
	minio "github.com/minio/minio/cmd"
)

func FileOrDirExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		//fmt.Printf("Name: %s, isDir: %t", r.Name(), r.IsDir())
		return true, nil
	}

	if os.IsNotExist(err) {

		return false, nil
	}

	return true, err
}

//TODO: implement
func GetObjectLayer() (minio.ObjectLayer, error) {
	defaultConfig, err := config.ReadDefaultConfig(true)
	if err != nil {
		return nil, err
	}

	logger := l.StdOutLogger
	mirroring := &gateway.Mirroring{Logger: &logger, Config: defaultConfig}
	objLayer, err := mirroring.NewGatewayLayer(auth.Credentials{})

	if err != nil {
		return nil, err
	}

	return objLayer, nil
}

func GetObjectName(fname, prefix, delimiter string) string {
	if prefix == "" {
		return fname
	}

	if delimiter == "" {
		return fname
	}

	return strings.Join([]string{prefix, fname}, delimiter)
}
