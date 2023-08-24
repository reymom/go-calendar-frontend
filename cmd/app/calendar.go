package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/reymom/go-calendar-frontend/internal/config"
	"github.com/reymom/go-calendar-frontend/internal/routes"
	"github.com/rs/zerolog/log"
)

const ConfigJsonName = "calendar_app.json"

func main() {
	log.Info().Msgf(" ---------- Calendar App, Version %s, Build date %s ----------", config.Version, config.BuildDate)

	configFilePath := flag.String("c", path.Join("conf/", config.ConfigJsonName), "config file path")
	conf, e := config.GenerateConfig(*configFilePath)
	if e != nil {
		log.Err(e).Msgf("Error while generating config")
		os.Exit(1)
	}

	handler, e := routes.GenerateRoutes(conf)
	if e != nil {
		log.Err(e).Msg("Could not initialize route handlers")
		os.Exit(1)
	}

	appServer := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	log.Info().Msg("Server is started and listens on port 8080")
	go startServer(appServer)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	log.Trace().Msg("Shutting down..")

	e = appServer.Shutdown(context.Background())
	if e != nil {
		log.Err(e).Msg("Error while shutting down the app listener")
	}

	os.Exit(0)
}

func startServer(server *http.Server) {

	if err := server.ListenAndServe(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

}
