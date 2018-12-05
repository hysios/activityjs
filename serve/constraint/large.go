package constraint

import (
	"strconv"

	"activityjs.io/serve/errors"
)

type large struct {
	field string
	val   int
}

type largeThan struct {
	field string
	val   int
}

// Large 约束某选择器的条件大于
func Large(name string, val int) *large {
	return &large{name, val}
}

func (e *large) Name() string {
	return e.field
}

func (e *large) Evaluate(ctx Context) (err error) {
	if val, ok := ctx.ExecCtx.GetInt(e.field); ok {
		if val > e.val {
			return nil
		}
	}

	return errors.Wrap(errors.ErrExit, e.field+" 要大于 "+strconv.Itoa(e.val))
}

// LargeThan 约束某选择器的条件大于
func LargeThan(name string, val int) *largeThan {
	return &largeThan{name, val}
}

func (e *largeThan) Name() string {
	return e.field
}

func (e *largeThan) Evaluate(ctx Context) (err error) {
	if val, ok := ctx.ExecCtx.GetInt(e.field); ok {
		if val >= e.val {
			return nil
		}
	}

	return errors.Wrap(errors.ErrExit, e.field+" 要大于或等于 "+strconv.Itoa(e.val))
}
