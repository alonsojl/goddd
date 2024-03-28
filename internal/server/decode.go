package server

import (
	"encoding/json"
	"goddd/internal"
	"goddd/pkg/errx"
	"net/http"
)

var ErrDecode = errx.NewErrorf(internal.CodeInvalidArgument, "fail decode request")

func decode(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return ErrDecode.SetOrigin(err)
	}
	return nil
}
