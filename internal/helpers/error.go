package helpers

import (
	"context"
	"lowerthirdsapi/internal/apierrors"
	"net/http"
)

// WriteError pulls an error response from the context, merges it with the provided error (if applicable),
// and writes the error to the http response writer
func WriteError(ctx context.Context, err error, w http.ResponseWriter) {
	errResp, ok := ctx.Value(ErrorsResponseKey).(*apierrors.Response)
	if ok {
		errResp.Add(err)
	} else {
		errResp = apierrors.NewResponse(err)
	}
	_ = errResp.Write(w)
}
