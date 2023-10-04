package main

import (
	"log"
	"net/http"
	"os"

	console "github.com/itzsBananas/mc-server-monitor/internal/data"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	remoteConsole console.ConsoleInterface
}

func main() {
	serverAddress := getEnv("SERVER_ADDRESS", ":4000")
	rconAddress := getEnv("RCON_ADDRESS", "127.0.0.1:25575")
	rconPassword := getEnv("RCON_PASSWORD", "password")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	con, err := console.Open(rconAddress, rconPassword)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer con.Close()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		remoteConsole: con,
	}

	srv := &http.Server{
		Addr:     serverAddress,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", serverAddress)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
