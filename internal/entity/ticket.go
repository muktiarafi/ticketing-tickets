package entity

type Ticket struct {
	ID      int64
	Title   string
	Price   float64
	OrderID int64
	UserID  int64
	Version int64
}
