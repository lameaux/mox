package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lameaux/mox/internal/admin"
	"github.com/lameaux/mox/internal/banner"
	"github.com/lameaux/mox/internal/mock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	appName    = "mox"
	appVersion = "v0.0.1"

	defaultPortMocks = 8080
	defaultPortAdmin = 9090

	stopTimeout = 5 * time.Second
)

var GitHash string //nolint:gochecknoglobals

func main() {
	var (
		debug      = flag.Bool("debug", false, "enable debug mode")
		logJSON    = flag.Bool("logJson", false, "log as json")
		accessLog  = flag.Bool("accessLog", false, "enable access log")
		skipBanner = flag.Bool("skipBanner", false, "skip banner")
		port       = flag.Int("port", defaultPortMocks, "port for mock server")
		adminPort  = flag.Int("adminPort", defaultPortAdmin, "port for admin server")
		configPath = flag.String("configPath", "", "path to config (mappings, templates, files)")
	)

	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if !*logJSON {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if !*skipBanner {
		fmt.Print(banner.Banner) //nolint:forbidigo
	}

	log.Info().Str("version", appVersion).Str("build", GitHash).Msg(appName)

	h, err := mock.NewHandler(*configPath, *accessLog)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize mock handler")
	}

	mockServer := mock.StartServer(*port, h)
	adminServer := admin.StartServer(*adminPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handleSignals(func() {
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
	timedOutCtx, cancel := context.WithTimeout(ctx, stopTimeout)
	defer cancel()

	if err := server.Shutdown(timedOutCtx); err != nil {
		log.Error().Str("addr", server.Addr).Err(err).Msg("failed to shutdown server")
	}

	log.Debug().Str("addr", server.Addr).Msg("server stopped")
}
