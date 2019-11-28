package orderstore

import (
	"github.com/ettec/open-trading-platform/execution-venue/model"
)

type OrderStore interface {
	Write(order *model.Order) error
	Close()
}
