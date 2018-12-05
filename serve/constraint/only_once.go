package constraint

import (
	"activityjs.io/serve/errors"
)

type onlyOnce struct {
	which Constrainter
	tick  int
}

// OnlyOnce 验证一次
func OnlyOnce() *onlyOnce {
	return &onlyOnce{}
}

func (e *onlyOnce) Name() string {
	return "onlyOnce_" + e.which.Name()
}

func (e *onlyOnce) Binding(bind Constrainter) error {
	if _, ok := bind.(*onlyOnce); ok {
		return errors.New("OnlyOnce 前置条件不能为 OnlyOnce")
	}

	e.which = bind
	return nil
}

func (e *onlyOnce) Reset() {
	e.which = nil
	e.tick = 0
}

func (e *onlyOnce) Evaluate(ctx Context) (err error) {
	n := len(ctx.queue)

	if n == 0 {
		return errors.New("OnlyOnce 前面必须有约束条件")
	}

	for i, _ := range ctx.queue {
		if ctx.queue[i] == e.which {
			e.tick++
		}

		if e.tick > 1 {
			return errors.New("约束只能执行一次")
		}
	}

	return
}
