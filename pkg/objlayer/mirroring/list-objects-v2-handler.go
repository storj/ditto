// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	"storj.io/ditto/pkg/utils"

	minio "github.com/minio/minio/cmd"
	l "storj.io/ditto/pkg/logger"
)

func NewListObjectsV2Handler(m   *MirroringObjectLayer,
							 ctx context.Context,
							 bucket, prefix, cntnToken, delimiter, startAfter string,
							 maxKeys int,
							 fetchOwner bool) *listObjectsV2Handler {

	h := &listObjectsV2Handler{}

	h.m = m
	h.ctx =  ctx
	h.bucket = bucket
	h.prefix = prefix
	h.cntnToken = cntnToken
	h.delimiter = delimiter
	h.maxKeys  = maxKeys
	h.startAfter = startAfter
	h.fetchOwner = fetchOwner

	return h
}

type listObjectsV2Handler struct {
	baseHandler
	bucket      string
	prefix      string
	cntnToken   string
	delimiter   string
	maxKeys 	int
	fetchOwner  bool
	startAfter  string
	primeInfo   *minio.ListObjectsV2Info
	alterInfo   *minio.ListObjectsV2Info
}

func (h *listObjectsV2Handler) execPrime() *listObjectsV2Handler {

	primeInfo, primeErr := h.m.Prime.ListObjectsV2(h.ctx, h.bucket, h.prefix, h.cntnToken, h.delimiter, h.maxKeys, h.fetchOwner, h.startAfter)

	h.primeInfo, h.primeErr = &primeInfo, primeErr

	return h
}

func (h *listObjectsV2Handler) execAlter() *listObjectsV2Handler {
	alterInfo, alterErr := h.m.Alter.ListObjectsV2(h.ctx, h.bucket, h.prefix, h.cntnToken, h.delimiter, h.maxKeys, h.fetchOwner, h.startAfter)

	h.alterInfo, h.alterErr = &alterInfo, alterErr

	return h
}

func (h *listObjectsV2Handler) Process () (minio.ListObjectsV2Info, error) {

	h.execPrime()

	switch {
	case h.m.Config.ListOptions.Merge:
		return h.merge()

	case !h.m.Config.ListOptions.DefaultOptions.ThrowImmediately:
		return h.retry()
	}

	return *h.primeInfo, h.primeErr
}

func (h *listObjectsV2Handler) retry() (minio.ListObjectsV2Info, error) {
	if h.primeErr != nil {

		h.m.Logger.LogE(h.primeErr)

		h.execAlter()

		if h.alterErr != nil {

			return minio.ListObjectsV2Info{}, utils.CombineErrors([]error{ h.alterErr, h.primeErr })
		}

		return *h.alterInfo, nil
	}

	return *h.primeInfo, nil
}

func (h *listObjectsV2Handler) merge() (minio.ListObjectsV2Info, error) {
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
		return minio.ListObjectsV2Info{}, utils.CombineErrors([]error{ h.alterErr, h.primeErr })
	}

	mergedObjects := utils.CombineObjectsDistinct(h.primeInfo.Objects, h.alterInfo.Objects)

	mergedResult := minio.ListObjectsV2Info{
		Objects:     		   mergedObjects,
		Prefixes:    	       h.primeInfo.Prefixes,
		IsTruncated: 	   	   h.primeInfo.IsTruncated,
		ContinuationToken: 	   h.primeInfo.ContinuationToken,
		NextContinuationToken: h.primeInfo.NextContinuationToken,
	}

	h.logDiff()

	return mergedResult, nil
}

func (h *listObjectsV2Handler) logDiff() {

	diff := utils.ListObjectsWithDifference(h.primeInfo.Objects, h.alterInfo.Objects)

	l, ok := h.m.Logger.(l.DiffLogger)

	if ok {
		l.LogDiff(diff)
	}
}