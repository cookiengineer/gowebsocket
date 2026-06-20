package main

import "example/chat"
import "log"
import "net/http"
import "os"

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	chat_server := chat.NewServer()

	http.Handle("/", chat_server)

	log.Printf("Chat Server listening on :%s", port)
	log.Printf("Open http://localhost:%s in your Web Browser", port)

	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
