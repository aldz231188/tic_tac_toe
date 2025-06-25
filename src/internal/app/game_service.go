package app

import (
	"errors"
	"github.com/google/uuid"
	"math"
	"t03/internal/domain"
)

type GameServiceImpl struct {
	repo domain.GameRepository
}

func NewGameService(repo domain.GameRepository) domain.GameService {
	return &GameServiceImpl{repo: repo}
}

func (svc *GameServiceImpl) NewGame() (string, error) {
	id := uuid.New()
	game := &domain.Game{
		ID:     id,
		Board:  domain.Board{},
		Status: domain.InProgress,
	}
	// svc.AIMove(game, domain.O)

	svc.repo.SaveGame(game)

	return id.String(), nil
}

func (svc *GameServiceImpl) ProcessGame(afterPlayerMove *domain.Game) (*domain.Game, error) {
	beforePlayerMove, err := svc.repo.GetGame(afterPlayerMove.ID.String())
	if err != nil {
		return beforePlayerMove, err
	}

	if beforePlayerMove.Status != domain.InProgress {
		return beforePlayerMove, errors.New("session finished")
	}

	err = svc.ValidateBoard(&beforePlayerMove.Board, &afterPlayerMove.Board)
	if err != nil {
		return beforePlayerMove, err
	}

	over, who := svc.CheckGameOver(afterPlayerMove.Board)

	if over && who != domain.Empty {
		afterPlayerMove.Status = domain.PlayerWon
		svc.repo.SaveGame(afterPlayerMove)
		return afterPlayerMove, errors.New("you win")
	} else if over && who == domain.Empty {
		afterPlayerMove.Status = domain.Draw
		svc.repo.SaveGame(afterPlayerMove)
		return afterPlayerMove, errors.New("played to a draw")
	}

	err = svc.AIMove(afterPlayerMove, domain.O)
	if err != nil {
		return afterPlayerMove, err
	}

	over, who = svc.CheckGameOver(afterPlayerMove.Board)

	if over && who != domain.Empty {
		afterPlayerMove.Status = domain.AIWon
		svc.repo.SaveGame(afterPlayerMove)
		return afterPlayerMove, errors.New("you lose")
	} else if over && who == domain.Empty {
		afterPlayerMove.Status = domain.Draw
		svc.repo.SaveGame(afterPlayerMove)
		return afterPlayerMove, errors.New("played to a draw")
	}
	svc.repo.SaveGame(afterPlayerMove)

	return afterPlayerMove, nil
}

func (svc *GameServiceImpl) AIMove(game *domain.Game, computer domain.Cell) error {
	bestScore := math.MinInt
	bestMove := [2]int{-1, -1}

	board := game.Board

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				board[i][j] = computer
				score := svc.minimax(board, 0, false, computer)
				board[i][j] = domain.Empty

				if score > bestScore {
					bestScore = score
					bestMove = [2]int{i, j}
				}
			}
		}
	}

	if bestMove[0] == -1 {
		return errors.New("no moves left")
	}
	game.Board[bestMove[0]][bestMove[1]] = domain.O

	return nil
}

func (svc *GameServiceImpl) ValidateBoard(oldBoard, newBoard *domain.Board) error {
	moveCount := 0

	for i := range oldBoard {
		for j := range oldBoard[i] {
			oldCell := oldBoard[i][j]
			newCell := newBoard[i][j]

			switch {
			case oldCell == newCell:
				continue

			case oldCell == domain.Empty && newCell == domain.X:
				moveCount++

			default:
				return errors.New("board is corrupted")
			}
		}
	}

	if moveCount == 0 {
		return errors.New("your turn")
	}
	if moveCount > 1 {
		return errors.New("only one move is allowed at a time")
	}

	return nil
}

func (svc *GameServiceImpl) CheckGameOver(board domain.Board) (bool, domain.Cell) {
	for i := 0; i < 3; i++ {
		if board[i][1] != domain.Empty && (board[i][1] == board[i][0] && board[i][1] == board[i][2]) {
			return true, board[i][1]
		} else if board[1][i] != domain.Empty && (board[1][i] == board[0][i] && board[1][i] == board[2][i]) {
			return true, board[1][i]
		} else {
			continue
		}
	}

	if board[1][1] != domain.Empty && ((board[0][0] == board[1][1] && board[1][1] == board[2][2]) || (board[0][2] == board[1][1] && board[1][1] == board[2][0])) {
		return true, board[1][1]
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				return false, domain.Empty
			}

		}

	}

	return true, domain.Empty
}

func (svc *GameServiceImpl) minimax(board domain.Board, depth int, isMaximizing bool, computer domain.Cell) int {
	isOver, winner := svc.CheckGameOver(board)
	if isOver {
		if winner == computer {
			return 10 - depth // победа ИИ
		} else if winner != domain.Empty {
			return depth - 10 // победа игрока
		}
		return 0 // ничья
	}

	if isMaximizing {
		best := math.MinInt
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.Empty {
					board[i][j] = computer
					score := svc.minimax(board, depth+1, false, computer)
					board[i][j] = domain.Empty
					best = max(best, score)
				}
			}
		}
		return best
	} else {
		best := math.MaxInt
		player := domain.X
		if computer == domain.X {
			player = domain.O
		}

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.Empty {
					board[i][j] = player
					score := svc.minimax(board, depth+1, true, computer)
					board[i][j] = domain.Empty
					best = min(best, score)
				}
			}
		}
		return best
	}
}
