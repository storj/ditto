// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package logger

import (
	"github.com/stretchr/testify/assert"
	"storj.io/ditto/pkg/models"
	"testing"
)

func TestDiffLogger(t *testing.T) {
	cases := []struct {
		testName string
		testFunc func()
	}{
		{
			testName: "Log diff",

			testFunc: func() {

				diff := []models.DiffModel {
					{
						Name: "name1",
						Diff: 1,
					},
					{
						Name: "name2",
						Diff: 2,
					},
					{
						Name: "name3",
						Diff: 3,
					},
				}

				DLogger.LogDiff(diff)

				assert.Equal(t, len(diff), len(DLogger.diff))
				assert.NotNil(t, DLogger.diff)

				for i := 0; i < len(DLogger.diff); i++ {
					assert.Equal(t, diff[i].Name, DLogger.diff[i].Name)
					assert.Equal(t, diff[i].Diff, DLogger.diff[i].Diff)
				}

			},
		},

	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			c.testFunc()
		})
	}
}