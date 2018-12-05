package constraint

import (
	"strconv"

	"activityjs.io/serve/errors"
)

type least struct {
	field string
	val   int
}

type leastFloat64 struct {
	field string
	val   float64
}

// Least 约束某选择器的条件至少
func Least(name string, val int) *least {
	return &least{name, val}
}

func (e *least) Name() string {
	return e.field
}

func (e *least) Evaluate(ctx Context) (err error) {
	if val, ok := ctx.ExecCtx.GetInt(e.field); ok {
		if val >= e.val {
			return nil
		}
	}

	return errors.Wrap(errors.ErrExit, e.field+" 至少 "+strconv.Itoa(e.val))
}

func (e *least) String() string {
	return "[Least] " + "字段 " + e.field + " 至少满足 " + strconv.Itoa(e.val)
}

func LeastFloat(name string, val float64) *leastFloat64 {
	return &leastFloat64{name, val}
}

func (e *leastFloat64) Name() string {
	return e.field
}

func (e *leastFloat64) Evaluate(ctx Context) (err error) {
	if val, ok := ctx.ExecCtx.GetFloat64(e.field); ok {
		if val >= e.val {
			return nil
		}
	}

	return errors.Wrap(errors.ErrExit, e.field+" 至少 "+strconv.FormatFloat(e.val, 'f', -1, 64))
}

func (e *leastFloat64) String() string {
	return "[Least(float64)] " + "字段 " + e.field + " 至少满足 " + strconv.FormatFloat(e.val, 'f', -1, 64)
}
