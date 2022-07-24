package barter

import "time"

type Good struct {
	ID        int
	Name      string
	OwnerID   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGood(trader Trader, goodName string) Good {
	return Good{
		Name:      goodName,
		OwnerID:   trader.ID,
		CreatedAt: time.Now(),
	}
}

func (d *Good) MyGood(trader Trader) bool {
	return d.OwnerID == trader.ID
}
