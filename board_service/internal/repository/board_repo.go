package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"board_service/internal/domain"
)

type BoardRepository interface {
	Create(ctx context.Context, tx pgx.Tx, b domain.Board) error
	Get(ctx context.Context, id string) (*domain.Board, error)
}

type boardRepo struct {
	db *pgx.Conn
}

func NewBoardRepo(db *pgx.Conn) BoardRepository {
	return &boardRepo{db}
}

func (r *boardRepo) Create(ctx context.Context, tx pgx.Tx, b domain.Board) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO boards (id, title, description, owner_id)
		 VALUES ($1,$2,$3,$4)`,
		b.ID, b.Title, b.Description, b.OwnerID,
	)
	return err
}

func (r *boardRepo) Get(ctx context.Context, id string) (*domain.Board, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id,title,description,owner_id FROM boards WHERE id=$1`, id)

	var b domain.Board
	err := row.Scan(&b.ID, &b.Title, &b.Description, &b.OwnerID)
	return &b, err
}