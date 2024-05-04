package main

import (
	"context"
	"log"
)

func cooking(ctx context.Context, order Order) error {
	_, span := tracer.Start(ctx, "cooking")
	defer span.End()

	// myErr := errors.New("failed cooking")
	// span.SetStatus(codes.Error, myErr.Error())
	// body, err := json.Marshal(order)
	// if err != nil {
	// 	span.SetStatus(codes.Error, myErr.Error())
	// 	span.RecordError(err)
	// 	return err
	// }
	// span.SetAttributes(
	// 	attribute.String("err.message", myErr.Error()),
	// 	attribute.String("message.payload", string(body)),
	// )
	// span.RecordError(myErr)

	log.Printf("cooking: %s\n", order.Food)
	randomSleep()
	return nil
}
