package reward

import (
	"activityjs.io/serve/context"
	"activityjs.io/serve/errors"
	"activityjs.io/serve/log"
	"activityjs.io/serve/model"
	"activityjs.io/serve/utils"
)

// Handler 奖励回调函数的定义
type Handler func(context.RewardContext) error

// MinusPrice 减单价
func MinusPrice(minus float64, free bool) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx
		ctxPrice, ok := ctx.GetFloat64("Price")
		if !ok {
			return errors.New("没有价格字段")
		}

		if ctxPrice-minus <= 0 && !free {
			return errors.New("减价不能小于 0")
		} else if ctxPrice-minus < 0 && free {
			return errors.New("减价不能为负")
		}

		return rctx.WithEffect("MinusPrice", "Price", -minus, func() error {
			ctxPrice -= minus
			// ctx.StopPropagation()
			log.Log("Minus Price %v", ctxPrice)
			ctx.Set("Price", ctxPrice)
			return nil
		})
	}
}

// ChangePrice 改单价
func ChangePrice(price float64, free bool) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx

		if price <= 0 && !free {
			return errors.New("改单价不能小于 0")
		} else if price < 0 && free {
			return errors.New("改单价不能为负")
		}

		ctx.Set("Price", price)
		return nil
	}
}

// PriceDiscount 单价折扣
func PriceDiscount(precent float64, free bool) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx

		ctxPrice, ok := ctx.GetFloat64("Price")
		if !ok {
			return errors.New("没有价格字段")
		}

		if ctxPrice*precent <= 0 && !free {
			return errors.New("价格不能小于等于 0")
		} else if ctxPrice*precent < 0 && free {
			return errors.New("价格不能为负")
		}

		ctxPrice *= precent
		ctx.Set("Price", ctxPrice)
		return nil
	}
}

// TotalMinus 总计上优惠
func TotalMinus(minus float64, free bool) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx
		orderTotal, ok := ctx.GetFloat64("Order.Total")
		if !ok {
			return errors.New("没有总计字段")
		}

		if orderTotal-minus <= 0 && !free {
			return errors.New("总计不能小于 0")
		} else if orderTotal-minus < 0 && free {
			return errors.New("总计不能为负")
		}

		orderTotal -= minus
		ctx.Set("Order.Total", orderTotal)
		return nil
	}
}

func SubtotalMinus(minus float64, free bool) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx
		subtotal, ok := ctx.GetFloat64("Subtotal")
		if !ok {
			return errors.New("没有小计字段")
		}

		if subtotal-minus <= 0 && !free {
			return errors.New("小计不能小于 0")
		} else if subtotal-minus < 0 && free {
			return errors.New("小计不能为负")
		}

		subtotal -= minus
		ctx.Set("Subtotal", subtotal)
		return nil
	}
}

// TotalDiscount 小计上折扣
func SubtotalDiscount(precent float64, free bool) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx
		ctxQuantity, ok := ctx.GetInt("Quantity")
		if !ok {
			return errors.New("没有数量字段")
		}
		if ctxQuantity <= 0 {
			return errors.New("数量不能少于 0")
		}

		ctxPrice, ok := ctx.GetFloat64("Price")
		if !ok {
			return errors.New("没有价格字段")
		}

		total := ctxPrice * float64(ctxQuantity)
		total = total * precent
		if total <= 0 && !free {
			return errors.New("价格不能小于等于 0")
		} else if total < 0 && free {
			return errors.New("价格不能为负")
		}
		ctx.Set("Subtotal", total)
		ctx.Set("Price", total/float64(ctxQuantity))
		// ctx.Subtotal = total
		// ctx.Price = total / float64(ctx.Quantity)
		return nil
	}
}

// AddGiftProduct 赠送一个产品
func AddGiftProduct(item *model.OrderItem, price float64, quantity int) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx
		if quantity == 0 {
			quantity = 1
		}

		order := ctx.Order()

		if item, insert := utils.FindItem(order.Items, item); insert {
			item.Price = price
			item.Quantity = quantity
			order.Items = append(order.Items, *item)
		} else {
			if item.Price == price {
				item.Quantity += quantity
			} else {
				item.Price = price
				item.Quantity = quantity
				order.Items = append(order.Items, *item)
			}
		}
		return nil
	}
}

// AddGift 赠送一个产品
func AddGift(gift *model.Gift, price float64, quantity int) Handler {
	return func(rctx context.RewardContext) error {
		ctx := rctx.Ctx
		if quantity == 0 {
			quantity = 1
		}
		order := ctx.Order()

		if gift, insert := utils.FindGift(order.Gifts, gift); insert {
			gift.Price = price
			gift.Quantity = quantity
			order.Gifts = append(order.Gifts, *gift)
		} else {
			if gift.Price == price {
				gift.Quantity += quantity
			} else {
				gift.Price = price
				gift.Quantity = quantity
				order.Gifts = append(order.Gifts, *gift)
			}
		}
		return nil
	}
}
