package mirroring

import (
	minio "github.com/minio/minio/cmd"
	"context"
	"io"
		)

type getAsyncHandler struct {
	ol minio.ObjectLayer
}

func(h getAsyncHandler) GetObjectAsync(ctx context.Context, bucket string, object string, startOffset int64, length int64, writer io.Writer, etag string, opts minio.ObjectOptions) <-chan error {
	errc := make(chan error)
	getTask := func(errc chan<- error) {
		err := h.ol.GetObject(ctx, bucket, object, startOffset, length, writer, etag, opts)
		errc <- err
	}

	go getTask(errc)
	return errc
}

type getHandler struct {
	prime, alter getAsyncHandler
	throwImmediately bool
}

func newGetHandler(prime, alter minio.ObjectLayer, thrImm bool) *getHandler {
	return &getHandler{getAsyncHandler{prime}, getAsyncHandler{alter}, thrImm}
}

func (h *getHandler) process(ctx context.Context, bucket string, object string, startOffset int64, length int64, writer io.Writer, etag string, opts minio.ObjectOptions) (err error) {
	wrtwrap := &writeCounter{w : writer}

	err = <-h.prime.GetObjectAsync(ctx, bucket, object, startOffset, length, wrtwrap, etag, opts)

	if h.throwImmediately {
		return
	}

	if err != nil {
		err = <-h.alter.GetObjectAsync(ctx, bucket, object, startOffset + wrtwrap.bcount, length - wrtwrap.bcount, wrtwrap, etag, opts)
	}

	return
}

type writeCounter struct {
	w io.Writer
	bcount int64
}

func (c *writeCounter) Write(b []byte) (int, error) {
	n, err := c.w.Write(b)
	if err == nil {
		c.bcount += int64(n)
	}

	return n, err
}