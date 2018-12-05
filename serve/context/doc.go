package context

import "activityjs.io/serve/model"

func New(item *model.OrderItem, user *model.User, order *model.Order, opt ...ExecutionOptions) Context {
	var (
		exeCtx = NewExecutionContext(item, user, order, opt...)
		ctx    = NewSelector(exeCtx)
	)
	return ctx
}

// func NewBlank() Context {
// 	var (
// 		exeCtx = NewBlankExecutionContext()
// 		ctx    = NewSelector(exeCtx)
// 	)
// 	return ctx
// }
