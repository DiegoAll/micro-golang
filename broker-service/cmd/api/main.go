package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
}

// VUE STRUCTURE
// type config struct {
// 	port int
// }

// type application struct {
// 	config   config
// 	infoLog  *log.Logger
// 	errorLog *log.Logger
// 	// db       *driver.DB
// 	models data.Models
// 	//routes   *mux.Router //error ./cmd/api/
// 	environment string
// }

func main() {

	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
