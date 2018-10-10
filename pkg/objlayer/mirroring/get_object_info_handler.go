// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package mirroring

import (
	"context"
	minio "github.com/minio/minio/cmd"
)

func NewGetObjectInfoHandler(m 	    *MirroringObjectLayer,
							 ctx    context.Context,
							 bucket string,
							 object string,
							 opts   minio.ObjectOptions) *getObjectInfoHandler {

	h := &getObjectInfoHandler{}

	h.m      = m
	h.ctx    = ctx
	h.bucket = bucket
	h.object = object
	h.opts   = opts

	return h
}

type getObjectInfoHandler struct {
	baseHandler
	bucket      string
	object      string
	primeInfo   minio.ObjectInfo
	alterInfo   minio.ObjectInfo
	opts        minio.ObjectOptions
}

func (h *getObjectInfoHandler) execPrime() *getObjectInfoHandler {
	h.primeInfo, h.primeErr = h.m.Prime.GetObjectInfo(h.ctx, h.bucket, h.object, h.opts)

	return h
}

func (h *getObjectInfoHandler) execAlter() *getObjectInfoHandler {
	h.alterInfo, h.alterErr = h.m.Alter.GetObjectInfo(h.ctx, h.bucket, h.object, h.opts)

	return h
}

func (h *getObjectInfoHandler) Process () (objInfo minio.ObjectInfo, err error) {

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
