// +build debug

package log

import (
	"github.com/kr/pretty"
)

func Log(pat string, args ...interface{}) {
	pretty.Logf(pat, args...)
}
