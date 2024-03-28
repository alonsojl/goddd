package internal

import (
	"goddd/pkg/errx"
)

const (
	CodeUnknown errx.ErrorCode = iota
	CodeInvalidArgument
	CodeInvalidToken
	CodeNotFound
	CodeNoRows
	CodeDecode
)

var (
	ErrNoRows = errx.NewErrorf(CodeNoRows, "no rows")
)
