package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lameaux/mox/internal/app"
	"github.com/lameaux/mox/internal/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

const (
	logo = `
 █████████████    ██████  █████ █████
░░███░░███░░███  ███░░███░░███ ░░███ 
 ░███ ░███ ░███ ░███ ░███ ░░░█████░  
 ░███ ░███ ░███ ░███ ░███  ███░░░███ 
 █████░███ █████░░██████  █████ █████
░░░░░ ░░░ ░░░░░  ░░░░░░  ░░░░░ ░░░░░           
`
	appName    = "mox"
	appVersion = "v0.0.1"
)

var GitHash string

func main() {
	var debug = flag.Bool("debug", false, "enable debug mode")
	var skipBanner = flag.Bool("skipBanner", false, "skip banner")
	var port = flag.String("port", "8080", "port for server")
	var adminPort = flag.String("adminPort", "8081", "port for admin")
	var metricsPort = flag.String("metricsPort", "9090", "port for metrics")

	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if !*skipBanner {
		fmt.Print(logo)
	}

	log.Info().Str("version", appVersion).Str("build", GitHash).Msg(appName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := metrics.StartServer(*metricsPort)
	defer metrics.StopServer(ctx, server)

	handleSignals(func() {
		metrics.StopServer(ctx, server)
		cancel()
	})

	app.Run(ctx, *port, *adminPort)
}

func handleSignals(shutdownFn func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msgf("received signal")
		shutdownFn()
	}()
}
