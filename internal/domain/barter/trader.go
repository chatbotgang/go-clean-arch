package barter

import "time"

type Trader struct {
	ID        int
	UID       string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTrader(uid string, email string, name string) Trader {
	return Trader{
		UID:   uid,
		Email: email,
		Name:  name,
	}
}
