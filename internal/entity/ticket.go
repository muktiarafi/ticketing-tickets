package entity

type Ticket struct {
	ID      int     `json:"id" msgpack:"id"`
	Title   string  `json:"title" msgpack:"title"`
	Price   float64 `json:"price" msgpack:"price"`
	OrderID int     `json:"orderId" msgpack:"orderId"`
	UserID  int     `json:"userId" msgpack:"userId"`
	Version int     `json:"version" msgpack:"version"`
}
