package activity_test

import (
	"testing"

	"github.com/kr/pretty"

	"activityjs.io/serve/activity"
	"activityjs.io/serve/context"

	"activityjs.io/serve/model"
)

func TestCalc(t *testing.T) {
	iphone := model.Product{
		// Name:         "apple iphone4",
		Title:        "苹果 IPhone4 掉渣天版",
		DefaultPrice: 4999.0,
	}

	ipad := model.Product{
		// Name:         "apple iPad pro",
		Title:        "苹果 IPad Pro版",
		DefaultPrice: 3999.0,
	}

	item := model.OrderItem{
		Title:    iphone.Title,
		Product:  iphone,
		Price:    iphone.DefaultPrice,
		Quantity: 20,
	}

	// 生成 ExecutionContext, 需要用户，订单等
	// 购买用户
	user := model.User{
		ID:       model.IntID(1),
		Username: "bobo",
		Level:    0,
	}

	// 设定订单
	order := model.Order{
		Items: []model.OrderItem{
			model.OrderItem{
				ID:       model.IntID(1),
				Title:    ipad.Title,
				Quantity: 4,
				Price:    ipad.DefaultPrice,
			},
			model.OrderItem{
				ID:       model.IntID(2),
				Title:    iphone.Title,
				Quantity: 5,
				Price:    iphone.DefaultPrice,
			},
		},
	}

	activity := activity.Calc{}

	// 执行上下文
	exeCtx := context.NewExecutionContext(
		&item,
		&user,
		&order,
	)

	if _, err := activity.Evaluate(exeCtx); err != nil {
		t.Fatalf("calc 不应该失败 %s", err)
	}

	t.Logf("ExecutionContext %# v", pretty.Formatter(exeCtx))

	if subtotal, ok := exeCtx.GetFloat64("Subtotal"); ok {
		if subtotal != 99980.0 {
			t.Fatalf("item 小计应该是 99980.0\n")
		}
	} else {
		t.Fatalf("没有获得 Subtotal 字段\n")
	}
}
