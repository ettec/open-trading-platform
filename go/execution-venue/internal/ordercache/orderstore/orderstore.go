package orderstore

import (
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/model"
)

type OrderStore interface {
	Write(order *model.Order) error
	Close()
}
