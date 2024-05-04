package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func decodeOrderRequest(r *http.Request) (Order, error) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	randFloat := rand.Float64()
	min := 5.0  // 300 milliseconds
	max := 20.0 // 500 milliseconds
	randomDuration := time.Duration(randFloat*(max-min)+min) * time.Millisecond
	time.Sleep(time.Duration(randomDuration))
}
