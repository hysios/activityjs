package constraint

import (
	"strconv"

	"activityjs.io/serve/errors"
)

type multi struct {
	fields []string
	val    int
}

type multiFloat struct {
	fields []string
	val    float64
}

// Multi 约束多个字段条件至少
func Multi(val int, fields ...string) *multi {
	return &multi{fields: fields, val: val}
}

func (e *multi) Name() string {
	return e.fields[0]
}

func (e *multi) Evaluate(ctx Context) (err error) {
	s := 0
	for _, field := range e.fields {
		if val, ok := ctx.ExecCtx.GetInt(field); ok {
			s += val
			if s >= e.val {
				return nil
			}
		}
	}

	return errors.Wrap(errors.ErrExit, "至少 "+strconv.Itoa(e.val))
}

func MultiFloat(val float64, fields ...string) *multiFloat {
	return &multiFloat{fields: fields, val: val}
}

func (e *multiFloat) Name() string {
	return e.fields[0]
}

func (e *multiFloat) Evaluate(ctx Context) (err error) {
	var s float64
	for _, field := range e.fields {
		if val, ok := ctx.ExecCtx.GetFloat64(field); ok {
			s += val
			if s >= e.val {
				return nil
			}
		}
	}

	return errors.Wrap(errors.ErrExit, "至少 "+strconv.FormatFloat(e.val, 'f', -1, 64))
}
