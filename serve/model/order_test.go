package model_test

import (
	"testing"

	"github.com/kr/pretty"

	"activityjs.io/serve/activity"
	"activityjs.io/serve/context"
	"activityjs.io/serve/utils"
)

func TestOrder(t *testing.T) {
	var (
		_, _, orderItems,
		order, _ = utils.Prepare()
		item = orderItems[0]
		// item2 = orderItems[1]
	)

	oldCount := order.Count
	calc := activity.Calc{}
	ctx := context.New(nil, nil, &order)
	calc.Ensure(ctx)

	order.AppendItem(&item)
	t.Logf("item %# v", pretty.Formatter(item))
	if order.Count == oldCount {
		t.Logf("Count 必须要加一")
	}
	t.Logf("order %# v", pretty.Formatter(order))

}
