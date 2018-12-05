package constraint

import (
	"activityjs.io/serve/context"
	"activityjs.io/serve/model"
)

// Context 约束的执行 Context
type Context struct {
	Item        *model.ActivityItem
	ActivityCtx *context.ActivityContext
	ExecCtx     context.Context
	queue       []Constrainter
}
