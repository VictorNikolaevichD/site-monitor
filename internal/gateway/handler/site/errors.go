package site

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, httpStatus int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func mapGRPCError(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, "internal server error"
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, st.Message()
	case codes.NotFound:
		return http.StatusNotFound, st.Message()
	case codes.AlreadyExists:
		return http.StatusConflict, st.Message()
	case codes.Unavailable:
		return http.StatusBadGateway, "monitor service unavailable"
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, "monitor service timeout"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
