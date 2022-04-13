package http

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

// response creates http response
type response struct {
	logger *zap.Logger
}

// NewResponse inits and returns Response struct
func NewResponse(logger *zap.Logger) *response {
	return &response{logger: logger}
}

// createResponse creates 200 response
func (r *response) createResponse(w http.ResponseWriter, req *http.Request, v interface{}) {
	w.WriteHeader(http.StatusOK)
	r.write(w, req, v)
}

// createInternalErrorResponse creates 500 response
func (r *response) createInternalErrorResponse(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	r.write(w, req, err)
}

// createBadRequestResponse creates 403 response
func (r *response) createBadRequestResponse(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	r.write(w, req, err)
}

func (r *response) write(w http.ResponseWriter, req *http.Request, v interface{}) {
	r.setCors(req)

	resp, err := json.Marshal(v)
	if err != nil {
		r.logger.Error("can't marshal http response", zap.Error(err))
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		r.logger.Error("can't write http response", zap.Error(err))
		return
	}
}

func (r *response) setCors(req *http.Request) {
	origin := req.Header.Get("origin")
	if origin != "" {
		req.Header.Set("Access-Control-Allow-Origin", origin)
		req.Header.Set("Access-Control-Allow-Credentials", "true")
		req.Header.Set("Access-Control-Allow-Headers", "Content-type, Authorization")
	}
}
