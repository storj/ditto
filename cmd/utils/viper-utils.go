// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package utils

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadFlagValueFromViperIfNotSet(cmd *cobra.Command, flagName string, viperName string) {
	if cmd == nil {
		return
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flagName != flag.Name {
			return
		}

		if flag.Changed {
			return
		}

		viperValue := viper.Get(viperName)
		if viperValue == nil {
			return
		}

		flagValue, err := cast.ToStringE(viperValue)
		if err != nil {
			// fmt.Printf("err set pflag %s from viper. err: %s\n", flag.Name, err)
			return
		}

		err = cmd.Flags().Set(flag.Name, flagValue)
		if err == nil {
			// fmt.Printf("set pflag %s from viper\n", flag.Name)
		} else {
			//
		}
	})
}
