package app

import (
	"context"
	"github.com/rs/zerolog/log"
)

func Run(_ context.Context, port string, adminPort string) {
	log.Info().Msgf("started on ports %s and %s", port, adminPort)
}
