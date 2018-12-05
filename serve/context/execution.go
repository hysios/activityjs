package context

import (
	"activityjs.io/serve/utils"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jsbuiltin"

	"activityjs.io/serve/errors"
	"activityjs.io/serve/model"
)

//go:generate ifacemaker -f execution.go -s ExecutionContext -i Context -p context -y "Context for Activity Evaluation" -c "DONT EDIT: Auto generated" -o context.go

// ExecutionContext 执行上下文
// type ExecutionContext map[interface{}]interface{}

type ExecutionContext struct {
	item       *model.OrderItem
	order      *model.Order
	user       *model.User
	errors     *errors.Errors
	options    ExecutionOptions
	upsideProp bool
}

var (
	itemkey    = "item"
	userkey    = "user"
	orderkey   = "order"
	errorkey   = "errors"
	optionskey = "options"
	applykey   = "apply"
)

type ExecutionOptions struct {
	NowHandler func() int64
}

type ActionHandler func(Context) error

var defaultOption = ExecutionOptions{}

var defaultActionHandler = AddToOrder

func NewExecutionContext(item *model.OrderItem, user *model.User, order *model.Order, opt ...ExecutionOptions) *ExecutionContext {
	var ectx = ExecutionContext{
		item:   item,
		user:   user,
		order:  order,
		errors: &errors.Errors{},
	}
	if len(opt) > 0 {
		ectx.options = opt[0]
	} else {
		ectx.options = defaultOption
	}

	return &ectx
}

func (ectx *ExecutionContext) Item() *model.OrderItem {
	return ectx.item
}

func (ectx *ExecutionContext) User() *model.User {
	return ectx.user
}

func (ectx *ExecutionContext) Order() *model.Order {
	return ectx.order
}

func (ectx *ExecutionContext) Errors() *errors.Errors {
	return ectx.errors
}

// Get 获取 Context 值，任意类型
func (ectx *ExecutionContext) Get(name string) (interface{}, bool) {
	item := ectx.Item()
	if item == nil {
		return nil, false
	}

	if val, ok := item.Get(name); ok {
		return val, ok
	}
	ectx.Errors().Add("builtins", errors.New("无效的字段 "+name))
	return nil, false
}

// Set 设置 Context 值
func (ectx *ExecutionContext) Set(name string, val interface{}) {
	ectx.sloppySet(name, val)
}

// sloppySet 数字类型宽松设置，自动转换 float <=> int, 或 bool, string 数字等类型
func (ectx *ExecutionContext) sloppySet(name string, val interface{}) {
	item := ectx.Item()
	if item == nil {
		ectx.Errors().Add("builtins", errors.New("没有 item 记录"))
		return
	}

	v, _ := item.Get(name)

	switch v.(type) {
	case int32, int64, int:
		if n, ok := utils.IsInt(val); ok {
			item.Set(name, n)
		} else {
			ectx.Errors().Add("builtins", errors.New("无效的格式转换成整型"))
		}
	case float32, float64:
		if n, ok := utils.IsFloat(val); ok {
			item.Set(name, n)
		} else {
			ectx.Errors().Add("builtins", errors.New("无效的格式转换成浮点"))
		}
	default:
		// ectx[name] = val
		item.Set(name, val)
	}

	if ectx.upsideProp {
		ectx.upsidePropagation()
	}
}

func (ectx *ExecutionContext) GetID(name string) (*model.Identity, bool) {
	item := ectx.Item()
	if item == nil {
		return nil, false
	}

	if val, ok := item.Get(name); ok {
		if v, ok := val.(*model.Identity); ok {
			return v, ok
		}
	}
	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return nil, false
}

func (ectx *ExecutionContext) GetInt(name string) (int, bool) {
	item := ectx.Item()
	if item == nil {
		return 0, false
	}
	if val, ok := item.Get(name); ok {
		if v, ok := val.(int); ok {
			return v, ok
		}
	}
	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return 0, false
}

func (ectx *ExecutionContext) GetInt64(name string) (int64, bool) {
	item := ectx.Item()
	if item == nil {
		return 0, false
	}

	if val, ok := item.Get(name); ok {
		if v, ok := val.(int64); ok {
			return v, ok
		}
	}
	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return 0, false
}

func (ectx *ExecutionContext) GetString(name string) (string, bool) {
	item := ectx.Item()
	if item == nil {
		return "", false
	}

	if val, ok := item.Get(name); ok {
		if v, ok := val.(string); ok {
			return v, ok
		}
	}
	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return "", false
}

func (ectx *ExecutionContext) GetBool(name string) (bool, bool) {
	item := ectx.Item()
	if item == nil {
		return false, false
	}

	if val, ok := ectx.Get(name); ok {
		if v, ok := val.(bool); ok {
			return v, ok
		}
	}
	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return false, false
}

func (ectx *ExecutionContext) GetFloat64(name string) (float64, bool) {
	item := ectx.Item()
	if item == nil {
		return 0, false
	}

	if val, ok := item.Get(name); ok {
		if v, ok := val.(float64); ok {
			return v, ok
		}
	}
	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return 0.0, false
}

func (ectx *ExecutionContext) GetTime(name string) (model.Time, bool) {
	item := ectx.Item()
	if item == nil {
		return model.Time(0), false
	}

	if val, ok := item.Get(name); ok {
		if v, ok := val.(model.Time); ok {
			return v, ok
		}
	}

	ectx.Errors().Add("builtins", errors.New("无效或不支持的字段 "+name))
	return model.Time(0), false
}

func (ectx *ExecutionContext) Clone() Context {
	var beclone = ExecutionContext{}

	if ectx.item != nil {
		becp := *ectx.item
		beclone.item = &becp
	}

	if ectx.order != nil {
		becp := *ectx.order
		becp.Items = make([]model.OrderItem, 0, len(ectx.order.Items))
		for _, itm := range ectx.order.Items {
			becp.Items = append(becp.Items, itm)
		}
		beclone.order = &becp
	}

	if ectx.user != nil {
		becp := *ectx.user
		beclone.user = &becp
	}

	beclone.options = ectx.options
	if ectx.errors != nil {
		becp := *ectx.errors
		beclone.errors = &becp
	}

	return &beclone
}

// func (ectx *ExecutionContext) Apply() (Context, error) {
// var (
// 	action ActionHandler = defaultActionHandler
// )
// if val, ok := ectx[applykey]; ok {
// 	if action, ok = val.(ActionHandler); !ok {
// 		action = defaultActionHandler // restore
// 	}
// }
// ctx := ectx.Clone()
// if err := action(ctx); err != nil {
// 	return ectx, err
// }
// return ectx, nil
// }

func (ectx *ExecutionContext) JSObject() map[string]interface{} {
	// ectx.unmarshalItem()

	var (
		jsobj = make(map[string]interface{})
		user  = ectx.User()
		order = ectx.Order()
		item  = ectx.Item()
		errs  = ectx.Errors()
		items = make([]map[string]interface{}, 0)
	)

	if user != nil {
		jsobj["user"] = map[string]interface{}{
			"username": user.Username,
			"avatar":   user.Avatar,
			"level":    user.Level,
		}

		if u, ok := jsobj["user"].(map[string]interface{}); ok {
			if !user.ID.Nil() {
				u["id"] = user.ID.Val
			}
		}
	}

	if item != nil {
		jitem := map[string]interface{}{
			"title":       item.Title,
			"price":       item.Price,
			"quantity":    item.Quantity,
			"subtotal":    item.Subtotal,
			"productDate": item.ProduceDate,
		}
		jsobj["item"] = jitem

		// jsobj["effects"] = make(map[string][]interface{})

		effects := makeEffects(item)
		if len(effects) > 0 {
			jitem["effects"] = effects
		}

		if it, ok := jsobj["item"].(map[string]interface{}); ok {
			if !item.ID.Nil() {
				it["id"] = item.ID.Val
			}

			if !item.ProductID.Nil() {
				it["productId"] = item.ProductID.Val
			}
		}
	}

	if order != nil {

		jsobj["order"] = map[string]interface{}{
			"count": order.Count,
			"total": order.Total,
			// "items": make([]map[string]interface{}),
			// "gifts": make([]map[string]interface{}),
		}

		for _, ite := range order.Items {
			iitm := map[string]interface{}{
				"title":       ite.Title,
				"price":       ite.Price,
				"quantity":    ite.Quantity,
				"subtotal":    ite.Subtotal,
				"productDate": ite.ProduceDate,
			}
			items = append(items, iitm)

			effects := makeEffects(&ite)
			if len(effects) > 0 {
				iitm["effects"] = effects
			}

			if !ite.ID.Nil() {
				iitm["id"] = ite.ID.Val
			}

			if !ite.ProductID.Nil() {
				iitm["productId"] = ite.ProductID.Val
			}
		}

		if ord, ok := jsobj["order"].(map[string]interface{}); ok {
			ord["items"] = items
		}
	}

	if errs != nil {
		jsobj["errors"] = errs.JSObject()
	}
	return jsobj
}

func makeEffects(item *model.OrderItem) map[string][]interface{} {
	effects := make(map[string][]interface{})
	for _, eff := range item.Effects() {
		jeff := make(map[string]interface{})
		jeff["field"] = eff.Field
		jeff["summary"] = eff.Summary
		jeff["value"] = eff.Val
		if _, ok := effects[eff.Field]; !ok {
			effects[eff.Field] = make([]interface{}, 0)
		}
		effects[eff.Field] = append(effects[eff.Field], jeff)
	}

	// if len(effects) > 0 {
	// 	jitem["effects"] = effects
	// }

	return effects
}
func isUndefined(obj *js.Object) bool {
	return jsbuiltin.TypeOf(obj) == "undefined"
}

func (ectx *ExecutionContext) LoadJS(state *js.Object) {
	if usr := state.Get("user"); !isUndefined(usr) {
		ectx.user = &model.User{
			ID:       fromId(usr.Get("id")),
			Username: usr.Get("username").String(),
			Avatar:   usr.Get("avatar").String(),
			Level:    usr.Get("level").Int(),
		}
	}

	if item := state.Get("item"); !isUndefined(item) {
		ite := &model.OrderItem{
			ID:          fromId(item.Get("id")),
			Price:       item.Get("price").Float(),
			Quantity:    item.Get("quantity").Int(),
			Subtotal:    item.Get("subtotal").Float(),
			ProductID:   fromId(item.Get("productId")),
			Title:       item.Get("title").String(),
			ProduceDate: model.Time(item.Get("productDate").Int()),
		}

		ectx.item = ite
		// ectx.marshalItem(ite)
	}

	if order := state.Get("order"); !isUndefined(order) {
		ord := &model.Order{
			// ID:    fromId(order.Get("id")),
			Total: order.Get("total").Float(),
			Count: order.Get("count").Int(),
		}
		items := order.Get("items")
		for i := 0; i < items.Length(); i++ {
			it := items.Index(i)
			ite := model.OrderItem{
				ID:          fromId(it.Get("id")),
				Price:       it.Get("price").Float(),
				Quantity:    it.Get("quantity").Int(),
				Subtotal:    it.Get("subtotal").Float(),
				ProductID:   fromId(it.Get("productId")),
				Title:       it.Get("title").String(),
				ProduceDate: model.Time(it.Get("productDate").Int()),
			}
			ord.Items = append(ord.Items, ite)
		}
		ectx.order = ord
	}
}

func fromId(val *js.Object) *model.Identity {
	typ := jsbuiltin.TypeOf(val)
	if typ == "string" {
		return model.StringID(val.String())
	} else if typ == "number" {
		return model.IntID(val.Int())
	} else {
		return nil
	}
}

func (ectx *ExecutionContext) getOptions() ExecutionOptions {
	return ectx.options
}

func (ectx *ExecutionContext) Now() model.Time {
	opt := ectx.getOptions()
	if opt.NowHandler == nil {
		ectx.Errors().Add("builtins", errors.New("没有设置定 NowHandler"))
		return 0
	}

	return model.Time(opt.NowHandler())
}

func (ectx *ExecutionContext) upsidePropagation() {
	order := ectx.Order()
	item := ectx.Item()
	if order != nil && item != nil {
		model.Propagation(order, item)
	}
}

func (ectx *ExecutionContext) Propagation() {
	ectx.upsideProp = true
}

func mustProduct(val interface{}, ok bool) model.Product {
	if ok {
		if prod, ok := val.(model.Product); ok {
			return prod
		}
	}
	return model.Product{}
}

// AddToOrder 添加/减少产品到订单
func AddToOrder(ctx Context) error {
	order := ctx.Order()
	item := ctx.Item()
	if idx, insert := utils.FindItemIndex(order.Items, item); insert {
		if item.Quantity > 0 {
			return order.AppendItem(item)
		} else {
			return errors.New("无效的数量")
		}
	} else {
		itm := &order.Items[idx]
		quantity := itm.Quantity + item.Quantity
		if quantity < 0 {
			return order.RemoveItem(item)
		}

		itm.Set("Quantity", quantity)
	}

	return nil
}

func remove(idx int, items []model.OrderItem) []model.OrderItem {
	if idx >= len(items)-1 {
		return items[:len(items)-1]
	}
	return append(items[:idx], items[idx+1:]...)
}

// RemoveFromOrder 移除产品
func RemoveFromOrder(ctx Context) error {
	order := ctx.Order()
	item := ctx.Item()

	if idx, _ := utils.FindItemIndex(order.Items, item); idx >= 0 {
		order.Items = remove(idx, order.Items)
	} else {
		return errors.New("产品没有在表单中")
	}

	return nil
}
