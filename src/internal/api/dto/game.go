package dto

type GameRequest struct {
	Board [][]string `json:"board"`
}

type GameResponse struct {
	ID      string     `json:"id"`
	Board   [][]string `json:"board"`
	Message string     `json:"message"`
}
