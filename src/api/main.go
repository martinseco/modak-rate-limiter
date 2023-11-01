package main

import (
	"fmt"
	"log"

	"github.com/martinseco/modak-rate-limiter/src/api/application"
)

func main() {
	app := application.BuildApplication()
	router := application.BuildRouter(app)
	server := application.BuildServer(router)

	fmt.Println("Serving application...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
