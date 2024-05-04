package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var (
	tracer = otel.Tracer("cook-handler")
)

func cookHandler(w http.ResponseWriter, r *http.Request) {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	_, span := tracer.Start(ctx, "cookHandler")
	defer span.End()

	if r.Method != http.MethodPost {
		errMessage := "Method not allowed"
		setErrorResponse(w, span, errMessage, http.StatusMethodNotAllowed)
		return
	}

	order, err := decodeOrderRequest(r)
	if err != nil {
		errMessage := "Invalid request body"
		setErrorResponse(w, span, errMessage, http.StatusBadRequest)
		return
	}

	cooking(ctx, order)
	successMessage := fmt.Sprintf("done cooking: %s", order.Food)
	log.Println(successMessage)
	setSuccessResponse(w, span, successMessage, http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
