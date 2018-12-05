package model

import "activityjs.io/serve/errors"

// OrderItem 订单项
type OrderItem struct {
	Effect
	ID          *Identity `header:"ID"`
	Product     Product
	ProductID   *Identity
	Title       string  `header:"title"`
	Price       float64 `header:"price"`
	Quantity    int     `header:"quantity"`
	Subtotal    float64 `header:"subtotal"`
	ProduceDate Time
	upsides     []propState

	propagation    bool
	disableCompute bool
	// parent      *Order
}

type propState struct {
	Type string
	Val  interface{}
}

func (item *OrderItem) Get(key string) (interface{}, bool) {

	switch key {
	case "ID":
		return item.ID, true
	case "Price":
		return item.Price, true
	case "ProduceDate":
		return item.ProduceDate, true
	case "Product":
		return item.Product, true
	case "Quantity":
		return item.Quantity, true
	case "Subtotal":
		return item.Subtotal, true
	case "Title":
		return item.Title, true
	default:
		return nil, false
	}
}

func (item *OrderItem) Set(key string, val interface{}) error {
	switch key {
	case "ID":
		item.ID = val.(*Identity)
	case "Price":
		item.compute(key, item.Price, func() interface{} {
			item.Price = val.(float64)
			return item.Price
		})
	case "ProduceDate":
		item.ProduceDate = val.(Time)
	case "Product":
		item.Product = val.(Product)
	case "Quantity":
		item.compute(key, item.Quantity, func() interface{} {
			item.Quantity = val.(int)
			return item.Quantity
		})
		// item.compute(key, item.Quantity)
		// item.Quantity = val.(int)
	case "Subtotal":
		item.compute(key, item.Subtotal, func() interface{} {
			item.Subtotal = val.(float64)
			return item.Subtotal
		})
		// oldSubtotal := item.Subtotal
		// item.orderCompute("diff.item.subtotal", oldSubtotal-item.Subtotal)
		// item.Subtotal = val.(float64)
	case "Title":
		item.Title = val.(string)
	default:
		return errors.New("无效的 Item 字段 " + key)
	}
	return nil
}

// func (item *OrderItem) SetOrder(order *Order) {
// 	item.parent = order
// }

func (item *OrderItem) compute(key string, oldVal interface{}, fn func() interface{}) {
	if disableCompute {
		fn()
		return
	}
	switch key {
	case "Price":
		oldSubtotal := item.Subtotal
		if newPrice, ok := fn().(float64); ok {
			item.Subtotal = float64(item.Quantity) * newPrice
			item.propCompute("diff.item.subtotal", oldSubtotal-item.Subtotal)
		}

	case "Quantity":
		oldSubtotal := item.Subtotal
		oldQuantity := item.Quantity
		if newQuantity, ok := fn().(int); ok {
			item.Subtotal = float64(newQuantity) * item.Price
			item.propCompute("diff.item.subtotal", oldSubtotal-item.Subtotal)
			item.propCompute("diff.item.quantity", oldQuantity-newQuantity)
		}
	case "Subtotal":
		oldSubtotal := item.Subtotal
		if newSubtotal, ok := fn().(float64); ok {
			item.propCompute("diff.item.subtotal", oldSubtotal-newSubtotal)
		}
	case "all":
		item.Subtotal = float64(item.Quantity) * item.Price
	}
}

func (item *OrderItem) propCompute(key string, newVal interface{}) {
	if !upsidePropagation {
		return
	}

	item.upsides = append(item.upsides, propState{key, newVal})
}

func (item *OrderItem) WithEffect(key, field, summary string, val interface{}, fn func() error) error {
	item.AddEffect(key, field, summary, val)
	return fn()
}

func (item *OrderItem) DisabledCompute() {
	item.disableCompute = true
}

func Propagation(order *Order, item *OrderItem) {
	if item.propagation {
		for _, state := range item.upsides {
			order.compute(state.Type, state.Val)
		}
	}

	item.upsides = nil
}
