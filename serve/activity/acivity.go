package activity

import (
	"activityjs.io/serve/constraint"

	"activityjs.io/serve/context"
	"activityjs.io/serve/errors"
	"activityjs.io/serve/model"
	"activityjs.io/serve/reward"
)

// ErrAbort 终止错误，停止验证检查，返回失败
var ErrAbort = errors.New("Abort Evalution")

// Activity 活动的描述类，只针对一个产品并不包括订单
type Activity struct {
	Constraints []constraint.Constrainter
	Item        *model.ActivityItem
	ctx         *context.ActivityContext
	Reward      reward.Handler
	name        string
}

// New 构造一个产品活动，
// 活动是对一个参与产品的描述，是有多个约束条件，一个上下文参数，一个奖励函数组成
// New(ite, []Constraint{
// 		constraint.Least("Quantity", 20),
// 		constraint.Monopolize() // 独占
// }, nil, func(ctx *context.ExecutionContext) error {
//      item, _ := ctx.Item()
//		item.Price -= 3
//		return nil
// })
//
// 使用奖励函数
// New(ite, []Constraint{
// 		constraint.Least("Quantity", 20),	// 至少 20 件
// 		constraint.StartAt("2018-12-12"), 	// 开始于 2018-12-12
// 		constraint.EndAt("2018-12-18"),		// 结束于 2018-12-18
// }, nil, PriceChange(30))
//
// 买一赠一
// New(ite, []Constraint{
// 		constraint.Least("Quantity", 20),	// 至少 20 件
// 		constraint.StartAt("2018-12-12"), 	// 开始于 2018-12-12
// 		constraint.EndAt("2018-12-18"),		// 结束于 2018-12-18
// }, nil, GiftProduct(item, 1))
func New(item *model.ActivityItem, cons []constraint.Constrainter, ctx *context.ActivityContext, fn reward.Handler) *Activity {
	if ctx == nil {
		ctx = &context.ActivityContext{}
	}

	act := &Activity{ctx: ctx, Constraints: cons, Item: item, Reward: fn}
	if len(ctx.Name) == 0 {
		ctx.Name = act.Name()
	}

	return act
}

func (act *Activity) Name() string {
	if len(act.name) > 0 {
		return act.name
	}
	for _, con := range act.Constraints {
		act.name += con.Name() + ";"
	}
	return act.name
}

// Valid 验证活动有效性
func (act *Activity) Valid(exeCtx context.Context) (errs *errors.Errors) {
	ctx := constraint.Context{Item: act.Item, ActivityCtx: act.ctx, ExecCtx: exeCtx}
	errs = exeCtx.Errors()
	actErrs := errors.Errors{}

	for _, con := range act.Constraints {
		err := con.Evaluate(ctx)
		if err == nil {
			continue
		} else if err != ErrAbort {
			actErrs.Add(con.Name(), err)
		} else {
			actErrs.Add("exit", err)
			return
		}
	}

	if actErrs.Empty() {
		return nil
	}

	errs.Add(act.Name(), &actErrs)
	return &actErrs
}

// Evaluate 活动规则运算
func (act *Activity) Evaluate(exeCtx context.Context) (context.Context, error) {
	var (
		ctx  = exeCtx.Clone()
		errs = act.Valid(exeCtx)
	)

	if errs == nil {
		act.Reward(context.RewardContext{Ctx: ctx, ActivityContext: act.ctx})
		return ctx, nil
	} else if errs != nil && errs.NotStop() {
		act.Reward(context.RewardContext{Ctx: ctx, ActivityContext: act.ctx})
	}
	return ctx, errs
}

func (act *Activity) Related(item *model.OrderItem) bool {
	if act.Item == nil {
		return true
	}

	if !act.Item.ItemID.Nil() && act.Item.ItemID.Compare(item.ID) {
		return true
	} else {
		return false
	}
}

func (act *Activity) Metadata() *context.ActivityMetadata {
	if act.ctx == nil {
		return nil
	}

	return &act.ctx.Metadata
}
