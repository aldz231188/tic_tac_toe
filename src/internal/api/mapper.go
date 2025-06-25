package api

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"t03/internal/api/dto"
	"t03/internal/domain"
)

func ToDomainGame(pasedId string, s [][]string) (*domain.Game, error) {
	var game domain.Game
	id, err := uuid.Parse(pasedId)
	if err != nil {
		return &game, err
	}
	game.ID = id

	if len(s) != 3 {
		return &game, errors.New("board must have 3 rows")
	}

	for i := 0; i < 3; i++ {
		if len(s[i]) != 3 {
			return &game, errors.New("each row must have 3 columns")
		}

		for j := 0; j < 3; j++ {
			switch s[i][j] {
			case "X":
				game.Board[i][j] = domain.X
			case "O":
				game.Board[i][j] = domain.O
			case "":
				game.Board[i][j] = domain.Empty
			default:
				return &game, fmt.Errorf("invalid cell value at [%d][%d]: %s", i, j, s[i][j])
			}
		}
	}

	return &game, nil
}

func ToGameResponse(game *domain.Game) dto.GameResponse {
	board := make([][]string, 3)
	for i := range board {
		board[i] = make([]string, 3)
		for j := 0; j < 3; j++ {
			switch game.Board[i][j] {
			case domain.X:
				board[i][j] = "X"
			case domain.O:
				board[i][j] = "O"
			default:
				board[i][j] = ""
			}
		}
	}

	return dto.GameResponse{
		ID:    game.ID.String(),
		Board: board,
	}
}
