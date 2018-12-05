package main

import (
	"activityjs.io/serve"
	"activityjs.io/serve/activity"
	"activityjs.io/serve/context"
	"activityjs.io/serve/model"
	"github.com/gopherjs/gopherjs/js"
)

//go:react:products
var products = []model.Product{
	model.Product{
		// Name:         "apple iphone4",
		Title:        "苹果 IPhone4 掉渣天版",
		DefaultPrice: 4999.0,
	},
	model.Product{
		// Name:         "apple iPad pro",
		Title:        "苹果 IPad Pro版",
		DefaultPrice: 3999.0,
	},
}

var activityItems = []model.ActivityItem{
	model.ActivityItem{
		Title:   products[0].Title,
		Product: products[0],
		Price:   products[0].DefaultPrice,
	},
	model.ActivityItem{
		Title:   products[1].Title,
		Product: products[1],
		Price:   products[1].DefaultPrice,
	},
}

var orderItems = []model.OrderItem{
	model.OrderItem{
		ID:       model.IntID(1),
		Title:    products[1].Title,
		Product:  products[1],
		Quantity: 1,
		Price:    products[1].DefaultPrice,
	},
	// 无效提交
	model.OrderItem{
		ID:       model.IntID(2),
		Title:    products[0].Title,
		Product:  products[0],
		Price:    products[0].DefaultPrice,
		Quantity: 19, // 无效数量，小于20
	},
}

// 购买用户
var user = model.User{
	ID:       model.StringID("uuid"),
	Username: "bobo",
	Level:    0,
}

// 原有订单
var order = model.Order{
	Items: orderItems,
}

func log(args string) {
	js.Global.Get("console").Call("log", args)
}

var mach *serve.Machine

func main() {
	var (
		actIphone = activityItems[0]
		actIpad   = activityItems[1]
		item      = orderItems[0]
	)

	// defultCfg := serve.DefaultConfig()
	var (
		ctx = context.New(&item, &user, &order)
	)

	mach = serve.NewMachine(ctx)
	mach.AddActivity(activity.SpecialPrice(&actIphone, 20, 3.0))
	mach.AddActivity(activity.OverDecrease(&actIpad, 5000.0, 100))

	// ctx, err = mach.Evaluate(ctx)
	// if err != nil {
	// 	log("Machine Evaluate: " + err.Error())
	// }

	js.Global.Set("ActivityMachine", map[string]interface{}{
		"evaluate": Evaluate,
		"context":  ctx.JSObject(),
		"load":     LoadCtx,
	})
}

func LoadCtx(state *js.Object) context.Context {
	ctx := context.New(nil, nil, nil)
	ctx.LoadJS(state)
	return ctx
}

func Evaluate(state *js.Object) map[string]interface{} {
	ctx := LoadCtx(state)
	ctx, err := mach.Evaluate(ctx)
	if err != nil {
		log("Machine Evaluate: " + err.Error())
	}
	return ctx.JSObject()
}
