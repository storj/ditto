// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	minio "github.com/minio/minio/cmd"
)

func NewGetBucketInfoHandler(m 	     *MirroringObjectLayer,
							 ctx     context.Context,
							 bucket  string) *getBucketInfoHandler {

	h := &getBucketInfoHandler{}

	h.m = m
	h.ctx =  ctx
	h.bucket = bucket

	return h
}

type getBucketInfoHandler struct {
	baseHandler
	bucket      string
	primeInfo   minio.BucketInfo
	alterInfo   minio.BucketInfo
}

func (h *getBucketInfoHandler) execPrime() *getBucketInfoHandler {
	h.primeInfo, h.primeErr = h.m.Prime.GetBucketInfo(h.ctx, h.bucket)

	return h
}

func (h *getBucketInfoHandler) execAlter() *getBucketInfoHandler {
	h.alterInfo, h.alterErr = h.m.Alter.GetBucketInfo(h.ctx, h.bucket)

	return h
}

func (h *getBucketInfoHandler) Process () (objInfo minio.BucketInfo, err error) {

	h.execPrime()

	if h.primeErr == nil {
		return h.primeInfo, nil
	}

	h.m.Logger.LogE(h.primeErr)

	h.execAlter()

	if h.alterErr != nil {

		h.m.Logger.LogE(h.alterErr)

		return h.alterInfo, h.primeErr
	}

	return h.alterInfo, nil
}


