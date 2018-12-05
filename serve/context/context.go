// DONT EDIT: Auto generated

package context

import (
	"activityjs.io/serve/errors"
	"activityjs.io/serve/model"
	"github.com/gopherjs/gopherjs/js"
)

// Context for Activity Evaluation
type Context interface {
	Item() *model.OrderItem
	User() *model.User
	Order() *model.Order
	Errors() *errors.Errors
	// Get 获取 Context 值，任意类型
	Get(name string) (interface{}, bool)
	// Set 设置 Context 值
	Set(name string, val interface{})
	GetID(name string) (*model.Identity, bool)
	GetInt(name string) (int, bool)
	GetInt64(name string) (int64, bool)
	GetString(name string) (string, bool)
	GetBool(name string) (bool, bool)
	GetFloat64(name string) (float64, bool)
	GetTime(name string) (model.Time, bool)
	Clone() Context
	JSObject() map[string]interface{}
	LoadJS(state *js.Object)
	Now() model.Time
}
