package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"jayant/database"
	"jayant/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTimeOut = 10 * time.Second

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.SetupRouter()
	if err := database.ConnectAndMigrate(
		"localhost",
		"5454",
		"my_todo",
		"local",
		"local",
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: #{err}")
	}
	logrus.Print("migration successfully")
	go func() {
		if err := http.ListenAndServe(":8080", srv); err != nil {
			fmt.Println("Error:", err)
		}
	}()
	logrus.Print("Server started at :8080")

	<-done

	logrus.Info("shutting down server")
	if err := database.ShutdownDatabase(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}
	srv = http.TimeoutHandler(srv, shutDownTimeOut, "failed to gracefully shutdown server")
}
