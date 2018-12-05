package model

import "strconv"

// ID 类型
type Identity struct {
	Val interface{}
}

// ID 接口
type ID interface {
	Compare(be *Identity) bool
}

// Compare ID 多态类型比较
func (a *Identity) Compare(b *Identity) bool {
	switch v := a.Val.(type) {
	case int:
		if a, ok := b.Val.(int); ok {
			return v == a
		}
		return false
	case string:
		if a, ok := b.Val.(string); ok {
			return v == a
		}
		return false
	default:
		return false
	}
}

func (a *Identity) String() string {
	switch v := a.Val.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return "[unknown]"
	}
}

func (a *Identity) Nil() bool {
	if a == nil {
		return true
	}
	return a.Val == nil
}

// IntID 整型 ID
func IntID(i int) *Identity {
	return &Identity{Val: i}
}

// StringID 字符型 ID
func StringID(i string) *Identity {
	return &Identity{Val: i}
}
