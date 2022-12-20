package args

import "reflect"

var factories map[reflect.Kind]optionFactory

func init() {
	factories = map[reflect.Kind]optionFactory{
		reflect.Bool: singleValuedOptionFactory(
			reflect.TypeOf(BooleanOption(false)),
			func(opt *reflect.Value, val interface{}) {
				opt.SetBool(val.(bool))
			}),
		reflect.Int: singleValuedOptionFactory(
			reflect.TypeOf(IntOption(0)),
			func(opt *reflect.Value, val interface{}) {
				opt.SetInt(int64(val.(int)))
			}),
		reflect.String: singleValuedOptionFactory(
			reflect.TypeOf(StringOption("")),
			func(opt *reflect.Value, val interface{}) {
				opt.SetString(val.(string))
			}),
		reflect.Slice: factoryFn(multiOptionsFactory),
	}
}

type optionFactory interface {
	create(value interface{}) interface{}
}

type factoryFn func(interface{}) interface{}

func (cf factoryFn) create(value interface{}) interface{} {
	return cf(value)
}

func singleValuedOptionFactory(typ reflect.Type, setter func(opt *reflect.Value, val interface{})) factoryFn {
	optVal := reflect.New(typ).Elem()
	f := func(v interface{}) interface{} {
		setter(&optVal, v)
		return optVal.Interface()
	}
	return factoryFn(f)
}

func multiOptionsFactory(value interface{}) interface{} {
	singleValuedOptions := map[reflect.Kind]interface{}{
		reflect.Bool:   BooleanOption(false),
		reflect.Int:    IntOption(0),
		reflect.String: StringOption(""),
	}

	vals := value.([]interface{})
	for _, val := range vals {
		kind := reflect.ValueOf(val).Kind()
		singleValuedOptions[kind] = factories[kind].create(val)
	}

	return MultiOptions{
		BooleanOption: singleValuedOptions[reflect.Bool].(BooleanOption),
		IntOption:     singleValuedOptions[reflect.Int].(IntOption),
		StringOption:  singleValuedOptions[reflect.String].(StringOption),
	}
}
