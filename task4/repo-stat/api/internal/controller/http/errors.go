package http

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func writeGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)

	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)

	case codes.Unavailable:
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)

	case codes.Internal:
		http.Error(w, "internal server error", http.StatusInternalServerError)

	default:
		http.Error(w, "unknown error", http.StatusInternalServerError)
	}
}
