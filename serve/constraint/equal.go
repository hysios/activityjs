package constraint

import (
	"strconv"

	"activityjs.io/serve/errors"
)

type equal struct {
	field string
	val   int
}

// Equal 约束某选择器的条件必须等于
func Equal(name string, val int) *equal {
	return &equal{name, val}
}

func (e *equal) Name() string {
	return e.field
}

func (e *equal) Evaluate(ctx Context) (err error) {
	if val, ok := ctx.ExecCtx.GetInt(e.field); ok {
		if e.val == val {
			return nil
		}
	}

	return errors.Wrap(errors.ErrExit, e.field+" 不等于 "+strconv.Itoa(e.val))
}
