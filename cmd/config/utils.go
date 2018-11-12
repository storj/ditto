// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import "storj/ditto/pkg/config"

func containsKey(key string) bool {
	for _, value := range config.GetKeysArray() {
		if value == key {
			return true
		}
	}

	return false
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}