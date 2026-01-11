package board

type BoardCreatedEvent struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerId     string    `json:"owner_id"`
}