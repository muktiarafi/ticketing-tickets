package entity

type Ticket struct {
	ID      int64   `json:"id"`
	Title   string  `json:"title"`
	Price   float64 `json:"price"`
	OrderID int64   `json:"orderId"`
	UserID  int64   `json:"userId"`
	Version int64   `json:"version"`
}
