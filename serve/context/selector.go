package context

import (
	"strconv"
	"strings"

	"activityjs.io/serve/errors"
	"activityjs.io/serve/model"
)

type Selector struct {
	*ExecutionContext
}

type SelectorPart int

const (
	Nosub SelectorPart = iota
	Subonly
	Subindex
)

var Subitems = []string{
	"User",
	"Order",
	"Order.Items",
	"Errors",
}

type (
	OutIDFunc      func(interface{}, string) (*model.Identity, bool)
	OutIntFunc     func(interface{}, string) (int, bool)
	OutInt64Func   func(interface{}, string) (int64, bool)
	OutStringFunc  func(interface{}, string) (string, bool)
	OutBoolFunc    func(interface{}, string) (bool, bool)
	OutFloat64Func func(interface{}, string) (float64, bool)
	OutTimeFunc    func(interface{}, string) (model.Time, bool)

	InputFunc func(interface{}, string) (interface{}, bool)
)

func NewSelector(exeCtx *ExecutionContext) *Selector {
	return &Selector{exeCtx}
}

func (s *Selector) getUser(val interface{}, key string) (interface{}, bool) {
	usr, ok := val.(*model.User)
	if !ok {
		return nil, false
	}

	return usr.Get(key)
}

func (s *Selector) setUser(usr *model.User, key string, val interface{}) error {

	return usr.Set(key, val)
}

func (s *Selector) getOrder(val interface{}, key string) (interface{}, bool) {
	order, ok := val.(*model.Order)
	if !ok {
		return nil, false
	}
	return order.Get(key)
}

func (s *Selector) setOrder(order *model.Order, key string, val interface{}) error {
	return order.Set(key, val)
}

func (s *Selector) getItem(val interface{}, key string) (interface{}, bool) {
	item, ok := val.(*model.OrderItem)
	if !ok {
		return nil, false
	}

	return item.Get(key)
}

func (s *Selector) setItem(item *model.OrderItem, key string, val interface{}) error {

	return item.Set(key, val)
}

func (s *Selector) getErrors(val interface{}, key string) (interface{}, bool) {
	errs, ok := val.(*errors.Errors)
	if !ok {
		return nil, false
	}

	err := errs.Get(key)
	if err != nil {
		return err, true
	}
	return nil, false
}

func (s *Selector) addErrors(errors *errors.Errors, key string, err error) error {
	errors.Add(key, err)
	return nil
}

func (s *Selector) parseSelector(sel string) []string {
	pos := strings.LastIndex(sel, ".")
	if pos > 0 {
		return []string{sel[:pos], sel[pos+1:]}
	}
	return nil
}

func (s *Selector) parseSubindex(sel string) ([]string, int) {
	lpos := strings.LastIndex(sel, "[")
	rpos := strings.LastIndex(sel, "]")
	if rpos-lpos > 1 {
		sub, key := sel[:lpos], sel[rpos+2:]
		if idx, err := strconv.Atoi(sel[lpos+1 : rpos]); err != nil {
			return nil, -1
		} else {
			return []string{sub, key}, idx
		}
	}

	return nil, -1
}

func (s *Selector) isSubindex(sel string) (ok bool) {
	lpos := strings.LastIndex(sel, "[")
	rpos := strings.LastIndex(sel, "]")
	return rpos-lpos > 1
}

func (s *Selector) subString(sel string) (sub, key string, ok bool) {
	ss := s.parseSelector(sel)
	if ss == nil {
		return "", "", false
	}

	return ss[0], ss[1], true
}

func (s *Selector) subindexString(sel string) (sub, key string, idx int, ok bool) {
	ss, idx := s.parseSubindex(sel)
	if ss == nil || idx == -1 {
		return "", "", -1, false
	}

	return ss[0], ss[1], idx, true
}

func (s *Selector) isSubkey(sub string) bool {
	for _, key := range Subitems {
		if key == sub {
			return true
		}
	}

	return false
}

func (s *Selector) selector(sel string, fn func(par SelectorPart, sub, key string, idx int)) error {
	if s.isSubindex(sel) {
		if sub, key, idx, ok := s.subindexString(sel); ok {
			fn(Subindex, sub, key, idx)
		} else {
			return errors.New("无效的 Subindex 格式")
		}
	} else if sub, key, ok := s.subString(sel); ok {
		fn(Subonly, sub, key, -1)
	} else {
		fn(Nosub, "", sel, -1)
	}
	return nil
}

func (s *Selector) MakeID(fn InputFunc) OutIDFunc {
	return func(obj interface{}, key string) (*model.Identity, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(*model.Identity); ok {
				return v, true
			} else {
				return nil, false
			}
		}
		return nil, false
	}
}

func (s *Selector) MakeInt(fn InputFunc) OutIntFunc {
	return func(obj interface{}, key string) (int, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(int); ok {
				return v, true
			} else {
				return 0, false
			}
		}
		return 0, false
	}
}

func (s *Selector) MakeInt64(fn InputFunc) OutInt64Func {
	return func(obj interface{}, key string) (int64, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(int64); ok {
				return v, true
			} else {
				return 0, false
			}
		}
		return 0, false
	}
}

func (s *Selector) MakeString(fn InputFunc) OutStringFunc {
	return func(obj interface{}, key string) (string, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(string); ok {
				return v, true
			} else {
				return "", false
			}
		}
		return "", false
	}
}

func (s *Selector) MakeBool(fn InputFunc) OutBoolFunc {
	return func(obj interface{}, key string) (bool, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(bool); ok {
				return v, true
			} else {
				return false, false
			}
		}
		return false, false
	}
}

func (s *Selector) MakeTime(fn InputFunc) OutTimeFunc {
	return func(obj interface{}, key string) (model.Time, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(model.Time); ok {
				return v, true
			} else {
				return 0, false
			}
		}
		return 0, false
	}
}

func (s *Selector) MakeFloat64(fn InputFunc) OutFloat64Func {
	return func(obj interface{}, key string) (float64, bool) {
		if val, ok := fn(obj, key); ok {
			if v, ok := val.(float64); ok {
				return v, true
			} else {
				return 0, false
			}
		}
		return 0, false
	}
}

func (s *Selector) subget(sub string, fn func(interface{}, InputFunc)) {
	exeCtx := s.ExecutionContext

	switch sub {
	case "User":
		usr := exeCtx.User()
		fn(usr, s.getUser)
	case "Order":
		order := exeCtx.Order()
		fn(order, s.getOrder)
	case "Errors":
		errs := exeCtx.Errors()
		fn(errs, s.getErrors)
	default:
	}
}

func (s *Selector) subidxget(sub string, idx int, fn func(interface{}, InputFunc)) (err error) {
	exeCtx := s.ExecutionContext

	switch sub {
	case "Order.Items":
		order := exeCtx.Order()
		items := order.Items
		if !(idx >= 0 && idx < len(items)) {
			err = errors.New("无效的 index")
			s.Errors().Add("builtins", err)
			return
		}
		fn(&items[idx], s.getItem)
	default:
		err = errors.New("无效的 index 字段 " + sub)
		s.Errors().Add("builtins", err)
		return
	}

	return nil
}

func (s *Selector) Set(name string, val interface{}) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			switch sub {
			case "Order.Items":
				order := s.Order()
				items := order.Items

				if !(idx >= 0 && idx < len(items)) {
					err := errors.New("无效的 index")
					s.Errors().Add("builtins", err)
					return
				}
				s.setItem(&items[idx], key, val)
			}
		case Subonly:
			switch sub {
			case "User":
				usr := s.User()
				s.setUser(usr, key, val)
			case "Order":
				order := s.Order()
				s.setOrder(order, key, val)
			case "Errors":
				errs := s.Errors()
				err := val.(error)
				s.addErrors(errs, key, err)
			default:
			}
		case Nosub:
			ctx := s.ExecutionContext
			ctx.Set(name, val)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
	}
}

func (s *Selector) Get(name string) (val interface{}, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				val, ok = getMth(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				val, ok = getMth(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.Get(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return nil, false
	}
	return
}

func (s *Selector) GetID(name string) (val *model.Identity, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				idGet := s.MakeID(getMth)
				val, ok = idGet(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				idGet := s.MakeID(getMth)
				val, ok = idGet(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.GetID(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return nil, false
	}
	return nil, false
}

func (s *Selector) GetInt(name string) (val int, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				intGet := s.MakeInt(getMth)
				val, ok = intGet(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				intGet := s.MakeInt(getMth)
				val, ok = intGet(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.GetInt(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return 0, false
	}
	return
}

func (s *Selector) GetInt64(name string) (val int64, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				intGet := s.MakeInt64(getMth)
				val, ok = intGet(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				intGet := s.MakeInt64(getMth)
				val, ok = intGet(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.GetInt64(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return 0, false
	}

	return
}

func (s *Selector) GetString(name string) (val string, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				typGet := s.MakeString(getMth)
				val, ok = typGet(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				typGet := s.MakeString(getMth)
				val, ok = typGet(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.GetString(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return "", false
	}

	return
}

func (s *Selector) GetBool(name string) (val bool, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				typGet := s.MakeBool(getMth)
				val, ok = typGet(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				typGet := s.MakeBool(getMth)
				val, ok = typGet(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.GetBool(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return false, false
	}

	return
}

func (s *Selector) GetFloat64(name string) (val float64, ok bool) {
	if err := s.selector(name, func(par SelectorPart, sub, key string, idx int) {
		switch par {
		case Subindex:
			s.subidxget(sub, idx, func(obj interface{}, getMth InputFunc) {
				typGet := s.MakeFloat64(getMth)
				val, ok = typGet(obj, key)
			})
		case Subonly:
			s.subget(sub, func(obj interface{}, getMth InputFunc) {
				typGet := s.MakeFloat64(getMth)
				val, ok = typGet(obj, key)
			})
		case Nosub:
			exeCtx := s.ExecutionContext
			val, ok = exeCtx.GetFloat64(key)
		}
	}); err != nil {
		s.Errors().Add("builtins", err)
		return 0, false
	}

	return
}

func (s *Selector) Clone() Context {
	if ectx, ok := s.ExecutionContext.Clone().(*ExecutionContext); ok {
		return &Selector{ectx}
	}
	return nil
}

// func (s *Selector) Apply() (Context, error) {
// 	ctx, err := s.ExecutionContext.Apply()
// 	if err != nil {
// 		return s, err
// 	}

// 	if ectx, ok := ctx.(*ExecutionContext); ok {

// 		return &Selector{ectx}, nil
// 	}

// 	return nil, nil
// }
