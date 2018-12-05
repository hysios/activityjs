// +build !js

package utils

//TODO: 这个包要在编译后排除
import "activityjs.io/serve/model"

func Prepare() ([]model.Product, []model.ActivityItem, []model.OrderItem, model.Order, model.User) {
	var products = []model.Product{
		model.Product{
			ID: model.IntID(10),
			// Name:         "apple iphone4",
			Title:        "苹果 IPhone4 掉渣天版",
			DefaultPrice: 4999.0,
		},
		model.Product{
			ID: model.IntID(11),
			// Name:         "apple iPad pro",
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
