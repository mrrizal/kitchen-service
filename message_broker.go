package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

func declareQueue() error {
	conn, err := amqp.Dial(BrokerURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func initRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial(BrokerURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeRabbitMQ(conn *amqp.Connection) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

type AMQPHeaderCarrier struct {
	Headers amqp.Table
}

func (a AMQPHeaderCarrier) Set(key, value string) {
	a.Headers[key] = value
}

func (a AMQPHeaderCarrier) Get(key string) string {
	v, ok := a.Headers[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (a AMQPHeaderCarrier) Keys() []string {
	i := 0
	r := make([]string, len(a.Headers))

	for k := range a.Headers {
		r[i] = k
		i++
	}

	return r
}

func getMessages(ctx context.Context, ch *amqp.Channel) {
	msgs, err := ch.Consume(
		QueueName, // Queue name
		"",        // Consumer name
		true,      // Auto-ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Go routine to consume messages
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	go func() {
		for msg := range msgs {
			ctx = propagator.Extract(ctx, &AMQPHeaderCarrier{msg.Headers})
			_, span := tracer.Start(ctx, "getMessages")
			defer span.End()
			var order Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				if err != nil {
					span.SetStatus(codes.Error, err.Error())
					span.RecordError(err)
				}
			}
			cooking(ctx, order)
			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	// Wait for termination signal
	<-sigs
}
