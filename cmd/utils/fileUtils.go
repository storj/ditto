package utils

import (
	"strings"
	"github.com/minio/minio/pkg/auth"
	"os"

	"storj.io/ditto/pkg/config"
	"storj.io/ditto/pkg/gateway"
	"storj.io/ditto/utils"

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
	defaultConfig, err := config.ReadDefaultConfig()
	if err != nil {
		return nil, err
	}

	logger := utils.StdOutLogger
	mirroring := &gateway.Mirroring{Logger: &logger, Config: defaultConfig}
	objLayer, err := mirroring.NewGatewayLayer(auth.Credentials{})

	if err != nil {
		return nil, err
	}

	return objLayer, nil
}

func CheckIfDir(lpath string) (isDir bool, err error) {
	fi, err := os.Stat(lpath)
	if err != nil {
		return
	}

	return fi.IsDir(), err
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
