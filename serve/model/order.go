package model

import (
	"activityjs.io/serve/errors"
)

// Order 订单
type Order struct {
	Items []OrderItem // 订单列表
	Gifts []Gift
	Count int     `header:"count"`
	Total float64 `header:"total"`
}

func (order *Order) Get(key string) (interface{}, bool) {
	switch key {
	case "Count":
		return order.Count, true
	case "Total":
		return order.Total, true
	default:
		return nil, false
	}
}

func (order *Order) Set(key string, val interface{}) error {
	switch key {
	case "Count":
		order.Count = val.(int)
	case "Total":
		order.Total = val.(float64)
	default:
		return errors.New("无效的 User 字段 " + key)
	}
	return nil
}

func (order *Order) AppendItem(item *OrderItem) error {
	order.Items = append(order.Items, *item)
	// item.parent = order

	order.compute("additem", item)
	return nil
}

func (order *Order) RemoveItem(item *OrderItem) error {
	var ok bool
	if order.Items, ok = RemoveItem(order.Items, item); !ok {
		return errors.New("没有找到这个元素")
	}
	// item.parent = nil
	order.compute("removeitem", item)
	return nil
}

// compute Order 的计算属性比较特殊，我们主要使用 key 做为事件名，
// 其中包括 diff.item.quantity, additem, removeitem 等
func (order *Order) compute(key string, val interface{}) error {
	if disableCompute {
		return nil
	}

	switch key {
	case "diff.item.subtotal": // 更改某 item 的小计，val 是 subtotal 的差价
		order.Total -= val.(float64)
	case "diff.item.quantity": // 更改某 item 的数据，val 是 quantity 的差价
		order.Count -= val.(int)
	case "additem", "removeitem":
		if item, ok := val.(*OrderItem); !ok {
			return errors.New("无效的  additem 参数")
		} else {
			// subtotal := item.Price * float64(item.Quantity)
			order.Total += item.Subtotal
			order.Count += item.Quantity
		}
	case "all":
		var (
			count int
			sum   float64
		)

		for _, item := range order.Items {
			// item.compute("all", item)
			count += item.Quantity
			sum += item.Subtotal
		}
		// order.Items
		// item.Subtotal = float64(item.Quantity) * newPrice
		// item.orderCompute("all", item)
	default:
		return errors.New("无效的 compute 消息 id" + key)
	}

	return nil
}

func RemoveItem(items []OrderItem, ite *OrderItem) ([]OrderItem, bool) {

	var idx = -1
	for i, m := range items {
		if m.ID.Compare(ite.ID) {
			idx = i
			break
		}
	}

	if idx >= 0 {
		if idx == len(items)-1 {
			return items[:idx], true
		} else {
			return append(items[:idx], items[idx+1:]...), true
		}
	} else {
		return items, false
	}
}
