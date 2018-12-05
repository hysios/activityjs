package context

// ActivityContext 活动上下文
type ActivityContext struct {
	Instock      int // 库存
	DefaultPrice int // 默认价格
	Name         string
	summaries    map[string]string
	Metadata     ActivityMetadata
	AppliedLimit int
}

type ActivityMetadata struct {
	Condition string      // 条件
	Kind      string      // 类型
	Offer     interface{} // 优惠
	Target    string
}

func NewActivityContext() *ActivityContext {
	return &ActivityContext{
		summaries: make(map[string]string),
	}
}

func (actx *ActivityContext) SetSummary(key, summary string) {
	actx.summaries[key] = summary
}

type RewardContext struct {
	ActivityContext *ActivityContext
	Ctx             Context
}

func (rctx *RewardContext) GetSummary(field string) string {
	if rctx.ActivityContext == nil {
		return ""
	}

	if summary, ok := rctx.ActivityContext.summaries[field]; ok {
		return summary
	}
	return ""
}

func (rctx *RewardContext) WithEffect(name, field string, val interface{}, fn func() error) error {
	var (
		ctx     = rctx.Ctx
		key     = rctx.ActivityContext.Name + name
		summary = rctx.GetSummary(field)
	)

	return ctx.Item().WithEffect(key, field, summary, val, func() error {
		return fn()
	})
}

func (rctx *RewardContext) UpsidePropagation(fn func()) {
	// model.UpsidePropagation(fn)
}
