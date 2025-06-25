package domain

import "github.com/google/uuid"

type Status int

const (
	InProgress Status = iota
	PlayerWon
	AIWon
	Draw
)

type Cell int

const (
	Empty Cell = iota
	X
	O
)

type Board [3][3]Cell

type Game struct {
	ID     uuid.UUID
	Board  Board
	Status Status
}
