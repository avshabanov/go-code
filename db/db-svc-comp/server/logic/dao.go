package logic

import (
	"io"
)

type UserProfile struct {
	Name string
}

type LineItem struct {
	ID int
}

type Order struct {
	ID        string
	LineItems []*LineItem
}

type OrderPage struct {
	Items       []*Order
	OffsetToken string
}

type Dao interface {
	io.Closer
	Add(p *UserProfile) (string, error)
	QueryOrders(userID string, from string, limit int) (*OrderPage, error)
}
