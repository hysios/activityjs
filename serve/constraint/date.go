package constraint

import (
	"activityjs.io/serve/model"

	"activityjs.io/serve/errors"
)

type startAt struct {
	val model.Time
}

// StartAt 约束时间必须大于条件
func StartAt(val model.Time) *startAt {
	return &startAt{val}
}

func (e *startAt) Name() string {
	return "StartAt"
}

func (e *startAt) Evaluate(ctx Context) (err error) {
	now := ctx.ExecCtx.Now()
	if now >= e.val {
		return nil
	}

	return errors.Wrap(errors.ErrExit, "时间必须要在 "+e.val.String()+" 之后")
}

type endAt struct {
	val model.Time
}

// EndAt 约束时间必须大于条件
func EndAt(val model.Time) *endAt {
	return &endAt{val}
}

func (e *endAt) Name() string {
	return "EndAt"
}

func (e *endAt) Evaluate(ctx Context) (err error) {
	now := ctx.ExecCtx.Now()
	if now <= e.val {
		return nil
	}

	return errors.Wrap(errors.ErrExit, "时间必须要在 "+e.val.String()+" 之前")
}
