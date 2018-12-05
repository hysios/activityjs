package serve

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"activityjs.io/serve/activity"
	"activityjs.io/serve/context"
	"activityjs.io/serve/model"

	"activityjs.io/serve/utils"
	"github.com/fatih/color"
)

// TestSceneInWaterfall 在活动流中的工作
func TestFail(t *testing.T) {
	var (
		_, activityItems, orderItems,
		user, order = utils.Prepare()
		actIphone = activityItems[0]
		actIpad   = activityItems[1]
		item      = orderItems[0]
		pp        = context.TableEx(os.Stderr)
	)

	defultCfg := DefaultConfig()
	var (
		ctx context.Context = context.New(&item, &order, &user, defultCfg.DefaultExecutionOptions)
		err error
	)

	mach := NewMachine(ctx)
	mach.AddActivity(activity.SpecialPrice(&actIphone, 20, 3.0))
	mach.AddActivity(activity.OverDecrease(&actIpad, 5000.0, 100))

	ctx, err = mach.Evaluate(ctx)
	if err != nil {
		t.Logf("Machine Evaluate: %s", err)
	}
	pp(ctx)
}

func TestSuccess(t *testing.T) {
	var (
		_, activityItems, orderItems,
		user, order = utils.Prepare()
		// actIphone = activityItems[0]
		actIpad = activityItems[1]
		item    = orderItems[0]
		pp      = utils.MakeTable(os.Stderr)
	)

	defultCfg := DefaultConfig()
	var (
		ctx context.Context = context.New(&item, &order, &user, defultCfg.DefaultExecutionOptions)
		err error
	)

	mach := NewMachine(ctx)
	mach.AddActivity(activity.SpecialPrice(&actIpad, 2, 3.0))
	// mach.AddActivity(activity.OverDecrease(&actIpad, 5000.0, 100))

	ctx.Set("Quantity", 2)
	ctx, err = mach.Evaluate(ctx)
	if err != nil {
		t.Logf("Machine Evaluate: %s", err)
	}
	color.Set(color.FgYellow)

	// t.Logf("Order: %# v", pretty.Formatter(ctx.Order()))
	log.Printf("ItemTable")
	pp.Print(fromItem(ctx.Item()))
	log.Printf("OrderTable")
	pp.Print(ctx.Order())

	pp.Print(ctx.Order().Items)
	log.Printf("MachineTable")
	ppm(pp, mach)
}

type PrintMachine struct {
	Activities []PrintActivity `header:"activities"`
}

type PrintItem struct {
	ID       string `header:"id"`
	Title    string `header:"title"`
	Price    string `header:"price"`
	Quantity string `header:"quantity"`
	Subtotal string `header:"subtotal"`
	// Effects  string `header:"effects"`
}

type PrintActivityItem struct {
	Name string `header:"item name"`
	ID   string `header:"item id"`
}

type PrintActivity struct {
	Name            string            `header:"name"`
	Item            PrintActivityItem `header:"inline"`
	Type            string            `header:"type"`
	ConstraintCount int               `header:"constraint"`
	Constraint      []PrintConstraint
	Condition       string `header:"condition"`
	Offer           string `header:"Offer"`
}

type PrintConstraint struct {
	Name        string      `header:"name"`
	Description interface{} `header:"description"`
}

func fromAct(act *activity.Activity) PrintActivity {
	pact := PrintActivity{
		Name:            act.Name(),
		Item:            PrintActivityItem{Name: act.Item.Title, ID: act.Item.ItemID.String()},
		ConstraintCount: len(act.Constraints),
	}

	metadata := act.Metadata()
	if metadata != nil {
		pact.Type = metadata.Kind
		pact.Condition = metadata.Condition
		pact.Offer = fmt.Sprintf("优惠%s %+v", metadata.Target, metadata.Offer)
	}

	return pact
}

func fromItem(itm *model.OrderItem) PrintItem {
	pitem := PrintItem{
		ID:       fmt.Sprintf("%s", itm.ID),
		Title:    itm.Title,
		Quantity: withEffect(itm, "Quantity", itm.Quantity),
		Price:    withEffect(itm, "Price", itm.Price),
		Subtotal: withEffect(itm, "Subtotal", itm.Subtotal),
	}

	return pitem
}

func withEffect(itm *model.OrderItem, field string, val interface{}) string {
	var (
		effs   = itm.FieldEffects(field)
		ss     = make([]string, 0)
		effstr string
	)

	for _, eff := range effs {
		ss = append(ss, fmt.Sprintf("%s %v", eff.Summary, eff.Val))
	}

	if len(ss) == 0 {
		effstr = ""
	} else {
		s := strings.Join(ss, ",")
		effstr = fmt.Sprintf(" (%s)", color.RedString(s))
	}

	switch v := val.(type) {
	case int:
		return fmt.Sprintf("% 5d%s", v, effstr)
	case float64:
		return fmt.Sprintf("%0.2f%s", v, effstr)
	default:
		return fmt.Sprintf("%v%s", v, effstr)
	}
}

func ppm(pp *utils.Printer, mach *Machine) {
	var acts = make([]PrintActivity, 0)
	for _, act := range mach.activities {
		acts = append(acts, fromAct(act))
	}

	for _, act := range mach.posterActivities {
		acts = append(acts, fromAct(act))
	}
	pp.Print(acts)

	var conts = make([]PrintConstraint, 0)
	for _, act := range mach.activities {
		for _, con := range act.Constraints {
			conts = append(conts, PrintConstraint{con.Name(), con})
		}
	}
	for _, act := range mach.posterActivities {
		for _, con := range act.Constraints {
			conts = append(conts, PrintConstraint{con.Name(), con})
		}
	}

	pp.Print(conts)

}
