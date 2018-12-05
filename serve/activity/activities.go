package activity

import (
	"strconv"

	. "activityjs.io/serve/constraint"
	"activityjs.io/serve/context"
	"activityjs.io/serve/model"
	. "activityjs.io/serve/reward"
)

// SpecialPrice 特价活动
func SpecialPrice(item *model.ActivityItem, quantity int, saving float64) *Activity {
	actx := context.NewActivityContext()
	actx.Metadata = context.ActivityMetadata{
		Kind:      "SpecialPrice",
		Condition: "数量满足 " + strconv.Itoa(quantity) + "件",
		Target:    "单价",
		Offer:     saving,
	}
	actx.AppliedLimit = 1
	return New(item, []Constrainter{
		Least("Quantity", quantity),
	}, actx, AppliedItemOnce("MinusPrice", MinusPrice(saving, false)))
}

// OverDecrease 满减活动
func OverDecrease(item *model.ActivityItem, total float64, saving float64) *Activity {
	return New(item, []Constrainter{
		Poster(LeastFloat("Order.Total", total)),
	}, nil, Joint(SubtotalMinus(saving, false), TotalMinus(saving, false)))
}
