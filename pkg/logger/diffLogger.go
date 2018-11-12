// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package logger

import (
	"storj/ditto/pkg/models"
)

type DiffLogger interface {
	Logger
	LogDiff([]models.DiffModel)
	GetDiff() []models.DiffModel
}

type diffLogger struct {
	lg
	diff []models.DiffModel
}

func (d *diffLogger) LogDiff(diff []models.DiffModel) {
	d.diff = diff
}

func (d *diffLogger) GetDiff() ([]models.DiffModel) {
	return d.diff
}

var DLogger = diffLogger{
	diff: []models.DiffModel{},
	lg:   StdOutLogger,
}