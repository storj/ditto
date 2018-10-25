// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/pkg/testutil"
	"testing"
)

func TestValidateSetArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		err  string
	}{
		{
			name: "Empy args",
			args: []string{},
			err:  "Two arguments expected",
		},
		{
			name: "Too many args",
			args: []string{"a", "a", "a"},
			err:  "Two arguments expected",
		},
		{
			name: "Valid case",
			args: []string{"a", "a"},
			err:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateSetArgs(nil, test.args)

			if test.err == "" {
				testutil.AssertNil(t, err)
			} else {
				assert.Equal(t, err.Error(), test.err)
			}
		})
	}
}

func TestExecuteSetCmd(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name: "Two arguments expected",
			args: []string{
				"a",
			},
			expectedError: "Two arguments expected",
		},
		{
			name: "Key is not exist",
			args: []string{
				"a",
				"a",
			},
			expectedError: "Key is not exist",
		},
		{
			name: "Valid case for Credentials",
			args: []string{
				"Server1.Endpoint",
				"testEndpoint",
			},
			expectedError: "",
		},
		{
			name: "Valid case for Predefined options",
			args: []string{
				"DefaultOptions.DefaultSource",
				"server1",
			},
			expectedError: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writeConfigMethod = func(key string, value string) error {
				return nil
			}

			err := executeSetCmd(nil, test.args)

			if test.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, err.Error(), test.expectedError)
			}

		})
	}
}

func TestSetPredefinedOptions(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		value          string
		possibleValues []string
		expectedResult string
	}{
		{
			name:           "Wrong arguments",
			key:            "",
			value:          "a",
			possibleValues: []string{"b"},
			expectedResult: "Only these arguments accepted: [b]",
		},
		{
			name:           "Valid case",
			key:            "a",
			value:          "a",
			possibleValues: []string{"a"},
			expectedResult: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writeConfigMethod = func(key string, value string) error {
				return nil
			}

			err := setPredefinedOptions(test.key, test.value, test.possibleValues)

			if test.expectedResult == "" {
				testutil.AssertNil(t, err)
			} else {
				assert.Equal(t, err.Error(), test.expectedResult)
			}
		})
	}
}
