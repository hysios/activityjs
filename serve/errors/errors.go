package errors

// Errorer 定制错误接口
type Errorer string

func (e Errorer) Error() string {
	return string(e)
}

// New 模拟 runtime new errors
func New(msg string) error {
	return Errorer(msg)
}

var (
	// ErrAbort 终止错误，停止验证检查，返回失败
	ErrAbort = New("Abort Evalution")
	// ErrExit Errors 退出标志
	ErrExit   = New("Exit")
	ErrNotice = New("Notice")
)

// Errors 错误集合
type Errors map[string]error

// MarshalJSON json 格式化输出
func (errs *Errors) JSObject() map[string]interface{} {
	var result = make(map[string]interface{})
	for key, err := range *errs {
		result[key] = err.Error()
	}

	return result
}

// Error 错误接口
func (errs *Errors) Error() string {

	output := ""
	if len(*errs) > 0 {
		output += "发生以下错误\n"
	}
	output = errs.msg(output, 1)
	return output
}

// Empty Errors 为空
func (errs Errors) Empty() bool {
	return len(errs) == 0
}

func (errs *Errors) msg(output string, depth int) string {

	for key, err := range *errs {
		for i := 0; i < depth; i++ {
			output += "\t"
		}
		output += "[" + key + "] \t"
		if errs := AsErrors(err); errs != nil {
			depth++
			output = errs.msg(output, depth)
		} else {
			output += err.Error() + "\n"
		}
	}
	return output
}

func (errs *Errors) traversal(fn func(key string, err error) error) error {
	for key, err := range *errs {
		if errs := AsErrors(err); errs != nil {
			if err := errs.traversal(fn); err != nil {
				return err
			}
		} else {
			if err := fn(key, err); err != nil {
				return err
			}
		}
	}

	return nil
}

func AsErrors(err error) *Errors {
	if errs, ok := err.(*Errors); ok {
		return errs
	}
	return nil
}

// NotStop 没有停止的错误，一般性 error 在这里被当成警告，
// 只有派生出 ErrExit 的错误被视为停止
func (errs Errors) NotStop() bool {
	if len(errs) == 0 {
		return true
	}

	var stop = false

	errs.traversal(func(key string, err error) error {
		if err := AsWrap(err); err != nil {
			if err.Is(ErrExit) {
				stop = true
				return err
			}
		}

		if err == ErrExit {
			stop = true
			return err
		}
		return nil
	})

	return !stop
}

// Add 增加子节点错误
func (errs Errors) Add(name string, err error) {
	errs[name] = err
}

func (errs Errors) Get(name string) error {
	return errs[name]
}

func Wrap(err error, msg string) error {
	return &WrapError{msg, err}
}

type WrapError struct {
	msg  string
	wrap error
}

func AsWrap(err error) *WrapError {
	if err, ok := err.(*WrapError); ok {
		return err
	}
	return nil
}

func (err *WrapError) Error() string {
	return err.msg + " " + err.wrap.Error()
}

func (err *WrapError) Is(cmp error) bool {
	werr := err.Unwrap()

	if werr == cmp {
		return true
	}

	if werr := AsWrap(werr); err != nil {
		return werr.Is(cmp)
	}

	return false
}

func (err *WrapError) Unwrap() error {
	return err.wrap
}
