package game

type Repository interface {
	FindById(id uint64) (*Game, error)
	FindAll() ([]Game, error)
	Create(g *Game) error
}
