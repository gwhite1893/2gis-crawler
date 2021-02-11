package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gwhite1893/2gis-crawler/config"
	"github.com/gwhite1893/2gis-crawler/internal/app"
)

// @title crawler
// @version 1.0.0
// @description 2-gis-crawler

// @host 127.0.0.1
// @BasePath /api/crawler/v1
func main() {
	configPath := flag.String("c", "config.yaml", "path to config")

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Println(err.Error())

		return
	}

	newApp, err := app.NewApp(
		app.WithHTTPServer(cfg),
	)
	if err != nil {
		return
	}

	defer newApp.Shutdown()

	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-ch
}
