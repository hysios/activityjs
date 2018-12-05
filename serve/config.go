//+build !js

package serve

import (
	"time"

	"activityjs.io/serve/context"
)

type Config struct {
	DefaultExecutionOptions context.ExecutionOptions
}

func DefaultConfig() Config {
	return Config{
		DefaultExecutionOptions: context.ExecutionOptions{
			NowHandler: NowHandler,
		},
	}
}

func NowHandler() int64 {
	return time.Now().Unix() * 1000
}
