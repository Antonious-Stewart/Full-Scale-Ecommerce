package shared

import (
	"github.com/google/uuid"
	"time"
)

type AuthEntity struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tokens       []uint8
	Password     []byte
	OrderHistory []uint8
}
