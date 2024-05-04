package main

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func setErrorResponse(w http.ResponseWriter, span trace.Span, errorMessage string, statusCode int) {
	http.Error(w, errorMessage, http.StatusBadRequest)
	span.SetStatus(codes.Error, errorMessage)
	span.RecordError(errors.New(errorMessage))
}

func setSuccessResponse(w http.ResponseWriter, span trace.Span, successMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	span.SetStatus(codes.Ok, successMessage)
}
