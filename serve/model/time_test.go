package model_test

import (
	"testing"

	"activityjs.io/serve/model"
)

func TestTime(t *testing.T) {
	tt := model.Time(1542249110944)
	t.Logf("tt %s", tt)
	tt = tt.Timezone(8)

	t.Logf("tt %s", tt)

	tt = model.Time(1542268606879)
	t.Logf("tt %s", tt)
	tt = tt.Timezone(8)

	t.Logf("tt %s", tt)
}
