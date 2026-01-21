package sender

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (s *EmailSender) Send(
	_ context.Context,
	userID uuid.UUID,
	message string,
) error {

	log.Info().
		Str("user_id", userID.String()).
		Msg("Email отправлен пользователю: " + message)

	return nil
}