// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/pkg/testutil"
	"storj.io/ditto/pkg/config"
	"testing"
)

func TestValidateGetArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		err  string
	}{
		{
			name: "Empy args",
			args: []string{},
			err:  "Only one argument accepted",
		},
		{
			name: "Too many args",
			args: []string{"a", "a"},
			err:  "Only one argument accepted",
		},
		{
			name: "Valid case",
			args: []string{"a"},
			err:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateGetArgs(nil, test.args)

			if test.err != "" {
				assert.Equal(t, err.Error(), test.err)
			}
		})
	}
}

func TestGetValueFromConfigFile(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		viperInitFunc  func(useDefaults bool) (config *config.Config, err error)
		expectedResult string
	}{
		{
			name: "Key not set",
			key:  "a",
			viperInitFunc: func(useDefaults bool) (config *config.Config, err error) {
				viper.Reset()
				return nil, nil
			},
			expectedResult: "Key is not set",
		},
		{
			name: "Key exist",
			key:  "a",
			viperInitFunc: func(useDefaults bool) (config *config.Config, err error) {
				viper.Reset()
				viper.Set("a", "b")
				return nil, nil
			},
			expectedResult: "b",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			readConfigMethod = test.viperInitFunc
			readConfigMethod(false)
			value := getValueFromConfigFile(test.key)

			assert.Equal(t, value, test.expectedResult)
		})
	}
}

func TestExecuteGetCmd(t *testing.T) {
	tests := []struct {
		name           string
		viperInitFunc  func(useDefaults bool) (config *config.Config, err error)
		args           []string
		expectedResult string
	}{
		{
			name:"Unsupported key",
			viperInitFunc: func(useDefaults bool) (config *config.Config, err error) {
				viper.Reset()

				return nil, nil
			},
			args: []string{"a"},
			expectedResult:"Key unsupported",
		},
		{
			name:"Valid case",
			viperInitFunc: func(useDefaults bool) (config *config.Config, err error) {
				viper.Reset()
				viper.Set("Server1.Endpoint", "b")
				return nil, nil
			},
			args: []string{"Server1.Endpoint"},
			expectedResult:"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			readConfigMethod = test.viperInitFunc
			err := executeGetCmd(nil, test.args)

			if test.expectedResult == "" {
				testutil.AssertNil(t, err)
			} else {
				assert.Equal(t, err.Error(), test.expectedResult)
			}


		})
	}
}
