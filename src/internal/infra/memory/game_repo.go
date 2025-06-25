package memory

import (
	"errors"
	"t03/internal/domain"
)

type GameRepositoryImpl struct {
	storage *Storage
}

func NewGameRepository(storage *Storage) domain.GameRepository {
	return &GameRepositoryImpl{storage: storage}
}

func (repo *GameRepositoryImpl) SaveGame(game *domain.Game) {
	entity := toEntity(game)
	repo.storage.games.Store(entity.ID, entity)
}

func (repo *GameRepositoryImpl) GetGame(id string) (*domain.Game, error) {
	value, ok := repo.storage.games.Load(id)
	if !ok {
		return nil, errors.New("game not found")
	}

	entity, ok := value.(*GameEntity)
	if !ok {
		return nil, errors.New("invalid game data")
	}

	return toDomain(entity)
}
