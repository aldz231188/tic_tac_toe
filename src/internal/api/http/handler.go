package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"t03/internal/api"
	"t03/internal/api/dto"
	"t03/internal/domain"
)

type GameHandler struct {
	Service domain.GameService
}

func NewGameHandler(service domain.GameService) *GameHandler {
	return &GameHandler{
		Service: service,
	}
}

func (h *GameHandler) HandleNewGame(w http.ResponseWriter, r *http.Request) {

	id, err := h.Service.NewGame()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *GameHandler) HandleGameMove(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/game/")

	var gameReq dto.GameRequest
	if err := json.NewDecoder(r.Body).Decode(&gameReq); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	playerMove, err := api.ToDomainGame(id, gameReq.Board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerMove, err = h.Service.ProcessGame(playerMove)

	response := api.ToGameResponse(playerMove)
	if err != nil {
		response.Message = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
