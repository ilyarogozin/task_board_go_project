package usecase

import (
	"context"

	"github.com/google/uuid"

	"notification_service/internal/domain"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *domain.Notification) error
}

type NotificationSender interface {
	Send(ctx context.Context, userID uuid.UUID, message string) error
}

type NotificationUseCase struct {
	repo   NotificationRepository
	sender NotificationSender
}

func NewNotificationUseCase(
	repo NotificationRepository,
	sender NotificationSender,
) *NotificationUseCase {
	return &NotificationUseCase{
		repo:   repo,
		sender: sender,
	}
}

func (uc *NotificationUseCase) HandleBoardCreated(
	ctx context.Context,
	userID uuid.UUID,
	boardTitle string,
) error {

	message := "Доска создана: " + boardTitle

	notification := domain.NewNotification(userID, message)

	if err := uc.sender.Send(ctx, userID, message); err != nil {
		return err
	}

	if err := uc.repo.Save(ctx, notification); err != nil {
		return err
	}

	return nil
}