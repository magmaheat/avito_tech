package entity

import (
	"github.com/google/uuid"
	"time"
)

type House struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	Year      int64     `json:"year"`
	Developer string    `json:"developer"`
	CreatedFl time.Time `json:"created_at"`
	UpdateFl  time.Time `json:"update_at"`
}

type Flat struct {
	ID      int64     `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	HouseID int64     `json:"house_id"`
	Number  int64     `json:"number"`
	Price   int64     `json:"price"`
	Rooms   int64     `json:"rooms"`
	Status  string    `json:"status"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	UserType string    `json:"user_type"`
}

type Subscription struct {
	HouseID int    `json:"house_id"`
	Email   string `json:"email"`
}
