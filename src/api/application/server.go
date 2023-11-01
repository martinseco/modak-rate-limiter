package application

import (
	"fmt"
	"net/http"
	"time"
)

const (
	portNumber = ":8080"
)

func BuildServer(router http.Handler) *http.Server {
	fmt.Println("Building server...")
	return &http.Server{
		Addr:         portNumber,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
