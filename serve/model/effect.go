package model

type Effect struct {
	items []EffectItem
}

type EffectItem struct {
	Key     string
	Field   string
	Summary string
	Val     interface{}
}

func (e *Effect) AddEffect(key, field, summary string, val interface{}) {
	e.items = append(e.items, EffectItem{key, field, summary, val})
}

func (e *Effect) FieldEffects(field string) []EffectItem {
	var rslt = make([]EffectItem, 0)
	for _, itm := range e.items {
		if itm.Field == field {
			rslt = append(rslt, itm)
		}
	}

	return rslt
}

func (e *Effect) Effects() []EffectItem {
	return e.items
}
