package reward

import "activityjs.io/serve/context"

func Joint(fns ...Handler) Handler {
	return func(ctx context.RewardContext) error {
		for _, fn := range fns {
			err := fn(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
