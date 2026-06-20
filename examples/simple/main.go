package main

import gows "github.com/cookiengineer/gowebsocket"
import "log"
import "os"
import "time"

import "fmt"

func main() {

	// WebSocket Server Usage
	go func() {

		logger := log.New(os.Stdout, "[server] ", log.LstdFlags)
		server := &gows.Server{
			Addr:    ":8080",
			Handler: func(websocket *gows.WebSocket) {

				logger.Print("Client connected!")

				websocket.OnMessage = func(message []byte) {
					logger.Printf("Received message: %s", message)
				}

				websocket.OnClose = func() {
					logger.Print("Client disconnected!")
				}

			},
			TLSConfig: nil,
			ErrorLog:  logger,
		}

		server.Listen()

	}()

	// WebSocket Client Usage
	go func() {

		time.Sleep(100 * time.Millisecond)

		logger := log.New(os.Stdout, "[client] ", log.LstdFlags)
		client, err0 := gows.NewClient("ws://localhost:8080")

		if err0 == nil {

			err1 := client.Connect()

			if err1 == nil {

				logger.Print(client)

				time.Sleep(100 * time.Millisecond)
				err := client.Socket.Send([]byte("Hello, world!"))
				fmt.Println(err)

				time.Sleep(100 * time.Millisecond)
				client.Socket.Send([]byte("Second Hello, world!"))

				time.Sleep(1 * time.Second)
				client.Socket.Close(gows.StatusGoingAway, "Goodbye!")

			}

		} else {
			logger.Fatal(err0)
		}

	}()

	time.Sleep(10 * time.Second)

}
