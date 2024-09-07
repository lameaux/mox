package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lameaux/mox/internal/admin"
	"github.com/lameaux/mox/internal/metrics"
	"github.com/lameaux/mox/internal/mock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	var logJson = flag.Bool("logJson", false, "log as json")
	var accessLog = flag.Bool("accessLog", false, "enable access log")
	var skipBanner = flag.Bool("skipBanner", false, "skip banner")
	var port = flag.Int("port", 8080, "port for mock server")
	var adminPort = flag.Int("adminPort", 8181, "port for admin server")
	var metricsPort = flag.Int("metricsPort", 9090, "port for metrics server")
	var configPath = flag.String("configPath", "./config", "path to config (mappings, templates, files)")

	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if !*logJson {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if !*skipBanner {
		fmt.Print(logo)
	}

	log.Info().Str("version", appVersion).Str("build", GitHash).Msg(appName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h, err := mock.NewHandler(*configPath, *accessLog)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize mock handler")
	}

	mockServer := mock.StartServer(*port, h)
	adminServer := admin.StartServer(*adminPort)
	metricsServer := metrics.StartServer(*metricsPort)

	handleSignals(func() {
		stopServer(ctx, metricsServer)
		stopServer(ctx, adminServer)
		stopServer(ctx, mockServer)
		cancel()
	})
}

func handleSignals(shutdownFn func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Info().Str("signal", sig.String()).Msgf("received signal")
	shutdownFn()
}

func stopServer(ctx context.Context, server *http.Server) {
	timedOutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timedOutCtx); err != nil {
		log.Error().Str("addr", server.Addr).Err(err).Msg("failed to shutdown server")
	}

	log.Debug().Str("addr", server.Addr).Msg("server stopped")
}
