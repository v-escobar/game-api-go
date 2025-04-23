package dto

type Error struct {
	Message string `json:"message"`
}

type Game struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
