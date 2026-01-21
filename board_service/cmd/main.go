package main

import (
	"github.com/rs/zerolog/log"

	"board_service/app"
)

func main() {
  if err := app.Run(); err != nil {
	log.Error().Err(err).Msg("failed to run app")
  }
}