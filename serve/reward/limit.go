package reward

import (
	"activityjs.io/serve/context"
	"activityjs.io/serve/log"
)

func AppliedItemOnce(name string, fn Handler) Handler {
	return func(rctx context.RewardContext) error {
		var (
			ctx   = rctx.Ctx
			key   = rctx.ActivityContext.Name + name
			count int
		)

		// if rctx.ActivityContext.AppliedLimit == 0 {
		// 	return fn(rctx)
		// }

		itm := ctx.Item()
		if itm == nil {
			return fn(rctx)
		}

		for _, eff := range itm.Effects() {
			if eff.Key == key {
				count++
			}
		}

		log.Log("Applied Count %d", count)
		if count < 1 {

			return fn(rctx)
		}

		return nil
	}
}

func WrapEffect(name, field string, val interface{}, fn Handler) Handler {
	return func(rctx context.RewardContext) error {
		var (
			ctx     = rctx.Ctx
			key     = rctx.ActivityContext.Name + name
			summary = rctx.GetSummary(field)
		)

		return ctx.Item().WithEffect(key, field, summary, val, func() error {
			return fn(rctx)
		})
	}
}
