// Package cflag provides a flag-like means of declaring configurables.
//
// The functions in this package which take a Registerable argument can
// have that argument passed as non-nil, in which case the configurable
// becomes a child of the configurable passed, or nil, in which case
// the configurable is registered at the top level.
//
// You should call Value() to get the value of a flag configurable.
package cflag

import "fmt"
import "strconv"
import "regexp"
import "strings"
import "gopkg.in/hlandau/configurable.v1"

// Group

type Registerable interface {
	Register(configurable configurable.Configurable)
}

type noReg struct{}

// Dummy Registerable implementation which does not do anything.
//
// Can be used to inhibit autoregistration.
var NoReg noReg

func (r *noReg) Register(configurable configurable.Configurable) {

}

func register(r Registerable, c configurable.Configurable) {
	if r == nil {
		configurable.Register(c)
	} else {
		r.Register(c)
	}
}

type Group struct {
	configurables []configurable.Configurable
	name          string
}

func (ig *Group) CfName() string {
	return ig.name
}

func (ig *Group) CfChildren() []configurable.Configurable {
	return ig.configurables
}

func (ig *Group) String() string {
	return fmt.Sprintf("%s", ig.name)
}

// Register a child configurable to the group.
func (ig *Group) Register(cfg configurable.Configurable) {
	ig.configurables = append(ig.configurables, cfg)
}

// Creates a flag group. A Group is itself a configurable and can hold multiple
// flags.
func NewGroup(reg Registerable, name string) *Group {
	ig := &Group{
		name: name,
	}
	register(reg, ig)
	return ig
}

// String

type StringFlag struct {
	name, curValue, summaryLine, defaultValue string
	curValuep                                 *string
	priority                                  configurable.Priority
	onChange                                  []func(*StringFlag)
}

func (sf *StringFlag) String() string {
	return fmt.Sprintf("SimpleFlag(%s: %#v)", sf.name, *sf.curValuep)
}

func (sf *StringFlag) CfSetValue(v interface{}) error {
	defer sf.notify()

	vs, ok := v.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	*sf.curValuep = vs
	return nil
}

func (sf *StringFlag) notify() {
	for _, f := range sf.onChange {
		f(sf)
	}
}

func (sf *StringFlag) CfValue() interface{} {
	return *sf.curValuep
}

func (sf *StringFlag) CfName() string {
	return sf.name
}

func (sf *StringFlag) CfUsageSummaryLine() string {
	return sf.summaryLine
}

func (sf *StringFlag) CfDefaultValue() interface{} {
	return sf.defaultValue
}

// Get the flag's current value.
func (sf *StringFlag) Value() string {
	return *sf.curValuep
}

// Set the flag's current value.
func (sf *StringFlag) SetValue(value string) {
	*sf.curValuep = value
}

func (sf *StringFlag) RegisterOnChange(f func(*StringFlag)) {
	sf.onChange = append(sf.onChange, f)
}

func (sf *StringFlag) CfSetPriority(priority configurable.Priority) {
	sf.priority = priority
}

func (sf *StringFlag) CfGetPriority() configurable.Priority {
	return sf.priority
}

// Creates a flag of type string. The variable referenced by pointer v is used as
// the storage location for the value of the configurable.
func StringVar(reg Registerable, v *string, name, defaultValue, summaryLine string) *StringFlag {
	sf := &StringFlag{
		name:         name,
		summaryLine:  summaryLine,
		defaultValue: defaultValue,
		curValue:     defaultValue,
		curValuep:    v,
	}
	if sf.curValuep == nil {
		sf.curValuep = &sf.curValue
	}

	register(reg, sf)
	return sf
}

// Creates a flag of type string.
//
// reg: See package-level documentation.
//
// summaryLine: One-line usage summary.
func String(reg Registerable, name, defaultValue, summaryLine string) *StringFlag {
	return StringVar(reg, nil, name, defaultValue, summaryLine)
}

// Int

type IntFlag struct {
	name, summaryLine      string
	curValue, defaultValue int
	curValuep              *int
	priority               configurable.Priority
	onChange               []func(*IntFlag)
}

func (sf *IntFlag) String() string {
	return fmt.Sprintf("IntFlag(%s: %#v)", sf.name, *sf.curValuep)
}

func (sf *IntFlag) CfSetValue(v interface{}) error {
	defer sf.notify()

	vi, ok := v.(int)
	if ok {
		*sf.curValuep = vi
		return nil
	}

	vs, ok := v.(string)
	if ok {
		vs = strings.TrimSpace(vs)
		n, err := strconv.ParseInt(vs, 0, 32)
		if err != nil {
			return err
		}

		*sf.curValuep = int(n)
		return nil
	}

	return fmt.Errorf("invalid value for configurable %#v, expecting int: %v", sf.name, v)
}

func (sf *IntFlag) notify() {
	for _, f := range sf.onChange {
		f(sf)
	}
}

func (sf *IntFlag) CfValue() interface{} {
	return sf.curValue
}

func (sf *IntFlag) CfName() string {
	return sf.name
}

func (sf *IntFlag) CfUsageSummaryLine() string {
	return sf.summaryLine
}

func (sf *IntFlag) CfDefaultValue() interface{} {
	return sf.defaultValue
}

// Get the flag's current value.
func (sf *IntFlag) Value() int {
	return *sf.curValuep
}

// Set the flag's current value.
func (sf *IntFlag) SetValue(value int) {
	*sf.curValuep = value
}

func (sf *IntFlag) RegisterOnChange(f func(*IntFlag)) {
	sf.onChange = append(sf.onChange, f)
}

func (sf *IntFlag) CfSetPriority(priority configurable.Priority) {
	sf.priority = priority
}

func (sf *IntFlag) CfGetPriority() configurable.Priority {
	return sf.priority
}

// Creates a flag of type int. The variable referenced by pointer v is used as
// the storage location for the value of the configurable.
func IntVar(reg Registerable, v *int, name string, defaultValue int, summaryLine string) *IntFlag {
	sf := &IntFlag{
		name:         name,
		summaryLine:  summaryLine,
		defaultValue: defaultValue,
		curValue:     defaultValue,
		curValuep:    v,
	}
	if sf.curValuep == nil {
		sf.curValuep = &sf.curValue
	}

	register(reg, sf)
	return sf
}

// Creates a flag of type int.
//
// reg: See package-level documentation.
//
// summaryLine: One-line usage summary.
func Int(reg Registerable, name string, defaultValue int, summaryLine string) *IntFlag {
	return IntVar(reg, nil, name, defaultValue, summaryLine)
}

// Bool

type BoolFlag struct {
	name, summaryLine      string
	curValue, defaultValue bool
	curValuep              *bool
	priority               configurable.Priority
	onChange               []func(*BoolFlag)
}

func (sf *BoolFlag) String() string {
	return fmt.Sprintf("BoolFlag(%s: %#v)", sf.name, sf.curValue)
}

var re_no = regexp.MustCompilePOSIX(`^(00?|no?|f(alse)?)$`)

func (sf *BoolFlag) CfSetValue(v interface{}) error {
	defer sf.notify()

	vb, ok := v.(bool)
	if ok {
		*sf.curValuep = vb
		return nil
	}

	vi, ok := v.(int)
	if ok {
		*sf.curValuep = (vi != 0)
		return nil
	}

	vs, ok := v.(string)
	if ok {
		vs = strings.TrimSpace(vs)
		*sf.curValuep = !re_no.MatchString(vs)
		return nil
	}

	return fmt.Errorf("invalid value for configurable %#v, expecting bool: %v", sf.name, v)
}

func (sf *BoolFlag) notify() {
	for _, f := range sf.onChange {
		f(sf)
	}
}

func (sf *BoolFlag) CfValue() interface{} {
	return sf.curValue
}

func (sf *BoolFlag) CfName() string {
	return sf.name
}

func (sf *BoolFlag) CfUsageSummaryLine() string {
	return sf.summaryLine
}

func (sf *BoolFlag) CfDefaultValue() interface{} {
	return sf.defaultValue
}

// Call to get the flag's current value.
func (sf *BoolFlag) Value() bool {
	return *sf.curValuep
}

// Set the flag's current value.
func (sf *BoolFlag) SetValue(value bool) {
	*sf.curValuep = value
}

func (sf *BoolFlag) RegisterOnChange(f func(*BoolFlag)) {
	sf.onChange = append(sf.onChange, f)
}

func (sf *BoolFlag) CfSetPriority(priority configurable.Priority) {
	sf.priority = priority
}

func (sf *BoolFlag) CfGetPriority() configurable.Priority {
	return sf.priority
}

// Creates a flag of type bool.
//
// reg: See package-level documentation.
//
// summaryLine: One-line usage summary.
func Bool(reg Registerable, name string, defaultValue bool, summaryLine string) *BoolFlag {
	return BoolVar(reg, nil, name, defaultValue, summaryLine)
}

// Creates a flag of type bool. The variable referenced by pointer v is used as
// the storage location for the value of the configurable.
func BoolVar(reg Registerable, v *bool, name string, defaultValue bool, summaryLine string) *BoolFlag {
	sf := &BoolFlag{
		name:         name,
		summaryLine:  summaryLine,
		defaultValue: defaultValue,
		curValue:     defaultValue,
		curValuep:    v,
	}
	if sf.curValuep == nil {
		sf.curValuep = &sf.curValue
	}

	register(reg, sf)
	return sf
}
