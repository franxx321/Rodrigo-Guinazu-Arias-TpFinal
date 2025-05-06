package sales

import "time"

var pending string = "pending"
var aproved string = "aproved"
var rejected string = "rejected"

type Sale struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Amount    float32   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}
