// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	"storj.io/ditto/pkg/utils"

	minio "github.com/minio/minio/cmd"
	l "storj.io/ditto/pkg/logger"
)

func NewListObjectsHandler(m   *MirroringObjectLayer,
	                       ctx context.Context,
	                       bucket, prefix, marker, delimiter string,
                           maxKeys int) *listObjectsHandler {

	h := &listObjectsHandler{}

	h.m = m
	h.ctx =  ctx
	h.bucket = bucket
	h.prefix = prefix
	h.marker = marker
	h.delimiter = delimiter
	h.maxKeys  = maxKeys

	return h
}

type listObjectsHandler struct {
	baseHandler
	bucket      string
	prefix      string
	marker      string
	delimiter   string
	maxKeys 	int
	primeInfo   *minio.ListObjectsInfo
	alterInfo   *minio.ListObjectsInfo
}

func (h *listObjectsHandler) execPrime() *listObjectsHandler {

	primeInfo, primeErr := h.m.Prime.ListObjects(h.ctx, h.bucket, h.prefix, h.marker, h.delimiter, h.maxKeys)

	h.primeInfo, h.primeErr = &primeInfo, primeErr

	return h
}

func (h *listObjectsHandler) execAlter() *listObjectsHandler {
	alterInfo, alterErr := h.m.Alter.ListObjects(h.ctx, h.bucket, h.prefix, h.marker, h.delimiter, h.maxKeys)

	h.alterInfo, h.alterErr = &alterInfo, alterErr

	return h
}

func (h *listObjectsHandler) Process () (minio.ListObjectsInfo, error) {

	h.execPrime()

	switch {
		case h.m.Config.ListOptions.Merge:
			return h.merge()

		case !h.m.Config.ListOptions.DefaultOptions.ThrowImmediately:
			return h.retry()
	}

	return *h.primeInfo, h.primeErr
}

func (h *listObjectsHandler) retry() (minio.ListObjectsInfo, error) {
	if h.primeErr != nil {

		h.m.Logger.LogE(h.primeErr)

		h.execAlter()

		if h.alterErr != nil {

			return minio.ListObjectsInfo{}, utils.CombineErrors([]error{ h.alterErr, h.primeErr })
		}

		return *h.alterInfo, nil
	}

	return *h.primeInfo, nil
}

func (h *listObjectsHandler) merge() (minio.ListObjectsInfo, error) {
	h.execAlter()

	if h.primeErr != nil && h.alterErr == nil {

		h.m.Logger.LogE(h.primeErr)

		return *h.alterInfo, nil
	}

	if h.alterErr != nil && h.primeErr == nil {

		h.m.Logger.LogE(h.alterErr)

		return *h.primeInfo, nil
	}

	if h.alterErr != nil && h.primeErr != nil {
		return minio.ListObjectsInfo{}, utils.CombineErrors([]error{ h.alterErr, h.primeErr })
	}

	mergedObjects := utils.CombineObjectsDistinct(h.primeInfo.Objects, h.alterInfo.Objects)

	mergedResult := minio.ListObjectsInfo{
		Objects:     mergedObjects,
		Prefixes:    h.primeInfo.Prefixes,
		IsTruncated: h.primeInfo.IsTruncated,
		NextMarker:  h.primeInfo.NextMarker,
	}


	h.logDiff()

	return mergedResult, nil
}

func (h *listObjectsHandler) logDiff() {

	diff := utils.ListObjectsWithDifference(h.primeInfo.Objects, h.alterInfo.Objects)

	l, ok := h.m.Logger.(l.DiffLogger)

	if ok {
		l.LogDiff(diff)
	}
}