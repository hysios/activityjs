package context_test

import (
	"log"
	"testing"

	"github.com/kr/pretty"

	"activityjs.io/serve/context"

	"activityjs.io/serve/utils"
)

func TestClone(t *testing.T) {
	var (
		_, _, orderItems,
		user, order = utils.Prepare()
		item = orderItems[0]
	)

	exeCtx := context.NewExecutionContext(&item, &order, &user)
	ctx := context.NewSelector(exeCtx)
	newCtx := ctx.Clone()

	ctx.Set("Subtotal", 133)
	ctx.Set("Order.Count", 13)
	ctx.Set("Order.Items[1].Price", 3833.0)

	log.Printf("oldCtx: %# v", pretty.Formatter(ctx))
	log.Printf("newCtx: %# v", pretty.Formatter(newCtx))

	st1, _ := ctx.GetFloat64("Subtotal")
	st2, _ := newCtx.GetFloat64("Subtotal")
	if st1 == st2 {
		t.Fatalf("Subtotal 应该不相同 %f %f", st1, st2)
	}
	tt1, _ := ctx.GetInt("Order.Count")
	tt2, _ := newCtx.GetInt("Order.Count")
	if tt1 == tt2 {
		t.Fatalf("Subtotal 应该不相同 %d %d", tt1, tt2)
	}
	ipa1, _ := ctx.GetFloat64("Order.Items[1].Price")
	ipa2, _ := newCtx.GetFloat64("Order.Items[1].Price")
	if ipa1 == ipa2 {
		t.Fatalf("Subtotal 应该不相同 %f %f", ipa1, ipa2)
	}

	oldOrder := ctx.Order()
	oldOrder.Items = append(oldOrder.Items, item)
	newOrder := newCtx.Order()
	if len(oldOrder.Items) == len(newOrder.Items) {
		t.Fatalf("Order Items 长度 应该不相同 %d %d", len(oldOrder.Items), len(newOrder.Items))
	}
}

func TestJSObject(t *testing.T) {
	var (
		_, _, orderItems,
		user, order = utils.Prepare()
		item = orderItems[0]
	)

	exeCtx := context.NewExecutionContext(&item, &order, &user)
	ctx := context.NewSelector(exeCtx)
	t.Logf("% #v", pretty.Formatter(ctx.JSObject()))
}

// func TestApply(t *testing.T) {
// 	var (
// 		_, _, orderItems,
// 		user, order = utils.Prepare()
// 		item = orderItems[0]
// 	)

// 	ctx := context.New(&item, &order, &user)

// 	newCtx, err := ctx.Apply()
// 	if err != nil {
// 		t.Fatalf("error %s", err)
// 	}
// 	t.Logf("% #v", pretty.Formatter(newCtx))
// }
