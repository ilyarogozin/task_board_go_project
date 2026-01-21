package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	board "github.com/ilyarogozin/task_board_go_project/gen/go/board"
)

type BoardClient struct {
	conn   *grpc.ClientConn
	client board.BoardServiceClient
}

func NewBoardClient(addr string) (*BoardClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &BoardClient{
		conn:   conn,
		client: board.NewBoardServiceClient(conn),
	}, nil
}

func (c *BoardClient) CreateBoard(
	ctx context.Context,
	req *board.CreateBoardRequest,
) (*board.BoardResponse, error) {
	return c.client.CreateBoard(ctx, req)
}

func (c *BoardClient) Close() error {
	return c.conn.Close()
}