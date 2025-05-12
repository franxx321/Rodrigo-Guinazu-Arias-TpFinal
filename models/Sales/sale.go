package Sales

import "time"

var Pending string = "Pending"
var Aproved string = "Aproved"
var Rejected string = "Rejected"

type Sale struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Amount    float32   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}
