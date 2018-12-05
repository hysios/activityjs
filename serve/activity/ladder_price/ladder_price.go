package ladder_price

import (
	"activityjs.io/serve/context"
	"activityjs.io/serve/model"
	"activityjs.io/serve/utils"
)

// ActivityLadder 阶梯价格活动
type ActivityLadder struct {
	ctx *context.ExecutionContext
}

// NewActivityLadder 通过 ctx 创建 ActivityLadder
func New(ctx *context.ExecutionContext) *ActivityLadder {
	return &ActivityLadder{ctx}
}

// Add 添加产品到订单
func (act *ActivityLadder) Add(incItem *model.OrderItem) {
	order := act.ctx.Order()
	if item, insert := utils.FindItem(order.Items, incItem); insert {
		order.Items = append(order.Items, *incItem)
	} else {
		item.Quantity += incItem.Quantity
	}
}

func (act *ActivityLadder) Enter() {

}
