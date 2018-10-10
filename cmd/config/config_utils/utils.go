// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config_utils

import "storj.io/ditto/pkg/config"

func ContainsKey(key string) bool {
	for _, value := range config.GetKeysArray() {
		if value == key {
			return true
		}
	}

	return false
}