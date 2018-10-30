package context

import "context"

type GetContext struct {
	context.Context
	Path, Prefix     string
	Delimiter        string
	Recursive, Force bool
	MaxKeys          int
}

func Clone(ctx *GetContext) *GetContext {
	return &GetContext{
		ctx.Context,
		ctx.Path,
		ctx.Prefix,
		ctx.Delimiter,
		ctx.Recursive,
		ctx.Force,
		ctx.MaxKeys,
	}
}
