package activity_test

import (
	"testing"
	"time"

	"activityjs.io/serve/activity"
	"activityjs.io/serve/constraint"
	"activityjs.io/serve/context"
	"activityjs.io/serve/model"
	"activityjs.io/serve/reward"
	"activityjs.io/serve/utils"
	"github.com/kr/pretty"
)

type Constrainter = constraint.Constrainter

func TestActivity(t *testing.T) {
	var (
		_, activityItems, orderItems,
		order, user = utils.Prepare()
		aitem = activityItems[0]
		item  = orderItems[0]
		item2 = orderItems[1]
	)

	// 这是一个当数量大于等于20时， 单价会降 3 元的活动
	// 活动对象主体是 iphone4 手机
	activity := activity.New(&aitem, []Constrainter{
		constraint.Equal("Quantity", 20),
	}, nil, reward.MinusPrice(3, false))

	t.Logf("%# v", pretty.Formatter(activity))

	// 原有订单
	order = model.Order{
		Items: order.Items[:1],
	}

	// 执行上下文
	exeCtx := context.NewExecutionContext(
		&item,
		&user,
		&order,
	)

	// 运算
	ctx, err := activity.Evaluate(exeCtx)
	if err != nil {
		t.Logf("errors: %v", err)
	}

	// 运算后结果
	t.Logf("有效活动 Execution After: %# v", pretty.Formatter(ctx))

	// 生成新上下文
	exeCtx = context.NewExecutionContext(
		&item2,
		&user,
		&order,
	)

	// 运算验证
	_, err = activity.Evaluate(exeCtx)
	if err != nil {
		t.Logf("errors: %s", err)
	}

	errs := exeCtx.Errors()
	if err := errs.Get("Quantity"); err == nil {
		t.Logf("[预期]会提示数量错误[Quantity]!")
	}

	if price, ok := exeCtx.GetFloat64("Price"); ok {
		if price != 4999.0 {
			t.Logf("[Price] 价格不应该发生改变 %.2f!", price)
		}
	}

	t.Logf("Execution After: %# v", pretty.Formatter(exeCtx))
}

func TestActivityLeast(t *testing.T) {
	var (
		_, activityItems, orderItems,
		order, user = utils.Prepare()
		aitem = activityItems[0]
		item  = orderItems[0]
		item2 = orderItems[1]
	)

	// 这是一个当数量大于等于20时， 单价会降 3 元的活动
	// 活动对象主体是 iphone4 手机
	activity := activity.New(&aitem, []Constrainter{
		constraint.Least("Quantity", 20),
	}, nil, reward.MinusPrice(3, false))

	t.Logf("%# v", pretty.Formatter(activity))

	// 生成 ExecutionContext, 需要用户，订单等

	// 原有订单
	order = model.Order{
		Items: order.Items[:1],
	}
	item.Quantity = 20

	// 执行上下文
	exeCtx := context.NewExecutionContext(
		&item,
		&user,
		&order,
	)

	// 验证有效性
	errs := activity.Valid(exeCtx)

	// 运算
	_, err := activity.Evaluate(exeCtx)
	if err != nil {
		t.Fatalf("errors: %v", err)
	}

	// 运算后结果
	t.Logf("有效活动 Execution After: %# v", pretty.Formatter(exeCtx))

	// 无效提交
	item2.Quantity = 25

	// 生成新上下文
	exeCtx = context.NewExecutionContext(
		&item2,
		&user,
		&order,
	)

	// 运算验证
	ctx, err := activity.Evaluate(exeCtx)
	if err != nil {
		t.Fatalf("errors: %s", err)
	}

	errs = exeCtx.Errors()
	if err := errs.Get("Quantity"); err != nil {
		t.Fatalf("不能发生数量错误[Quantity]!")
	}

	if price, ok := ctx.GetFloat64("Price"); ok {
		if price != 4996.0 {
			t.Fatalf("[Price] 价格应该发生改变 %.2f!", price)
		}
	}

	t.Logf("Execution After: %# v", pretty.Formatter(exeCtx))
}

func TestActivityDate(t *testing.T) {
	var (
		_, activityItems, orderItems,
		order, user = utils.Prepare()
		aitem = activityItems[0]
		item2 = orderItems[1]
	)

	tt, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	// 这是一个当数量大于等于20时， 时间在大于 2016-01-02 15:04:05
	// 单价会降 3 元的活动
	// 活动对象主体是 iphone4 手机
	activity := activity.New(&aitem, []Constrainter{
		constraint.Least("Quantity", 20),
		constraint.StartAt(model.Time(tt.Unix() * 1000)),
	}, nil, reward.MinusPrice(3, false))

	t.Logf("%# v", pretty.Formatter(activity))

	// 生成 ExecutionContext, 需要用户，订单等

	// 原有订单
	order = model.Order{
		Items: order.Items[:1],
	}

	item2.Quantity = 20

	// 生成新上下文
	exeCtx := context.NewExecutionContext(
		&item2,
		&user,
		&order,
		context.ExecutionOptions{
			NowHandler: func() int64 {
				// 使用当前时间，会通过
				return time.Now().Unix() * 1000
			},
		},
	)

	// 运算验证
	_, err := activity.Evaluate(exeCtx)
	if err != nil {
		t.Fatalf("errors: %s", err)
	}

	// 生成新上下文
	exeCtx = context.NewExecutionContext(
		&item2,
		&user,
		&order,
		context.ExecutionOptions{
			NowHandler: func() int64 {
				tt, _ := time.Parse(time.RFC3339, "2006-01-01T15:04:05Z")
				// 使用当前时间，会通过
				return tt.Unix()
			},
		},
	)

	// 运算验证
	_, err = activity.Evaluate(exeCtx)
	if err == nil {
		t.Fatalf("时间验证不会通过: %s", err)
	}
	t.Logf("时间验证不能过的错误 %s", err)
}
