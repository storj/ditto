// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import "storj.io/ditto/pkg/config"

var argsMap = map[string][]string{
	config.SERVER_1_ENDPOINT:                 {},
	config.SERVER_1_ACCESS_KEY:               {},
	config.SERVER_1_SECRET_KEY:               {},
	config.SERVER_2_ENDPOINT:                 {},
	config.SERVER_2_ACCESS_KEY:               {},
	config.SERVER_2_SECRET_KEY:               {},
	config.DEFAULT_OPTIONS_DEFAULT_SOURCE:    {"server1", "server2"},
	config.DEFAULT_OPTIONS_THROW_IMMEDIATELY: {"true", "false"},
	config.LIST_DEFAULT_SOURCE:               {"server1", "server2"},
	config.LIST_THROW_IMMEDIATELY:            {"true", "false"},
	config.LIST_MERGE:                        {"true", "false"},
	config.PUT_DEFAULT_SOURCE:                {"server1", "server2"},
	config.PUT_THROW_IMMEDIATELY:             {"true", "false"},
	config.PUT_CREATE_BUCKET_IF_NOT_EXIST:    {"true", "false"},
	config.GET_OBJECT_DEFAULT_SOURCE:         {"server1", "server2"},
	config.GET_OBJECT_THROW_IMMEDIATELY:      {"true", "false"},
	config.COPY_DEFAULT_SOURCE:               {"server1", "server2"},
	config.COPY_THROW_IMMEDIATELY:            {"true", "false"},
	config.DELETE_DEFAULT_SOURCE:             {"server1", "server2"},
	config.DELETE_THROW_IMMEDIATELY:          {"true", "false"},
}
