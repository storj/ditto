// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package testing_utils

import "storj/ditto/pkg/models"

type MockDiffLogger struct {
	MockLogger

	diff []models.DiffModel
}

func (d *MockDiffLogger) LogDiff(diff []models.DiffModel) {
	d.diff = diff
}

func (d *MockDiffLogger) GetDiff() ([]models.DiffModel) {
	return d.diff
}