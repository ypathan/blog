package main

import (
	"log"
	"net/http"
	"os"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}


func main() {
	//routes
    http.HandleFunc("/", handleRoot)

	//logging
	file , err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err.Error())
	}
	log.SetOutput(file)

    // start sercer
    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Println("HTTP server error:", err)
    }
}
