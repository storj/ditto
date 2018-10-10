// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	"storj.io/ditto/pkg/utils"

	l "storj.io/ditto/pkg/logger"
	minio "github.com/minio/minio/cmd"
)

func NewListBucketsHandler(m   *MirroringObjectLayer,
						   ctx context.Context) *listBucketsHandler {

	h := &listBucketsHandler{}

	h.m = m
	h.ctx =  ctx

	return h
}

type listBucketsHandler struct {
	baseHandler
	primeBuckets []minio.BucketInfo
	alterBuckets []minio.BucketInfo
}

func (h *listBucketsHandler) execPrime() *listBucketsHandler {

	h.primeBuckets, h.primeErr = h.m.Prime.ListBuckets(h.ctx)

	return h
}

func (h *listBucketsHandler) execAlter() *listBucketsHandler {
	h.alterBuckets, h.alterErr = h.m.Alter.ListBuckets(h.ctx)

	return h
}

func (h *listBucketsHandler) Process () ([]minio.BucketInfo, error) {

	h.execPrime()

	switch {
		case h.m.Config.ListOptions.Merge:
			return h.withMerge()

		case !h.m.Config.ListOptions.DefaultOptions.ThrowImmediately:
			return h.withoutMerge()
	}

	return h.primeBuckets, h.primeErr
}

func (h *listBucketsHandler) withoutMerge() ([]minio.BucketInfo, error) {
	if h.primeErr != nil {

		h.m.Logger.LogE(h.primeErr)

		h.execAlter()

		if h.alterErr != nil {

			return nil, utils.CombineErrors([]error{ h.alterErr, h.primeErr })
		}

		return h.alterBuckets, nil
	}

	return h.primeBuckets, nil
}

func (h *listBucketsHandler) withMerge() ([]minio.BucketInfo, error) {
	h.execAlter()

	if h.primeErr != nil && h.alterErr == nil {

		h.m.Logger.LogE(h.primeErr)

		return h.alterBuckets, nil
	}

	if h.alterErr != nil && h.primeErr == nil {

		h.m.Logger.LogE(h.alterErr)

		return h.primeBuckets, nil
	}

	if h.alterErr != nil && h.primeErr != nil {
		return nil, utils.CombineErrors([]error{ h.alterErr, h.primeErr })
	}

	mergedBuckets := utils.CombineBucketsDistinct(h.primeBuckets, h.alterBuckets)

	h.logDiff()

	return mergedBuckets, nil
}

func (h *listBucketsHandler) logDiff() {

	diff := utils.ListBucketsWithDifference(h.primeBuckets, h.alterBuckets)

	l, ok := h.m.Logger.(l.DiffLogger)

	if ok {
		l.LogDiff(diff)
	}
}