package main

import (
	"activityjs.io/serve/activity"
	"activityjs.io/serve/context"
	"activityjs.io/serve/log"
	"activityjs.io/serve/model"
)

func Prepare() ([]model.Product, []model.ActivityItem, []model.OrderItem, model.Order, model.User) {
	var products = []model.Product{
		model.Product{
			ID:           model.IntID(10),
			Name:         "apple iphone4",
			Title:        "苹果 IPhone4 掉渣天版",
			DefaultPrice: 4999.0,
		},
		model.Product{
			ID:           model.IntID(11),
			Name:         "apple iPad pro",
			Title:        "苹果 IPad Pro版",
			DefaultPrice: 3999.0,
		},
	}

	var activityItems = []model.ActivityItem{
		model.ActivityItem{
			ID:      model.IntID(1),
			ItemID:  model.IntID(1),
			Title:   products[0].Title,
			Product: products[0],
			Price:   products[0].DefaultPrice,
		},
		model.ActivityItem{
			ID:      model.IntID(2),
			ItemID:  model.IntID(2),
			Title:   products[1].Title,
			Product: products[1],
			Price:   products[1].DefaultPrice,
		},
	}

	var orderItems = []model.OrderItem{
		model.OrderItem{
			ID:       activityItems[1].ItemID,
			Title:    products[1].Title,
			Product:  products[1],
			Quantity: 1,
			Price:    products[1].DefaultPrice,
		},
		// 无效提交
		model.OrderItem{
			ID:       activityItems[0].ItemID,
			Title:    products[0].Title,
			Product:  products[0],
			Price:    products[0].DefaultPrice,
			Quantity: 19, // 无效数量，小于20
		},
	}
	// 原有订单
	var order = model.Order{
		Items: orderItems,
	}

	// 购买用户
	var user = model.User{
		ID:       model.IntID(1),
		Username: "bobo",
		Level:    0,
	}

	return products, activityItems, orderItems, order, user
}

func main() {
	var (
		_, activityItems, orderItems,
		user, order = Prepare()
		actIphone = activityItems[0]
		actIpad   = activityItems[1]
		item      = orderItems[0]
	)

	// 这是一个当数量大于等于20时， 单价会降 3 元的活动
	// 活动对象主体是 iphone4 手机
	var (
		activities = []*activity.Activity{
			// 特价活动
			activity.SpecialPrice(&actIphone, 20, 3.0),
			// 满减活动
			activity.OverDecrease(&actIpad, 5000.0, 100),
		}

		exeCtx                 = context.NewExecutionContext(&item, &order, &user)
		ctx    context.Context = context.NewSelector(exeCtx)
		err    error
	)

	for _, activity := range activities {
		ctx, err = activity.Evaluate(ctx)
		if err != nil {
			log.Log("errors %s", err)
			// log.Printf("Errors %s", err)
		}
	}
	log.Log("ctx %# v", ctx)
}
