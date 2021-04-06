package orderbook

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Order strores information about request
type Order struct {
	side      Side
	id        string
	timestamp time.Time
	quantity  decimal.Decimal
	price     decimal.Decimal
}

// MarketView represents order book in a glance
type MarketView struct {
	Asks map[string]decimal.Decimal `json:"asks"`
	Bids map[string]decimal.Decimal `json:"bids"`
}

// NewOrder creates new constant object Order
func NewOrder(orderID string, side Side, quantity, price decimal.Decimal, timestamp time.Time) *Order {
	return &Order{
		id:        orderID,
		side:      side,
		quantity:  quantity,
		price:     price,
		timestamp: timestamp,
	}
}

// ID returns orderID field copy
func (o *Order) ID() string {
	return o.id
}

// Side returns side of the order
func (o *Order) Side() Side {
	return o.side
}

// Quantity returns quantity field copy
func (o *Order) Quantity() decimal.Decimal {
	return o.quantity
}

// Price returns price field copy
func (o *Order) Price() decimal.Decimal {
	return o.price
}

// Time returns timestamp field copy
func (o *Order) Time() time.Time {
	return o.timestamp
}

// String implements Stringer interface
func (o *Order) String() string {
	return fmt.Sprintf("\n\"%s\":\n\tside: %s\n\tquantity: %s\n\tprice: %s\n\ttime: %s\n", o.ID(), o.Side(), o.Quantity(), o.Price(), o.Time())
}

// MarshalJSON implements json.Marshaler interface
func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			S         Side            `json:"side"`
			ID        string          `json:"id"`
			Timestamp time.Time       `json:"timestamp"`
			Quantity  decimal.Decimal `json:"quantity"`
			Price     decimal.Decimal `json:"price"`
		}{
			S:         o.Side(),
			ID:        o.ID(),
			Timestamp: o.Time(),
			Quantity:  o.Quantity(),
			Price:     o.Price(),
		},
	)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (o *Order) UnmarshalJSON(data []byte) error {
	obj := struct {
		S         Side            `json:"side"`
		ID        string          `json:"id"`
		Timestamp time.Time       `json:"timestamp"`
		Quantity  decimal.Decimal `json:"quantity"`
		Price     decimal.Decimal `json:"price"`
	}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	o.side = obj.S
	o.id = obj.ID
	o.timestamp = obj.Timestamp
	o.quantity = obj.Quantity
	o.price = obj.Price
	return nil
}

// GetOrderSide gets the orderside along with its orders in one side of the market
func (ob *OrderBook) GetOrderSide(side Side) *OrderSide {
	switch side {
	case Buy:
		return ob.bids
	default:
		return ob.asks
	}
}

// MarketOverview gives an overview of the market including the quantities and prices of each side in the market
// asks:   qty   price       bids:  qty   price
//         0.2   14                 0.9   13
//         0.1   14.5               5     14
//         0.8   16                 2     16
func (ob *OrderBook) MarketOverview() *MarketView {
	return &MarketView{
		Asks: compileOrders(ob.asks.Orders()),
		Bids: compileOrders(ob.bids.Orders()),
	}
}

// compileOrders compiles orders in the following format
func compileOrders(orders []*list.Element) map[string]decimal.Decimal {

	q := make(map[string]decimal.Decimal)

	for _, o := range orders {
		if order, ok := o.Value.(*Order); ok {
			if qty, ok := q[order.Price().String()]; ok {
				qty = qty.Add(order.Quantity())
				q[order.Price().String()] = qty
			} else {
				q[order.Price().String()] = order.Quantity()
			}

		}
	}

	return q
}
