package utils

import (
	"github.com/minio/minio/pkg/auth"
	"strings"

	"storj/ditto/pkg/config"
	"storj/ditto/pkg/gateway"

	minio "github.com/minio/minio/cmd"
	l "storj/ditto/pkg/logger"
)

func GetObjectLayer() (minio.ObjectLayer, error) {
	config, err := config.ParseConfig()
	if err != nil {
		return nil, err
	}

	logger := l.StdOutLogger
	mirroring := &gateway.Mirroring{Logger: &logger, Config: config}
	objLayer, err := mirroring.NewGatewayLayer(auth.Credentials{})

	if err != nil {
		return nil, err
	}

	return objLayer, nil
}

type GetwayResolver func(l.Logger) (minio.Gateway, error)

func GetGateway(logger l.Logger) (minio.Gateway, error) {
	defaultConfig, err := config.ParseConfig()
	if err != nil {
		return nil, err
	}
	return &gateway.Mirroring{Logger: logger, Config: defaultConfig}, nil
}

func GetObjectName(fname, prefix, delimiter string) (string) {
	if prefix == "" {
		return fname
	}

	if delimiter == "" {
		return fname
	}

	return strings.Join([]string{prefix, fname}, delimiter)
}

func GetFileName(object, delimiter string) (string) {
	elems := strings.Split(object, delimiter)

	elemlen := len(elems)
	for i := elemlen; i > 0; i-- {
		value := elems[i-1]
		if value != "" {
			return value
		}
	}

	return object
}

func AppendPrefix(base, prefix, delimiter string) (string) {
	prefix = strings.Trim(prefix, delimiter)

	if base == "" {
		return strings.Join([]string{prefix, ""}, delimiter)
	}

	return strings.Join([]string{base, prefix, ""}, delimiter)
}

func AppendObject(prefix, object, delimiter string) (string) {
	return strings.TrimSuffix(AppendPrefix(prefix, strings.TrimPrefix(object, prefix), delimiter), delimiter)
}