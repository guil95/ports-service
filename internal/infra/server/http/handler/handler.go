package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/guil95/ports-service/internal/core/domain"
)

type HTTPHandler struct {
	portService domain.ServicePort
}

func NewHTTPHandler(portService domain.ServicePort) *HTTPHandler {
	return &HTTPHandler{portService: portService}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/ports/"):
		id := strings.TrimPrefix(r.URL.Path, "/ports/")
		if id == "" {
			writeResponse(w, http.StatusBadRequest, nil, methodNotAllowed)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.getPort(w, r, id)
		default:
			writeResponse(w, http.StatusMethodNotAllowed, nil, methodNotAllowed)
		}
	case r.URL.Path == "/ports":
		switch r.Method {
		case http.MethodPost:
			h.createPort(w, r)
		default:
			writeResponse(w, http.StatusMethodNotAllowed, nil, methodNotAllowed)
		}
	default:
		writeResponse(w, http.StatusNotFound, nil, notFound)
	}
}

func (h *HTTPHandler) createPort(w http.ResponseWriter, r *http.Request) {
	var p domain.Port
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, invalidRequest)
		return
	}

	if len(p.Unlocs) == 0 {
		writeResponse(w, http.StatusBadRequest, nil, invalidRequest)
		return
	}

	err := h.portService.CreateOrUpdate(r.Context(), p)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, internalServer)
		return
	}

	writeResponse(w, http.StatusOK, nil, nil)
}

func (h *HTTPHandler) getPort(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		writeResponse(w, http.StatusBadRequest, nil, missingIDParameter)
		return
	}

	port, err := h.portService.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPortNotFound) {
			writeResponse(w, http.StatusNotFound, nil, err)
			return
		}
		writeResponse(w, http.StatusInternalServerError, nil, internalServer)
		return
	}

	writeResponse(w, http.StatusOK, port, nil)
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err != nil {
		encoderError := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if encoderError != nil {
			return
		}

		return
	}

	if data != nil {
		encoderError := json.NewEncoder(w).Encode(data)
		if encoderError != nil {
			return
		}
	}
}
