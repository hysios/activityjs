package context

import (
	"testing"

	"activityjs.io/serve/utils"
)

func TestParseSelector(t *testing.T) {
	var (
		_, _, orderItems,
		order, user = utils.Prepare()
		item = orderItems[0]
		// item2 = orderItems[1]
	)

	// 执行上下文
	exeCtx := NewExecutionContext(
		&item,
		&user,
		&order,
	)
	selCtx := Selector{
		ExecutionContext: exeCtx,
	}

	ss := selCtx.parseSelector("Order.Count")
	if len(ss) != 2 {
		t.Fatalf("Order.Count 必须成功， 可以分成 [ Order, Count ] 但: %v", ss)
	}

	t.Logf("ParseSelector 结果: %#v", ss)

	ss = selCtx.parseSelector("XXX.Count")
	if len(ss) != 2 {
		t.Fatalf("Order.Count 必须成功， 可以分成 [ XXX, Count ] 但: %v", ss)
	}

	t.Logf("ParseSelector 结果: %#v", ss)

	ss = selCtx.parseSelector("User.Username")
	if len(ss) != 2 {
		t.Fatalf("Order.Count 必须成功， 可以分成 [ User, Username ] 但: %v", ss)
	}

	t.Logf("ParseSelector 结果: %#v", ss)

	ss, idx := selCtx.parseSubindex("Order.Items[3].Count")
	if len(ss) != 2 {
		t.Fatalf("Order.Items[3].Count 必须成功， 可以分成 [ Order.Items, Count ] 但: %v", ss)
	}

	if idx != 3 {
		t.Fatalf("Order.Items[3].Count idx 不正确，错误 %v", idx)
	}

	t.Logf("ParseSubindex 结果: %#v, %v", ss, idx)

}

func TestSelectorGet(t *testing.T) {
	var (
		_, _, orderItems,
		order, user = utils.Prepare()
		item = orderItems[0]
		// item2 = orderItems[1]
	)

	// 执行上下文
	exeCtx := NewExecutionContext(
		&item,
		&user,
		&order,
	)
	selCtx := NewSelector(exeCtx)

	itemquantity, ok := selCtx.GetInt("Quantity")
	if !ok {
		t.Fatalf("Quantity 必须成功 ")
	}
	t.Logf("GetInt -> Quantity 结果: %#v", itemquantity)

	// calc := activity.NewCalc(&item)
	// calc.Evaluate(exeCtx)

	count, ok := selCtx.GetInt("Order.Count")
	if !ok {
		t.Fatalf("Order.Count 必须成功 ")
	}

	t.Logf("GetInt -> Order.Count 结果: %#v", count)

	total, ok := selCtx.GetFloat64("Order.Total")
	if !ok {
		t.Fatalf("Order.Total 必须成功 ")
	}

	t.Logf("GetFloat64 -> Order.Total 结果: %#v", total)

	itemquan, ok := selCtx.GetInt("Order.Items[1].Quantity")
	if !ok {
		t.Logf("errors %# v", selCtx.Errors())
		t.Fatalf("Order.Items[1].Quantity 必须成功 ")
	}

	t.Logf("GetInt -> Order.Items[1].Quantity 结果: %#v", itemquan)

	// 超出边界
	_, ok = selCtx.GetInt("Order.Items[999].Quantity")
	if ok {
		t.Fatalf("Order.Items[999].Quantity  超出边界必须失败 ")
	}
	t.Logf("Order.Items[999].Quantity 错误 %s", selCtx.Errors())
}

func TestSelectorSet(t *testing.T) {
	var (
		_, _, orderItems,
		order, user = utils.Prepare()
		item = orderItems[0]
		// item2 = orderItems[1]
	)

	// 执行上下文
	exeCtx := NewExecutionContext(
		&item,
		&user,
		&order,
	)
	selCtx := Selector{
		ExecutionContext: exeCtx,
	}

	// calc := activity.NewCalc(&item)
	// calc.Evaluate(exeCtx)

	selCtx.Set("Order.Count", 10)

	count, ok := selCtx.Get("Order.Count")
	if !ok {
		t.Fatalf("Order.Count 必须成功 ")
	}

	if count != 10 {
		t.Fatalf("Order.Count 必须等于 10 ")
	}
	t.Logf("GetInt 结果: %#v", count)

	selCtx.Set("Order.Total", 100.0)

	total, ok := selCtx.GetFloat64("Order.Total")
	if !ok {
		t.Fatalf("Order.Total 必须成功 ")
	}
	if total != 100.0 {
		t.Fatalf("Order.Count 必须等于 100.0 ")
	}
	t.Logf("GetFloat64 结果: %#v", total)

	selCtx.Set("User.Username", "Jim")
	username, ok := selCtx.GetString("User.Username")
	if !ok {
		t.Fatalf("User.Username 必须成功 ")
	}

	if username != "Jim" {
		t.Fatalf("User.Username 必须等于 Jim ")
	}
	t.Logf("GetString 结果: %#v", username)

	selCtx.Set("Order.Items[0].Quantity", 10)
	itemquan, ok := selCtx.GetInt("Order.Items[0].Quantity")
	if !ok {
		t.Fatalf("Order.Items[0].Quantity 必须成功 ")
	}

	if itemquan != 10 {
		t.Fatalf("Order.Items[0].Quantity 必须等于 10")
	}
	t.Logf("GetInt -> Order.Items[0].Quantity结果: %#v", itemquan)

}
