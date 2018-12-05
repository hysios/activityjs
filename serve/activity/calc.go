package activity

import (
	"strconv"

	"activityjs.io/serve/model"

	"activityjs.io/serve/context"
	"activityjs.io/serve/errors"
)

// Calc 订单计算器
type Calc struct {
	Item *model.ActivityItem
}

func NewCalc(item *model.ActivityItem) *Calc {
	return &Calc{Item: item}
}

// Valid 订单计算器验证无效的数据
func (act *Calc) Valid(ctx context.Context) (errs *errors.Errors) {
	// ctx := constraint.Context{Item: act.Item, ExecCtx: exeCtx}
	errs = ctx.Errors()
	order := ctx.Order()

	itemLstErrs := &errors.Errors{}

	for i, item := range order.Items {
		var itmErrs = &errors.Errors{}
		if item.Quantity <= 0 {
			itmErrs.Add("Quantity", errors.New("数量不能为 0"))
		}

		if !itmErrs.Empty() {
			itemLstErrs.Add("Item-"+strconv.Itoa(i), itmErrs)
			continue
		}

	}
	if !itemLstErrs.Empty() {
		errs.Add("Items", itemLstErrs)
	}

	if quantity, ok := ctx.GetInt("Quantity"); !ok {
		errs.Add("Quantity", errors.New("Item 数量字段不存在"))

	} else if quantity <= 0 {
		errs.Add("Quantity", errors.New("Item 数量不能小于或等于 0 "))
	}

	if price, ok := ctx.GetFloat64("Price"); !ok {
		errs.Add("Price", errors.New("Item 单价字段不存在"))
	} else if price < 0 {
		errs.Add("Price", errors.New("单价不能为负"))
	}

	if errs.Empty() {
		return nil
	}
	return errs
}

// Evaluate 计算器验值过程
func (act *Calc) Evaluate(ctx context.Context) (context.Context, error) {

	errs := act.Valid(ctx)
	if errs != nil {
		return nil, errs
	}

	var (
		sum   float64
		count int
		order = ctx.Order()
	)

	for i, item := range order.Items {
		order.Items[i].Subtotal = float64(item.Quantity) * item.Price
		count += order.Items[i].Quantity
		sum += order.Items[i].Subtotal
	}

	if quantity, ok := ctx.GetInt("Quantity"); ok {
		if price, ok := ctx.GetFloat64("Price"); ok {
			total := float64(quantity) * price
			ctx.Set("Subtotal", total)
		}
	}
	order.Count = count
	order.Total = sum
	return ctx, nil
}

func (act *Calc) Ensure(ctx context.Context) context.Context {

	var (
		sum   float64
		count int
		order = ctx.Order()
		item  = ctx.Item()
	)

	if order != nil {
		for i, item := range order.Items {
			order.Items[i].Subtotal = float64(item.Quantity) * item.Price
			// order.Items[i].SetOrder(order)
			count += order.Items[i].Quantity
			sum += order.Items[i].Subtotal
		}
		order.Count = count
		order.Total = sum
	}

	if item != nil {

		if quantity, ok := ctx.GetInt("Quantity"); ok {
			if price, ok := ctx.GetFloat64("Price"); ok {
				total := float64(quantity) * price
				ctx.Set("Subtotal", total)
			}
		}
	}

	return ctx
}
