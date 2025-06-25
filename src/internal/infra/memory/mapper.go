package memory

import (
	"github.com/google/uuid"
	"t03/internal/domain"
)

func toEntity(game *domain.Game) *GameEntity {
	board := [3][3]int{}
	for i := range game.Board {
		for j := range game.Board[i] {
			board[i][j] = int(game.Board[i][j])
		}
	}
	return &GameEntity{
		ID:     game.ID.String(),
		Board:  board,
		Status: int(game.Status),
	}
}

func toDomain(entity *GameEntity) (*domain.Game, error) {
	id, err := uuid.Parse(entity.ID)
	if err != nil {
		return nil, err
	}

	board := domain.Board{}
	for i := range entity.Board {
		for j := range entity.Board[i] {
			board[i][j] = domain.Cell(entity.Board[i][j])
		}
	}

	return &domain.Game{
		ID:     id,
		Board:  board,
		Status: domain.Status(entity.Status),
	}, nil
}
