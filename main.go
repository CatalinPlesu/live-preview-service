package main

import (
	"context"
	"fmt"
	"os"
	"log"
	"os/signal"

	"github.com/CatalinPlesu/live-preview-service/application"
)

func main() {
	log.SetOutput(os.Stdout) // Make sure logs go to standard output

	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
