package api_errors

import "net/http"

var (
	ErrBadRequest = NewAPIError(http.StatusBadRequest, "bad request")
)
