package domain

// import "github.com/google/uuid"

type GameService interface {
	AIMove(game *Game, computer Cell) (err error)

	ValidateBoard(oldBoard, newBoard *Board) error

	CheckGameOver(board Board) (isOver bool, winner Cell)

	ProcessGame(game *Game) (*Game, error)

	NewGame() (string, error)
}

type GameRepository interface {
	SaveGame(game *Game)
	GetGame(id string) (*Game, error)
}
