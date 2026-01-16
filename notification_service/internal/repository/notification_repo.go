package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"notification_service/internal/domain"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

func (r *NotificationRepository) Save(
	ctx context.Context,
	n *domain.Notification,
) error {

	query := `
		INSERT INTO notifications (id, user_id, message, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		n.ID,
		n.UserID,
		n.Message,
		n.CreatedAt,
	)

	return err
}