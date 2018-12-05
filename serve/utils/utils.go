package utils

import (
	"activityjs.io/serve/model"
)

// FindItem 查找订单中的产品
func FindItem(items []model.OrderItem, ite *model.OrderItem) (found *model.OrderItem, insert bool) {
	for i, m := range items {
		if m.ID.Compare(ite.ID) {
			return &items[i], false
		}
	}

	return ite, true
}

// FindGift 查找订单中的礼品
func FindGift(items []model.Gift, ite *model.Gift) (found *model.Gift, insert bool) {
	for i, m := range items {
		if m.ID.Compare(ite.ID) {
			return &items[i], false
		}
	}

	return ite, true
}

func FindItemIndex(items []model.OrderItem, ite *model.OrderItem) (idx int, insert bool) {
	for i, m := range items {
		if m.ID.Compare(ite.ID) {
			return i, false
		}
	}

	return -1, true
}

func RemoveItem(items []model.OrderItem, ite *model.OrderItem) ([]model.OrderItem, bool) {
	if idx, _ := FindItemIndex(items, ite); idx >= 0 {
		if idx == len(items)-1 {
			return items[:idx], true
		} else {
			return append(items[:idx], items[idx+1:]...), true
		}
	} else {
		return items, false
	}
}
