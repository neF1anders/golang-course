package http

import (
	"net/http"
	"strings"

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

	case codes.Unknown:
		msg := st.Message()
		if strings.Contains(msg, "404") {
			http.Error(w, "repository not found", http.StatusNotFound)
			return
		}
		if strings.Contains(msg, "403") {
			http.Error(w, "access forbidden", http.StatusForbidden)
			return
		}
		if strings.Contains(msg, "429") {
			http.Error(w, "rate limit exceeded", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, msg, http.StatusInternalServerError)

	default:
		http.Error(w, "unknown error", http.StatusInternalServerError)
	}
}
