package helpers

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

// StartServerWithGracefulShutdown is starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App) {
	idleConnectionsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals
		<-sigint

		// Received an interrupt signal, shutdown
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server cannot be shutted down! Reason: %v", err)
		}

		close(idleConnectionsClosed)
	}()

	// Run server
	if err := a.Listen(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))); err != nil {
		log.Printf("Oops... Server cannot be started! Reason: %v", err)
	}

	<-idleConnectionsClosed
}
