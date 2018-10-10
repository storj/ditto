// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package logger

import (
	"storj.io/ditto/pkg/models"
)

type DiffLogger interface {
	Logger
	LogDiff([]models.DiffModel)
}

type diffLogger struct {
	lg

	Diff []models.DiffModel
}

func (d *diffLogger) LogDiff(diff []models.DiffModel) {
	d.bucketDiffBuffer = diff
}

var DLogger = diffLogger{
	Diff: []models.DiffModel{},
	lg: StdOutLogger,
}