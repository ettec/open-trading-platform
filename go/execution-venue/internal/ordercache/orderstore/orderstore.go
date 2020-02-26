package orderstore

import (
	"github.com/ettec/open-trading-platform/go/model"
)

type OrderStore interface {
	Write(order *model.Order) error
	Close()
}
