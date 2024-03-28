package server

import (
	"encoding/json"
	"errors"
	"goddd/internal"
	"goddd/pkg/errx"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type errorResponse struct {
	Error     string            `json:"error"`
	Status    string            `json:"status"`
	HTTPCode  int               `json:"http_code"`
	Datetime  string            `json:"datetime"`
	Timestamp int64             `json:"timestamp"`
	Details   validation.Errors `json:"details,omitempty" swaggertype:"object"`
}

func encodeError(w http.ResponseWriter, err error) {
	var (
		verr validation.Errors
		werr errx.Error
		resp = errorResponse{
			Status:    "fail",
			Datetime:  time.Now().Format("2006-01-02 15:04:05"),
			Timestamp: time.Now().Unix(),
		}
	)

	if !errors.As(err, &werr) {
		resp.HTTPCode = http.StatusInternalServerError
		resp.Error = "internal server error"
	} else {
		switch werr.Code() {
		case internal.CodeInvalidArgument:
			resp.HTTPCode = http.StatusBadRequest
			if errors.As(werr, &verr) {
				resp.Details = verr
			}
		case internal.CodeInvalidToken:
			resp.HTTPCode = http.StatusUnauthorized
		case internal.CodeNoRows, internal.CodeNotFound:
			resp.HTTPCode = http.StatusNotFound
		}
		resp.Error = werr.Message()
	}
	respond(w, resp, resp.HTTPCode)
}

func respond(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
