package cmd

import (
	"notification_service/internal/handler"
	"notification_service/internal/infra"
	"notification_service/internal/repository"
)

func main() {
	db := infra.MustPostgres()
	repo := repository.NewNotificationRepo(db)
	handler.StartConsumer(repo)
}