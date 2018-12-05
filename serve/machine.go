package serve

import (
	"activityjs.io/serve/activity"
	"activityjs.io/serve/constraint"
	"activityjs.io/serve/context"
	"activityjs.io/serve/log"
	"activityjs.io/serve/model"
)

// Machine 活动机
type Machine struct {
	activities       []*activity.Activity `header:"activities"`
	posterActivities []*activity.Activity `header:"Poster activities"`
	ctx              context.Context
}

// NewMachine 新建活动机
func NewMachine(ctx context.Context) *Machine {
	return &Machine{ctx: ctx}
}

// AddActivity 加入活动
func (mac *Machine) AddActivity(act *activity.Activity) {
	if mac.hasPoster(act) {
		mac.addPosterActivity(act)
	} else {
		mac.addActivity(act)
	}
}

func (mac *Machine) Evaluate(ctx context.Context) (context.Context, error) {
	var err error
	// Before:

	calc := activity.Calc{}
	ctx = calc.Ensure(ctx)

	if item := ctx.Item(); item != nil {
		ctx, err = mac.evaluateItem(item, ctx)
	}
	log.Log("ctx %# v", ctx)

	// itemCtx := ctx.Clone()
	if order := ctx.Order(); order != nil {
		for _, item := range order.Items {
			if ctx, err = mac.evaluateItem(&item, ctx); err != nil {
				return ctx, err
			}
			log.Log("ctx %# v", ctx)
		}
	}
	// ctx = itemCtx

	// activities := mac.relatedActivities(mac.activities,
	// for _, act := range mac.activities {
	// 	ctx, err = act.Evaluate(ctx)
	// }

	// // ctx, err = ctx.Apply() 不做 AddToCard
	// // Poster:
	// for _, act := range mac.posterActivities {
	// 	ctx, err = act.Evaluate(ctx)
	// }

	// ctx, err = ctx.Apply()
	return ctx, err
}

func (mac *Machine) evaluateItem(item *model.OrderItem, ctx context.Context) (context.Context, error) {
	var err error
	activities := mac.relatedActivities(mac.activities, item)
	for _, act := range activities {
		ctx, err = act.Evaluate(ctx)
	}

	activities = mac.relatedActivities(mac.posterActivities, item)
	for _, act := range activities {
		ctx, err = act.Evaluate(ctx)
	}

	return ctx, err
}

func (mac *Machine) hasPoster(act *activity.Activity) bool {
	for _, con := range act.Constraints {
		if constraint.IsPoster(con) {
			return true
		}
	}

	return false
}

func (mac *Machine) addActivity(act *activity.Activity) {
	mac.activities = append(mac.activities, act)
}

func (mac *Machine) addPosterActivity(act *activity.Activity) {
	mac.posterActivities = append(mac.posterActivities, act)
}

func (mac *Machine) relatedActivities(activities []*activity.Activity, item *model.OrderItem) []*activity.Activity {
	var relActivities = make([]*activity.Activity, 0, len(activities))
	for _, act := range activities {
		if act.Related(item) {
			relActivities = append(relActivities, act)
		}
	}
	return relActivities
}
