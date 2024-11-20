package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CatalinPlesu/live-preview-service/messaging"
	"github.com/CatalinPlesu/live-preview-service/model"
	"github.com/gorilla/websocket"
)

type App struct {
	config   Config
	hub      *Hub
	rabbitMQ *messaging.RabbitMQ
	upgrader websocket.Upgrader
}

// New initializes the application with RabbitMQ connection
func New(config Config) *App {
	rabbitMQ, err := messaging.NewRabbitMQ(config.RabitMQURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	return &App{
		config:   config,
		hub:      newHub(),
		rabbitMQ: rabbitMQ,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Start sets up the WebSocket server and RabbitMQ consumer
func (a *App) Start(ctx context.Context) error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("New WebSocket connection request received.")
		serveWs(a.hub, w, r)
	})

	go a.hub.run()

	// Start consuming messages from RabbitMQ in a goroutine
	go a.startConsumingMessages()

	// Run the HTTP server in a separate goroutine
	go func() {
		if err := http.ListenAndServe(fmt.Sprint(":", a.config.ServerPort), nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	log.Println("Shutting down WebSocket server...")
	return nil
}

// Function to consume messages from RabbitMQ and forward them to WebSocket clients
func (a *App) startConsumingMessages() {
	log.Println("Starting RabbitMQ message consumer...")

	err := a.rabbitMQ.ConsumeMessages("message", func(msg model.Message) {
		messageContent, err := json.Marshal(model.MessageMin{ParentID: msg.ParentID,
			UserID: msg.UserID, MessageText: msg.MessageText})
		if err != nil {
			log.Println("Error marshalling message:", err)
			return
		}
		log.Printf("Broadcasting message to WebSocket clients: %s", messageContent)
		a.hub.broadcast <- messageContent
	})
	if err != nil {
		log.Printf("Error consuming RabbitMQ messages: %v", err)
	}
}
