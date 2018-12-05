package constraint

type Constrainter interface {
	Name() string
	Evaluate(ctx Context) error
}

type poster struct {
	Constrainter
}

// Poster 后置条件，使约束从验证流中脱离，在其它约束完成后，再验证
func Poster(con Constrainter) Constrainter {
	return &poster{con}
}

// IsPoster 是否后置条件
func IsPoster(con Constrainter) bool {
	if _, ok := con.(*poster); ok {
		return true
	}

	return false
}
